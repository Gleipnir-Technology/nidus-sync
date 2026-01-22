package userfile

import (
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/Gleipnir-Technology/nidus-sync/config"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
)

func AudioFileContentPathRaw(audioUUID string) string {
	return fmt.Sprintf("%s/%s.m4a", config.FilesDirectoryUser, audioUUID)
}
func AudioFileContentPathMp3(audioUUID string) string {
	return fmt.Sprintf("%s/%s.mp3", config.FilesDirectoryUser, audioUUID)
}
func AudioFileContentPathNormalized(audioUUID string) string {
	return fmt.Sprintf("%s/%s-normalized.m4a", config.FilesDirectoryUser, audioUUID)
}
func AudioFileContentPathOgg(audioUUID string) string {
	return fmt.Sprintf("%s/%s.ogg", config.FilesDirectoryUser, audioUUID)
}
func AudioFileContentWrite(audioUUID uuid.UUID, body io.Reader) error {
	// Create file in configured directory
	filepath := AudioFileContentPathRaw(audioUUID.String())
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
func ImageFileContentPathRawUser(uid string) string {
	return imageFileContentPath(config.FilesDirectoryUser, uid, "raw")
}
func imageFileContentPathLogoPng(uid string) string {
	return imageFileContentPath(config.FilesDirectoryLogo, uid, "png")
}
func imageFileContentPath(dir string, uid string, ext string) string {
	return fmt.Sprintf("%s/%s.%s", dir, uid, ext)
}
func ImageFileContentWrite(uid uuid.UUID, body io.Reader) error {
	filepath := ImageFileContentPathRawUser(uid.String())

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
	image_path := imageFileContentPathLogoPng(uid.String())
	writeFileContent(w, image_path)
}

func PublicImageFileContentWrite(uid uuid.UUID, body io.Reader) error {
	// Create file in configured directory
	filepath := PublicImageFileContentPathRaw(uid.String())
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

func PublicImageFileContentPathRaw(uid string) string {
	return fmt.Sprintf("%s/%s.raw", config.FilesDirectoryPublic, uid)
}

func PublicImageFileToResponse(w http.ResponseWriter, uid string) {
	image_path := PublicImageFileContentPathRaw(uid)
	writeFileContent(w, image_path)
}

func writeFileContent(w http.ResponseWriter, image_path string) {
	// Open the file
	file, err := os.Open(image_path)
	if err != nil {
		if os.IsNotExist(err) {
			http.Error(w, "Image not found", http.StatusNotFound)
		} else {
			http.Error(w, "Failed to retrieve image", http.StatusInternalServerError)
		}
		return
	}
	defer file.Close()

	// Get file info for Content-Length header
	fileInfo, err := file.Stat()
	if err != nil {
		http.Error(w, "Failed to get image information", http.StatusInternalServerError)
		return
	}

	// Set appropriate headers
	w.Header().Set("Content-Type", "image/png")
	w.Header().Set("Content-Length", fmt.Sprintf("%d", fileInfo.Size()))

	// Copy file contents to response writer
	_, err = io.Copy(w, file)
	if err != nil {
		// Note: At this point, we've already started writing the response,
		// so we can't change the status code anymore. The best we can do
		// is log the error and abandon the connection.
		return
	}
}
