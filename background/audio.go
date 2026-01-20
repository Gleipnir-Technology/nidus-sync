package background

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"

	"github.com/Gleipnir-Technology/nidus-sync/db"
	"github.com/Gleipnir-Technology/nidus-sync/queue"
	"github.com/Gleipnir-Technology/nidus-sync/userfile"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
)

// StartAudioWorker initializes the audio job channel and starts the worker goroutine.
func StartWorkerAudio(ctx context.Context, audioJobChannel chan queue.JobAudio) {
	go func() {
		for {
			select {
			case <-ctx.Done():
				log.Info().Msg("Audio worker shutting down.")
				return
			case job := <-audioJobChannel:
				log.Info().Str("uuid", job.AudioUUID.String()).Msg("Processing audio job")
				err := processAudioFile(job.AudioUUID)
				if err != nil {
					log.Error().Err(err).Str("uuid", job.AudioUUID.String()).Msg("Error processing audio file")
				}
			}
		}
	}()
}

func normalizeAudio(audioUUID uuid.UUID) error {
	source := userfile.AudioFileContentPathRaw(audioUUID.String())
	_, err := os.Stat(source)
	if errors.Is(err, os.ErrNotExist) {
		log.Warn().Str("source", source).Msg("file doesn't exist, skipping normalization")
		return nil
	}
	log.Info().Str("sourcce", source).Msg("Normalizing")
	destination := userfile.AudioFileContentPathNormalized(audioUUID.String())
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

func processAudioFile(audioUUID uuid.UUID) error {
	// Normalize audio
	err := normalizeAudio(audioUUID)
	if err != nil {
		return fmt.Errorf("failed to normalize audio %s: %v", audioUUID, err)
	}

	// Transcode to OGG
	err = transcodeToOgg(audioUUID)
	if err != nil {
		return fmt.Errorf("failed to transcode audio %s to OGG: %v", audioUUID, err)
	}

	queue.EnqueueLabelStudioJob(queue.LabelStudioJob{
		UUID: audioUUID,
	})
	return nil
}

func transcodeToOgg(audioUUID uuid.UUID) error {
	source := userfile.AudioFileContentPathNormalized(audioUUID.String())
	_, err := os.Stat(source)
	if errors.Is(err, os.ErrNotExist) {
		log.Warn().Str("source", source).Msg("file doesn't exist, skipping OGG transcoding")
		return nil
	}
	log.Info().Str("source", source).Msg("Transcoding to ogg")
	destination := userfile.AudioFileContentPathOgg(audioUUID.String())
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
