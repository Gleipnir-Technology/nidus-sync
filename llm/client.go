package llm

import (
	"context"
	"fmt"
	//"github.com/rs/zerolog/log"
)

type Message struct {
	Content        string
	IsFromCustomer bool
}

func GenerateNextMessage(ctx context.Context, history []Message, customer_phone string) (Message, error) {
	next, err := client.continueConversation(ctx, history, customer_phone)
	if err != nil {
		return Message{}, fmt.Errorf("Failed to generate next message: %w", err)
	}
	return next, nil
}
