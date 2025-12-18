package userfile

import (
	"fmt"
	"io"
	"log"
	"os"

	"github.com/google/uuid"
)

var UserFilesDirectory string

func AudioFileContentPathRaw(audioUUID string) string {
	return fmt.Sprintf("%s/%s.m4a", UserFilesDirectory, audioUUID)
}
func AudioFileContentPathMp3(audioUUID string) string {
	return fmt.Sprintf("%s/%s.mp3", UserFilesDirectory, audioUUID)
}
func AudioFileContentPathNormalized(audioUUID string) string {
	return fmt.Sprintf("%s/%s-normalized.m4a", UserFilesDirectory, audioUUID)
}
func AudioFileContentPathOgg(audioUUID string) string {
	return fmt.Sprintf("%s/%s.ogg", UserFilesDirectory, audioUUID)
}
func AudioFileContentWrite(audioUUID uuid.UUID, body io.Reader) error {
	// Create file in configured directory
	filepath := AudioFileContentPathRaw(audioUUID.String())
	dst, err := os.Create(filepath)
	if err != nil {
		log.Printf("Failed to create audio file at %s: %v\n", filepath, err)
		return fmt.Errorf("Failed to create audio file at %s: %v", filepath, err)
	}
	defer dst.Close()

	// Copy rest of request body to file
	_, err = io.Copy(dst, body)
	if err != nil {
		return fmt.Errorf("Unable to save file to create audio file at %s: %v", filepath, err)
	}
	log.Printf("Saved audio content to %s\n", filepath)
	return nil
}
func ImageFileContentPathRaw(uid string) string {
	return fmt.Sprintf("%s/%s.raw", UserFilesDirectory, uid)
}
func ImageFileContentWrite(uid uuid.UUID, body io.Reader) error {
	filepath := ImageFileContentPathRaw(uid.String())

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
