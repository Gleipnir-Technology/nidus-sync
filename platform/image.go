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
	"time"

	"github.com/Gleipnir-Technology/bob"
	"github.com/Gleipnir-Technology/bob/dialect/psql"
	"github.com/Gleipnir-Technology/bob/dialect/psql/um"
	"github.com/Gleipnir-Technology/nidus-sync/db/models"
	"github.com/aarondl/opt/omit"
	"github.com/aarondl/opt/omitnull"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"github.com/rwcarlsen/goexif/exif"
	"github.com/rwcarlsen/goexif/tiff"
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

func saveImageUploads(ctx context.Context, tx bob.Tx, uploads []ImageUpload) (models.PublicreportImageSlice, error) {
	images := make(models.PublicreportImageSlice, 0)
	for _, u := range uploads {
		image, err := models.PublicreportImages.Insert(&models.PublicreportImageSetter{
			ContentType: omit.From(u.ContentType),

			Created: omit.From(time.Now()),
			//Location: 	psql.Raw("NULL"),
			Location:         omitnull.FromPtr[string](nil),
			ResolutionX:      omit.From(int32(u.Bounds.Max.X)),
			ResolutionY:      omit.From(int32(u.Bounds.Max.Y)),
			StorageUUID:      omit.From(u.UUID),
			StorageSize:      omit.From(int64(u.UploadFilesize)),
			UploadedFilename: omit.From(u.UploadFilename),
		}).One(ctx, tx)
		if err != nil {
			return images, fmt.Errorf("Failed to create photo records: %w", err)
		}

		// TODO: figure out how to do this via the setter...?
		if u.Exif != nil {
			if u.Exif.GPS != nil {
				_, err = psql.Update(
					um.Table("publicreport.image"),
					um.SetCol("location").To(fmt.Sprintf("ST_Point(%f, %f, 4326)", u.Exif.GPS.Longitude, u.Exif.GPS.Latitude)),
					um.Where(psql.Quote("id").EQ(psql.Arg(image.ID))),
				).Exec(ctx, tx)
			}

			exif_setters := make([]*models.PublicreportImageExifSetter, 0)
			for k, v := range u.Exif.Tags {
				to_save := trimQuotes(v)
				exif_setters = append(exif_setters, &models.PublicreportImageExifSetter{
					ImageID: omit.From(image.ID),
					Name:    omit.From(k),
					Value:   omit.From(to_save),
				})
			}
			if len(exif_setters) > 0 {
				_, err = models.PublicreportImageExifs.Insert(bob.ToMods(exif_setters...)).Exec(ctx, tx)
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
