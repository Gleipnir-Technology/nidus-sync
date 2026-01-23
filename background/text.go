package background

import (
	"context"
	"errors"
	"fmt"

	"github.com/Gleipnir-Technology/nidus-sync/comms"
	"github.com/Gleipnir-Technology/nidus-sync/config"
	"github.com/Gleipnir-Technology/nidus-sync/db/enums"
	"github.com/rs/zerolog/log"
)

var channelJobText chan jobText

func ReportSubscriptionConfirmationText(destination comms.E164, report_id string) {
	enqueueJobText(jobText{
		Destination: destination,
		ReportID:    report_id,
		Source:      config.RMOPhoneNumber,
		Type:        enums.CommsMessagetypetextReportSubscriptionConfirmation,
	})
}

type jobText struct {
	Destination comms.E164
	ReportID    string
	Source      comms.E164
	Type        enums.CommsMessagetypetext
}

func enqueueJobText(job jobText) {
	select {
	case channelJobText <- job:
		log.Info().Msg("Enqueued text job")
	default:
		log.Warn().Msg("sms job channel is full, dropping job")
	}
}

func startWorkerText(ctx context.Context, channel chan jobText) {
	go func() {
		for {
			select {
			case <-ctx.Done():
				log.Info().Msg("Email worker shutting down.")
				return
			case job := <-channel:
				err := jobProcessText(job)
				if err != nil {
					log.Error().Err(err).Str("type", string(job.Type)).Msg("Error processing text message job")
				}
			}
		}
	}()
}

func jobProcessText(job jobText) error {
	var message string
	switch job.Type {
	case enums.CommsMessagetypetextInitialContact:
		message = "This is Report Mosquitoes Online. We just got your number. Text \"YES\" to get texts, or \"STOP\" to stap."
	case enums.CommsMessagetypetextReportSubscriptionConfirmation:
		message = "Thanks for submitting a mosquito report. Text for any questions. We'll send you updates as we get them."
	default:
		return errors.New("No idea what message to send")
	}
	err := comms.SendText(job.Source, job.Destination, message)
	if err != nil {
		log.Error().Err(err).Msg("Failed to send text message")
		return fmt.Errorf("Failed to send message '%s' to '%s'", job.Type, job.Destination)
	}
	return nil
}
