package userfile

import (
	"fmt"
	"io"
	//"net/http"
	"os"

	//"github.com/Gleipnir-Technology/nidus-sync/config"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
)

func FileContentWrite(body io.Reader, collection Collection, uid uuid.UUID) error {
	// Create file in configured directory
	filepath := fileContentPath(collection, uid)
	dst, err := os.Create(filepath)
	if err != nil {
		log.Error().Err(err).Str("filepath", filepath).Msg("Failed to create audio file")
		return fmt.Errorf("Failed to create audio file at %s: %v", filepath, err)
	}
	defer dst.Close()

	// Copy rest of request body to file
	_, err = io.Copy(dst, body)
	if err != nil {
		return fmt.Errorf("Unable to save file to create audio file at %s: %v", filepath, err)
	}
	log.Info().Str("filepath", filepath).Msg("Save audio file content")
	return nil
}

func NewFileReader(collection Collection, uid uuid.UUID) (io.Reader, error) {
	path := fileContentPath(collection, uid)
	return os.Open(path)
}
