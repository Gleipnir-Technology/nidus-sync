package publicreport

import (
	"bytes"
	"fmt"
	"net/http"

	"github.com/Gleipnir-Technology/nidus-sync/userfile"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
)

type PhotoUpload struct {
	Filename string
	Size     int64
	UUID     uuid.UUID
}

func extractPhotoUploads(r *http.Request) (uploads []PhotoUpload, err error) {
	uploads = make([]PhotoUpload, 0)
	for _, fheaders := range r.MultipartForm.File {
		for _, headers := range fheaders {
			file, err := headers.Open()

			if err != nil {
				return uploads, fmt.Errorf("Failed to open header: %v", err)
			}

			defer file.Close()

			buff := make([]byte, 512)
			file.Read(buff)

			file.Seek(0, 0)
			contentType := http.DetectContentType(buff)
			var sizeBuff bytes.Buffer
			fileSize, err := sizeBuff.ReadFrom(file)
			if err != nil {
				return uploads, fmt.Errorf("Failed to read file: %v", err)
			}
			file.Seek(0, 0)
			log.Info().Int64("size", fileSize).Str("filename", headers.Filename).Str("content-type", contentType).Msg("Got an uploaded file")
			u, err := uuid.NewUUID()
			if err != nil {
				return uploads, fmt.Errorf("Failed to create quick report photo uuid", err)
			}
			err = userfile.PublicImageFileContentWrite(u, file)
			if err != nil {
				return uploads, fmt.Errorf("Failed to write image file to disk: %v", err)
			}
			uploads = append(uploads, PhotoUpload{
				Size:     fileSize,
				Filename: headers.Filename,
				UUID:     u,
			})
		}
	}
	return uploads, nil
}
