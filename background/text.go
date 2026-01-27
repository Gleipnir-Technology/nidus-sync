package background

import (
	"context"

	"github.com/Gleipnir-Technology/nidus-sync/config"
	"github.com/Gleipnir-Technology/nidus-sync/platform/text"
	"github.com/rs/zerolog/log"
)

var channelJobText chan text.Job

func ReportSubscriptionConfirmationText(destination text.E164, report_id string) {
	enqueueJobText(text.NewJobReportSubscriptionConfirmation(
		destination,
		report_id,
		config.PhoneNumberReport,
	))
}

func enqueueJobText(job text.Job) {
	select {
	case channelJobText <- job:
		log.Info().Msg("Enqueued text job")
	default:
		log.Warn().Msg("sms job channel is full, dropping job")
	}
}

func startWorkerText(ctx context.Context, channel chan text.Job) {
	go func() {
		log.Info().Msg("Text worker started")
		for {
			select {
			case <-ctx.Done():
				log.Info().Msg("Text worker shutting down.")
				return
			case job := <-channel:
				text.Handle(ctx, job)
			}
		}
	}()
}
