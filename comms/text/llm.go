package text

import (
	"github.com/rs/zerolog/log"
)

func SendTextFromLLM(content string) {
	log.Info().Str("content", content).Msg("Pretend I sent a message")
}
