package background

import (
	"context"

	"github.com/Gleipnir-Technology/nidus-sync/queue"
	"github.com/rs/zerolog/log"
)

func StartWorkerEmail(ctx context.Context, channel chan queue.JobEmail) {
	go func() {
		for {
			select {
			case <-ctx.Done():
				log.Info().Msg("Email worker shutting down.")
				return
			case job := <-channel:
				err := processJobEmail(job)
				if err != nil {
					log.Error().Err(err).Str("dest", job.Destination).Str("type", string(job.Type)).Msg("Error processing audio file")
				}
			}
		}
	}()
}

func StartWorkerSMS(ctx context.Context, channel chan queue.JobSMS) {
	go func() {
		for {
			select {
			case <-ctx.Done():
				log.Info().Msg("Email worker shutting down.")
				return
			case job := <-channel:
				err := processJobSMS(job)
				if err != nil {
					log.Error().Err(err).Str("dest", job.Destination).Str("type", string(job.Type)).Msg("Error processing audio file")
				}
			}
		}
	}()
}

func processJobEmail(job queue.JobEmail) error {
	log.Info().Str("dest", job.Destination).Str("type", string(job.Type)).Msg("Pretend doing email job")
	return nil
}

func processJobSMS(job queue.JobSMS) error {
	log.Info().Str("dest", job.Destination).Str("type", string(job.Type)).Msg("Pretend doing email job")
	return nil
}
