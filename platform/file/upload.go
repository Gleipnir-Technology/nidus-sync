package file

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"

	"github.com/Gleipnir-Technology/nidus-sync/lint"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
)

type Upload struct {
	ContentType string
	Name        string
	SizeBytes   int
	UUID        uuid.UUID
}

func SaveFileUploads(r *http.Request, collection Collection) ([]Upload, error) {
	results := make([]Upload, 0)
	for n, fheaders := range r.MultipartForm.File {
		log.Debug().Str("n", n).Msg("looking at header")
		for _, headers := range fheaders {
			f, err := saveFileUpload(headers, collection)
			if err != nil {
				return results, fmt.Errorf("Failed to extract photo upload: %w", err)
			}
			results = append(results, f)
		}
	}
	return results, nil
}
func saveFileUpload(headers *multipart.FileHeader, collection Collection) (upload Upload, err error) {
	file, err := headers.Open()
	if err != nil {
		return upload, fmt.Errorf("Failed to open header: %w", err)
	}
	defer lint.LogOnErr(file.Close, "close file")

	file_bytes, err := io.ReadAll(file)
	content_type := http.DetectContentType(file_bytes)

	u, err := uuid.NewUUID()
	if err != nil {
		return upload, fmt.Errorf("Failed to create uuid", err)
	}
	err = FileContentWrite(bytes.NewReader(file_bytes), collection, u)
	if err != nil {
		return upload, fmt.Errorf("Failed to write file to disk: %w", err)
	}
	log.Info().Int("size", len(file_bytes)).Str("uploaded_filename", headers.Filename).Str("content-type", content_type).Str("uuid", u.String()).Msg("Saved an uploaded file to disk")
	return Upload{
		ContentType: content_type,
		Name:        headers.Filename,
		SizeBytes:   len(file_bytes),
		UUID:        u,
	}, nil
}
