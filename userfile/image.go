package userfile

import (
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
)

func ImageFileContentWrite(uid uuid.UUID, body io.Reader) error {
	filepath := fileContentPath(CollectionImageRaw, uid)

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
	return nil
}
func ImageFileContentWriteLogo(w http.ResponseWriter, uid uuid.UUID) {
	//image_path := imageFileContentPathLogoPng(uid.String())
	image_path := fileContentPath(CollectionLogo, uid)
	writeFileContent(w, image_path)
}

func PublicImageFileContentWrite(uid uuid.UUID, body io.Reader) error {
	// Create file in configured directory
	//filepath := PublicImageFileContentPathRaw(uid.String())
	filepath := fileContentPath(CollectionPublicImage, uid)
	dst, err := os.Create(filepath)
	if err != nil {
		log.Error().Err(err).Str("filepath", filepath).Msg("Failed to create public image file")
		return fmt.Errorf("Failed to create public image file at %s: %v", filepath, err)
	}
	defer dst.Close()

	// Copy rest of request body to file
	_, err = io.Copy(dst, body)
	if err != nil {
		return fmt.Errorf("Unable to save file to create audio file at %s: %v", filepath, err)
	}
	log.Info().Str("filepath", filepath).Msg("Saved public report image file content")
	return nil
}

func PublicImageFileToResponse(w http.ResponseWriter, uid uuid.UUID) {
	//image_path := PublicImageFileContentPathRaw(uid)
	image_path := fileContentPath(CollectionPublicImage, uid)
	writeFileContent(w, image_path)
}
