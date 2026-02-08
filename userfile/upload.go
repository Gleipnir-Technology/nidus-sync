package userfile

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
)

type FileUpload struct {
	ContentType string

	UploadFilesize int
	UploadFilename string
	UUID           uuid.UUID
}

func SaveFileUpload(r *http.Request, name string, subdir string, extension string) ([]FileUpload, error) {
	results := make([]FileUpload, 0)
	for n, fheaders := range r.MultipartForm.File {
		log.Debug().Str("n", n).Msg("looking at header")
		if n != name {
			continue
		}
		for _, headers := range fheaders {
			f, err := saveFileUpload(headers, subdir, extension)
			if err != nil {
				return results, fmt.Errorf("Failed to extract photo upload: %w", err)
			}
			results = append(results, f)
		}
	}
	return results, nil
}
func saveFileUploads(r *http.Request, subdir string, extension string) ([]FileUpload, error) {
	results := make([]FileUpload, 0)
	for name, fheaders := range r.MultipartForm.File {
		for _, headers := range fheaders {
			upload, err := saveFileUpload(headers, subdir, extension)
			if err != nil {
				return results, fmt.Errorf("Failed to save upload '%s': %w", name, err)
			}
			results = append(results, upload)
		}
	}
	return results, nil
}
func saveFileUpload(headers *multipart.FileHeader, subdir string, extension string) (upload FileUpload, err error) {
	file, err := headers.Open()
	if err != nil {
		return upload, fmt.Errorf("Failed to open header: %w", err)
	}
	defer file.Close()

	file_bytes, err := io.ReadAll(file)
	content_type := http.DetectContentType(file_bytes)

	u, err := uuid.NewUUID()
	if err != nil {
		return upload, fmt.Errorf("Failed to create uuid", err)
	}
	err = fileContentWrite(bytes.NewReader(file_bytes), subdir, u, extension)
	if err != nil {
		return upload, fmt.Errorf("Failed to write file to disk: %w", err)
	}
	log.Info().Int("size", len(file_bytes)).Str("uploaded_filename", headers.Filename).Str("content-type", content_type).Str("uuid", u.String()).Msg("Saved an uploaded file to disk")
	return FileUpload{
		ContentType:    content_type,
		UploadFilename: headers.Filename,
		UploadFilesize: len(file_bytes),
		UUID:           u,
	}, nil
}
