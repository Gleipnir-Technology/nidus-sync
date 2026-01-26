package llm

import (
	"github.com/rs/zerolog/log"
)

type Message struct {
	Content        string
	IsFromCustomer bool
}

func GenerateNextMessage(history []Message, current Message) (Message, error) {
	// In general our history
	for i, msg := range history {
		log.Info().Int("i", i).Bool("is_customer", msg.IsFromCustomer).Msg("History")
	}

	return Message{
		Content:        "hey there. :)",
		IsFromCustomer: false,
	}, nil
}
