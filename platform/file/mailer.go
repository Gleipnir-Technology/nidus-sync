package file

import (
	"fmt"
	"io"
	"os"

	"github.com/rs/zerolog/log"
)

func MailerFromReader(public_id string, body io.Reader) error {
	filepath := MailerPath(public_id)

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
	log.Info().Str("filepath", filepath).Str("collection", collectionName(CollectionMailerPDF)).Msg("Saved image file content to collection")
	return nil
}
func MailerPath(public_id string) string {
	collection := CollectionMailerPDF
	return fileContentPath(collection, public_id)
}
