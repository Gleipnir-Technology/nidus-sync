package llm

import (
	"context"
	"fmt"
	"strings"

	"github.com/maruel/genai"
	"github.com/maruel/genai/adapters"
	"github.com/maruel/genai/providers/openaichat"
	"github.com/rs/zerolog"
)

func CreateOpenAIClient(ctx context.Context, logger *zerolog.Logger) error {
	linkLogger(logger)

	opts := genai.ProviderOptions{
		Model: genai.ModelCheap,
	}
	c, err := openaichat.New(ctx, &opts, nil)
	if err != nil {
		return fmt.Errorf("Failed to create genai client: %v", err)
	}
	client = &openAIClient{
		client:        c,
		conversations: make(map[string][]genai.Message),
		log:           logger,
	}
	return nil
}

type openAIClient struct {
	client        *openaichat.Client
	conversations map[string][]genai.Message
	log           *Logger
}

type ContactSupervisorInput struct {
	Reason string `json:"reason"`
}

type ContactDistrictInput struct {
	Reason string `json:"reason"`
}

type QueryReportStatusInput struct {
	ReportID string `json:"report_id"`
}

var client *openAIClient

func (c *openAIClient) continueConversation(ctx context.Context, tools genai.OptionsTools, msg genai.Message) (Message, error) {
	res, _, err := adapters.GenSyncWithToolCallLoop(ctx, c.client, genai.Messages{msg}, &tools)
	if err != nil {
		return Message{}, fmt.Errorf("Failed to continue conversation: %v", err)
	}

	for _, m := range res {
		// Empty responses are tool call related.
		if m.String() == "" {
			//log.Debug().Msg("Tool called")
		} else {
			var toSay string = m.String()
			toSay = strings.Replace(toSay, "report-mosquitoes-online: ", "", 1)
			return Message{
				Content:        toSay,
				IsFromCustomer: false,
			}, nil
		}
	}

	return Message{}, nil
}
