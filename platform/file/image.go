package file

import (
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
)

func ImageFileFromReader(collection Collection, uid uuid.UUID, body io.Reader) error {
	filepath := fileContentPathUUID(collection, uid)

	// Create file in configured directory
	dst, err := os.Create(filepath)
	if err != nil {
		return fmt.Errorf("Failed to create image file %s: %w", filepath, err)
	}
	defer dst.Close()

	// Copy rest of request body to file
	_, err = io.Copy(dst, body)
	if err != nil {
		return fmt.Errorf("Unable to save file %s: %w", filepath, err)
	}
	log.Info().Str("filepath", filepath).Int("collection", int(collection)).Msg("Saved image file content to collection")
	return nil
}
func ImageFileToWriter(collection Collection, uid uuid.UUID, w http.ResponseWriter) {
	image_path := fileContentPathUUID(collection, uid)
	writeFileContent(w, image_path)
}
