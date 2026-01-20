package queue

import (
	"github.com/Gleipnir-Technology/nidus-sync/db/enums"
	"github.com/rs/zerolog/log"
)

type JobEmail struct {
	Destination string
	Source      string
	Type        enums.CommsEmailmessagetype
}
type JobSMS struct {
	Destination string
	Source      string
	Type        enums.CommsSmsmessagetype
}

var ChannelJobEmail chan JobEmail
var ChannelJobSMS chan JobSMS

func EnqueueJobEmail(job JobEmail) {
	select {
	case ChannelJobEmail <- job:
		log.Info().Str("destination", job.Destination).Msg("Enqueued email job")
	default:
		log.Warn().Msg("email job channel is full, dropping job")
	}
}

func EnqueueJobSMS(job JobSMS) {
	select {
	case ChannelJobSMS <- job:
		log.Info().Str("destination", job.Destination).Msg("Enqueued sms job")
	default:
		log.Warn().Msg("sms job channel is full, dropping job")
	}
}
