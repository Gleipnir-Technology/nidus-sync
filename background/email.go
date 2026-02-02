package background

import (
	"context"

	"github.com/Gleipnir-Technology/nidus-sync/comms/email"
	"github.com/rs/zerolog/log"
)

var channelJobEmail chan email.Job

func ReportSubscriptionConfirmationEmail(destination, report_id string) {
	enqueueJobEmail(email.NewJobReportNotificationConfirmation(
		destination,
		report_id,
	))
}

func enqueueJobEmail(job email.Job) {
	select {
	case channelJobEmail <- job:
		return
	default:
		log.Warn().Msg("email job channel is full, dropping job")
	}
}

func startWorkerEmail(ctx context.Context, channel chan email.Job) {
	go func() {
		log.Info().Msg("Email worker started")
		for {
			select {
			case <-ctx.Done():
				log.Info().Msg("Email worker shutting down.")
				return
			case job := <-channel:
				err := email.Handle(ctx, job)
				if err != nil {
					log.Error().Err(err).Msg("Failed to handle email message")
				}
			}
		}
	}()
}
