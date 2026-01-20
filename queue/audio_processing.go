package queue

import (
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
)

// AudioJob represents a job to process an audio file.
type JobAudio struct {
	AudioUUID uuid.UUID
}

var ChannelJobAudio chan JobAudio

// EnqueueAudioJob sends an audio processing job to the worker.
func EnqueueAudioJob(job JobAudio) {
	select {
	case ChannelJobAudio <- job:
		log.Info().Str("uuid", job.AudioUUID.String()).Msg("Enqueued audio job")
	default:
		log.Warn().Str("uuid", job.AudioUUID.String()).Msg("Audio job channel is full, dropping job")
	}
}
