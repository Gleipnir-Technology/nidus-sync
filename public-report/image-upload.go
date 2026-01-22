package publicreport

import (
	"bytes"
	"context"
	"fmt"
	"image"
	_ "image/gif"  // register GIF format
	_ "image/jpeg" // register JPEG format
	_ "image/png"  // register PNG format
	"io"
	"mime/multipart"
	"net/http"
	"time"

	"github.com/Gleipnir-Technology/nidus-sync/db/models"
	"github.com/Gleipnir-Technology/nidus-sync/userfile"
	"github.com/aarondl/opt/omit"
	"github.com/aarondl/opt/omitnull"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"github.com/rwcarlsen/goexif/exif"
	"github.com/stephenafamo/bob"
	"github.com/stephenafamo/bob/dialect/psql"
	"github.com/stephenafamo/bob/dialect/psql/um"
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
	Exif        ExifCollection
	Format      string

	UploadFilesize int
	UploadFilename string
	UUID           uuid.UUID
}

func extractExif(content_type string, file_bytes []byte) (result ExifCollection, err error) {
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

	exif, err := exif.Decode(bytes.NewReader(file_bytes))
	if err != nil {
		return result, fmt.Errorf("Failed to decode image meta: %w", err)
	}
	lat, lng, _ := exif.LatLong()
	result.GPS = &GPS{
		Latitude:  lat,
		Longitude: lng,
	}
	return result, err
}

func extractImageUpload(headers *multipart.FileHeader) (upload ImageUpload, err error) {
	file, err := headers.Open()
	if err != nil {
		return upload, fmt.Errorf("Failed to open header: %w", err)
	}
	defer file.Close()

	file_bytes, err := io.ReadAll(file)
	content_type := http.DetectContentType(file_bytes)

	exif, err := extractExif(content_type, file_bytes)
	if err != nil {
		//return upload, fmt.Errorf("Failed to extract EXIF data: %w", err)
		log.Warn().Err(err).Msg("Failed to extract EXIF data")
	}
	log.Debug().Float64("lat", exif.GPS.Latitude).Float64("lng", exif.GPS.Longitude).Msg("extracted GPS from exif")

	i, format, err := image.Decode(bytes.NewReader(file_bytes))
	if err != nil {
		return upload, fmt.Errorf("Failed to decode image file: %w", err)
	}

	u, err := uuid.NewUUID()
	if err != nil {
		return upload, fmt.Errorf("Failed to create quick report photo uuid", err)
	}
	err = userfile.PublicImageFileContentWrite(u, bytes.NewReader(file_bytes))
	if err != nil {
		return upload, fmt.Errorf("Failed to write image file to disk: %w", err)
	}
	log.Info().Int("size", len(file_bytes)).Str("uploaded_filename", headers.Filename).Str("content-type", content_type).Str("uuid", u.String()).Msg("Saved an uploaded file to disk")
	return ImageUpload{
		Bounds:         i.Bounds(),
		ContentType:    content_type,
		Exif:           exif,
		Format:         format,
		UploadFilename: headers.Filename,
		UploadFilesize: len(file_bytes),
		UUID:           u,
	}, nil
}

func extractImageUploads(r *http.Request) (uploads []ImageUpload, err error) {
	uploads = make([]ImageUpload, 0)
	for _, fheaders := range r.MultipartForm.File {
		for _, headers := range fheaders {
			upload, err := extractImageUpload(headers)
			if err != nil {
				return make([]ImageUpload, 0), fmt.Errorf("Failed to extract photo upload: %w", err)
			}
			uploads = append(uploads, upload)
		}
	}
	return uploads, nil
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
		if u.Exif.GPS != nil {
			_, err = psql.Update(
				um.Table("publicreport.image"),
				um.SetCol("location").To(fmt.Sprintf("ST_GeometryFromText('Point(%f %f)')", u.Exif.GPS.Longitude, u.Exif.GPS.Latitude)),
				um.Where(psql.Quote("id").EQ(psql.Arg(image.ID))),
			).Exec(ctx, tx)
		}

		exif_setters := make([]*models.PublicreportImageExifSetter, 0)
		for k, v := range u.Exif.Tags {
			exif_setters = append(exif_setters, &models.PublicreportImageExifSetter{
				ImageID: omit.From(image.ID),
				Name:    omit.From(k),
				Value:   omit.From(v),
			})
		}
		if len(exif_setters) > 0 {
			_, err = models.PublicreportImageExifs.Insert(bob.ToMods(exif_setters...)).Exec(ctx, tx)
			if err != nil {
				return images, fmt.Errorf("Failed to create photo exif records: %w", err)
			}
		}
		images = append(images, image)
		log.Info().Int32("id", image.ID).Int("tags", len(u.Exif.Tags)).Msg("Saved an uploaded file to the database")
	}
	return images, nil
}
