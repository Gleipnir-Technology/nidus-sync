package email

import (
	"context"
	"errors"
	"fmt"

	"github.com/Gleipnir-Technology/nidus-sync/db/enums"
	"github.com/rs/zerolog/log"
)

type Job interface {
	destination() string
	messageType() enums.CommsMessagetypeemail
	renderHTML() (string, error)
	renderTXT() (string, error)
	subject() string
}

type jobEmailBase struct {
	destination string
	source      string
}

func Handle(ctx context.Context, job Job) error {
	var err error
	switch job.messageType() {
	case enums.CommsMessagetypeemailReportSubscriptionConfirmation:
		err = sendEmailReportConfirmation(ctx, job)
	default:
		return errors.New("not implemented")
	}
	if err != nil {
		log.Error().Err(err).Str("dest", job.destination()).Str("type", string(job.messageType())).Msg("Error processing email")
		return fmt.Errorf("Failed to handle email: %w", err)
	}
	return nil
	/*
		case enums.CommsMessagetypeemailReportStatusScheduled:
		case enums.CommsMessagetypeemailReportStatusComplete:

		}
	*/
}
