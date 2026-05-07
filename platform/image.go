package platform

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"image"
	_ "image/gif"  // register GIF format
	_ "image/jpeg" // register JPEG format
	_ "image/png"  // register PNG format
	"io"
	"math"
	"time"

	"github.com/Gleipnir-Technology/nidus-sync/db"
	"github.com/Gleipnir-Technology/nidus-sync/db/gen/nidus-sync/publicreport/model"
	querypublicreport "github.com/Gleipnir-Technology/nidus-sync/db/query/publicreport"
	"github.com/Gleipnir-Technology/nidus-sync/geomutil"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"github.com/rwcarlsen/goexif/exif"
	"github.com/rwcarlsen/goexif/tiff"
	"github.com/twpayne/go-geom"
	//exif "github.com/rwcarlsen/goexif/exif"
	//"github.com/dsoprea/go-exif-extra/format"
)

type GPS struct {
	Latitude  float64
	Longitude float64
}

type ExifCollection struct {
	GPS  *GPS
	Tags map[string]string
}

type ImageUpload struct {
	Bounds      image.Rectangle
	ContentType string
	Exif        *ExifCollection
	Format      string

	UploadFilesize int
	UploadFilename string
	UUID           uuid.UUID
}

func (e *ExifCollection) Walk(name exif.FieldName, tag *tiff.Tag) error {
	e.Tags[string(name)] = tag.String()
	return nil
}
func ImageExtractExif(content_type string, file_bytes []byte) (result *ExifCollection, err error) {
	/*
		Using "github.com/evanoberholster/imagemeta"
		meta, err := imagemeta.Decode(bytes.NewReader(file_bytes))
		if err != nil {
			return result, fmt.Errorf("Failed to decode image meta: %w", err)
		}
		result.GPS = &GPS{
			Latitude: meta.GPS.Latitude(),
			Longitude: meta.GPS.Longitude(),
		}
		return result, err
	*/

	e, err := exif.Decode(bytes.NewReader(file_bytes))
	if err != nil {
		if err.Error() == "exif: failed to find exif intro marker" {
			return nil, nil
		} else if errors.Is(err, io.EOF) {
			return nil, nil
		}
		return nil, fmt.Errorf("Failed to decode image meta: %w", err)
	}
	lat, lng, _ := e.LatLong()
	result = &ExifCollection{
		GPS: &GPS{
			Latitude:  lat,
			Longitude: lng,
		},
		Tags: make(map[string]string, 0),
	}
	err = e.Walk(result)
	return result, err
}

func saveImageUploads(ctx context.Context, txn db.Ex, uploads []ImageUpload) ([]model.Image, error) {
	images := make([]model.Image, 0)
	for _, u := range uploads {
		var location *geom.T
		if u.Exif != nil && u.Exif.GPS != nil && !(math.IsNaN(u.Exif.GPS.Longitude) || math.IsNaN(u.Exif.GPS.Latitude)) {
			l := geomutil.PointFromLngLat(u.Exif.GPS.Longitude, u.Exif.GPS.Latitude)
			location = &l
		}
		image := model.Image{
			// ID:
			ContentType:      u.ContentType,
			Created:          time.Now(),
			Location:         location,
			ResolutionX:      int32(u.Bounds.Max.X),
			ResolutionY:      int32(u.Bounds.Max.Y),
			StorageUUID:      u.UUID,
			StorageSize:      int64(u.UploadFilesize),
			UploadedFilename: u.UploadFilename,
		}
		image, err := querypublicreport.ImageInsert(ctx, txn, image)
		if err != nil {
			return images, fmt.Errorf("Failed to create photo records: %w", err)
		}

		// TODO: figure out how to do this via the setter...?
		if u.Exif != nil {
			exif_models := make([]model.ImageExif, len(u.Exif.Tags))
			i := 0
			for k, v := range u.Exif.Tags {
				to_save := trimQuotes(v)
				exif_models[i] = model.ImageExif{
					ImageID: image.ID,
					Name:    k,
					Value:   to_save,
				}
			}
			if len(exif_models) > 0 {
				_, err = querypublicreport.ImageExifInserts(ctx, txn, exif_models)
				if err != nil {
					return images, fmt.Errorf("Failed to create photo exif records: %w", err)
				}
			}
			log.Info().Int32("id", image.ID).Int("tags", len(u.Exif.Tags)).Msg("Saved an uploaded file to the database")
		} else {
			log.Info().Int32("id", image.ID).Int("tags", 0).Msg("Saved an uploaded file without EXIF data")
		}
		images = append(images, image)
	}
	return images, nil
}

// Given a string like "\"foo\"" return "foo".
func trimQuotes(s string) string {
	if len(s) >= 2 && s[0] == '"' && s[len(s)-1] == '"' {
		return s[1 : len(s)-1]
	}
	return s
}
