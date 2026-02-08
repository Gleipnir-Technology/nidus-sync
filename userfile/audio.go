package userfile

import (
	"errors"
	"fmt"
	"os"
	"os/exec"

	"github.com/Gleipnir-Technology/nidus-sync/db"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
)

func NormalizeAudio(audioUUID uuid.UUID) error {
	//source := AudioFileContentPathRaw(audioUUID.String())
	source := fileContentPath("user", audioUUID, "m4a")
	_, err := os.Stat(source)
	if errors.Is(err, os.ErrNotExist) {
		log.Warn().Str("source", source).Msg("file doesn't exist, skipping normalization")
		return nil
	}
	log.Info().Str("sourcce", source).Msg("Normalizing")
	//destination := AudioFileContentPathNormalized(audioUUID.String())
	destination := fileContentPathAudioNormalized(audioUUID)
	// Use "ffmpeg" directly, assuming it's in the system PATH
	cmd := exec.Command("ffmpeg", "-i", source, "-filter:a", "loudnorm", destination)
	out, err := cmd.CombinedOutput()
	if err != nil {
		log.Printf("FFmpeg output for normalization: %s", out)
		return fmt.Errorf("ffmpeg normalization failed: %v", err)
	}
	err = db.NoteAudioNormalized(audioUUID.String())
	if err != nil {
		return fmt.Errorf("failed to update database for normalized audio %s: %v", audioUUID, err)
	}
	log.Info().Str("destination", destination).Msg("Normalized audio")
	return nil
}

func TranscodeToOgg(audioUUID uuid.UUID) error {
	//source := AudioFileContentPathNormalized(audioUUID.String())
	source := fileContentPathAudioNormalized(audioUUID)
	_, err := os.Stat(source)
	if errors.Is(err, os.ErrNotExist) {
		log.Warn().Str("source", source).Msg("file doesn't exist, skipping OGG transcoding")
		return nil
	}
	log.Info().Str("source", source).Msg("Transcoding to ogg")
	//destination := userfile.AudioFileContentPathOgg(audioUUID.String())
	destination := fileContentPath("user", audioUUID, "ogg")
	// Use "ffmpeg" directly, assuming it's in the system PATH
	cmd := exec.Command("ffmpeg", "-i", source, "-vn", "-acodec", "libvorbis", destination)
	out, err := cmd.CombinedOutput()
	if err != nil {
		log.Error().Err(err).Bytes("out", out).Msg("FFmpeg output for OGG transcoding")
		return fmt.Errorf("ffmpeg OGG transcoding failed: %v", err)
	}
	err = db.NoteAudioTranscodedToOgg(audioUUID.String())
	if err != nil {
		return fmt.Errorf("failed to update database for OGG transcoded audio %s: %v", audioUUID, err)
	}
	log.Info().Str("destination", destination).Msg("Transcoded audio")
	return nil
}

func fileContentPathAudioNormalized(u uuid.UUID) string {
	//destination := AudioFileContentPathNormalized(audioUUID.String())
	return fileContentPath("user", u, "normalized.m4a")
}
