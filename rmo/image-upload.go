package rmo

import (
	"bytes"
	"fmt"
	"github.com/Gleipnir-Technology/nidus-sync/platform"
	"github.com/Gleipnir-Technology/nidus-sync/platform/file"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"image"
	_ "image/gif"  // register GIF format
	_ "image/jpeg" // register JPEG format
	_ "image/png"  // register PNG format
	"io"
	"mime/multipart"
	"net/http"
)

func extractImageUpload(headers *multipart.FileHeader) (upload platform.ImageUpload, err error) {
	f, err := headers.Open()
	if err != nil {
		return upload, fmt.Errorf("Failed to open header: %w", err)
	}
	defer f.Close()

	file_bytes, err := io.ReadAll(f)
	content_type := http.DetectContentType(file_bytes)

	exif, err := platform.ImageExtractExif(content_type, file_bytes)
	if err != nil {
		return upload, fmt.Errorf("Failed to extract EXIF data: %w", err)
	}
	i, format, err := image.Decode(bytes.NewReader(file_bytes))
	if err != nil {
		return upload, fmt.Errorf("Failed to decode image file: %w", err)
	}

	u, err := uuid.NewUUID()
	if err != nil {
		return upload, fmt.Errorf("Failed to create quick report photo uuid", err)
	}
	err = file.PublicImageFileContentWrite(u, bytes.NewReader(file_bytes))
	if err != nil {
		return upload, fmt.Errorf("Failed to write image file to disk: %w", err)
	}
	log.Info().Int("size", len(file_bytes)).Str("uploaded_filename", headers.Filename).Str("content-type", content_type).Str("uuid", u.String()).Msg("Saved an uploaded file to disk")
	return platform.ImageUpload{
		Bounds:         i.Bounds(),
		ContentType:    content_type,
		Exif:           exif,
		Format:         format,
		UploadFilename: headers.Filename,
		UploadFilesize: len(file_bytes),
		UUID:           u,
	}, nil
}

func extractImageUploads(r *http.Request) (uploads []platform.ImageUpload, err error) {
	uploads = make([]platform.ImageUpload, 0)
	for _, fheaders := range r.MultipartForm.File {
		for _, headers := range fheaders {
			upload, err := extractImageUpload(headers)
			if err != nil {
				return make([]platform.ImageUpload, 0), fmt.Errorf("Failed to extract photo upload: %w", err)
			}
			uploads = append(uploads, upload)
		}
	}
	return uploads, nil
}
