package text

import (
	"context"

	"github.com/rs/zerolog/log"
)

type MessageType int

const (
	ReportSubscription MessageType = iota
)

type Job interface {
	content() string
	destination() string
	messageType() MessageType
	messageTypeName() string
	source() string
}

func Handle(ctx context.Context, job Job) {
	var err error
	switch job.messageType() {
	case ReportSubscription:
		err = sendReportSubscription(ctx, job)
	}
	if err != nil {
		log.Error().Err(err).Str("dest", job.destination()).Str("type", string(job.messageTypeName())).Msg("Error processing email")
		return
	}
	/*
		case enums.CommsMessagetypeemailReportStatusScheduled:
		case enums.CommsMessagetypeemailReportStatusComplete:

		}
	*/
}
