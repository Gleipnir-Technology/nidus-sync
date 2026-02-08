package background

import (
	"context"
	"fmt"

	"github.com/Gleipnir-Technology/nidus-sync/userfile"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
)

// AudioJob represents a job to process an audio file.
type jobAudio struct {
	AudioUUID uuid.UUID
}

var channelJobAudio chan jobAudio

func AudioTranscode(audio_uuid uuid.UUID) {
	enqueueAudioJob(jobAudio{
		AudioUUID: audio_uuid,
	})
}

// startAudioWorker initializes the audio job channel and starts the worker goroutine.
func startWorkerAudio(ctx context.Context, audioJobChannel chan jobAudio) {
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

// EnqueueAudioJob sends an audio processing job to the worker.
func enqueueAudioJob(job jobAudio) {
	select {
	case channelJobAudio <- job:
		log.Info().Str("uuid", job.AudioUUID.String()).Msg("Enqueued audio job")
	default:
		log.Warn().Str("uuid", job.AudioUUID.String()).Msg("Audio job channel is full, dropping job")
	}
}

func processAudioFile(audioUUID uuid.UUID) error {
	// Normalize audio
	err := userfile.NormalizeAudio(audioUUID)
	if err != nil {
		return fmt.Errorf("failed to normalize audio %s: %v", audioUUID, err)
	}

	// Transcode to OGG
	err = userfile.TranscodeToOgg(audioUUID)
	if err != nil {
		return fmt.Errorf("failed to transcode audio %s to OGG: %v", audioUUID, err)
	}

	enqueueLabelStudioJob(jobLabelStudio{
		UUID: audioUUID,
	})
	return nil
}
