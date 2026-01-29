package api

import (
	"os"

	"github.com/rs/zerolog/log"
)

func debugSaveRequest(body []byte, err error, message string) {
	// TODO(eliribble): avoid using a single static filename and instead securely generate
	// this value
	if err != nil {
		log.Error().Err(err).Msg(message)
	}
	output, err := os.OpenFile("/tmp/request.body", os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		log.Info().Msg("Failed to open temp request.bady")
	}
	defer output.Close()
	output.Write(body)
	log.Info().Msg("Wrote request to /tmp/request.body")
}
