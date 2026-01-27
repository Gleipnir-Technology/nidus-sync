package llm

import (
	"context"
	"fmt"
	"strings"

	"github.com/maruel/genai"
	"github.com/maruel/genai/adapters"
	"github.com/maruel/genai/providers/openaichat"
	"github.com/rs/zerolog/log"
)

func CreateOpenAIClient(ctx context.Context) error {
	logger := createLogger()

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

type QueryReportStatusInput struct {
	ReportID string `json:"report_id"`
}

var client *openAIClient

func (c *openAIClient) continueConversation(ctx context.Context, history []Message, customer_phone string) (Message, error) {
	opts := genai.OptionsTools{
		Tools: []genai.ToolDef{
			{
				Name:        "query_report_status",
				Description: "This is used to answer any questions about the current state of the mosquito nuisance report.",
				Callback: func(ctx2 context.Context, input *QueryReportStatusInput) (string, error) {
					return c.queryReportStatus(ctx2, customer_phone)
				},
			},
		},
	}

	msg := c.convertHistory(history)
	res, _, err := adapters.GenSyncWithToolCallLoop(ctx, c.client, genai.Messages{msg}, &opts)
	if err != nil {
		return Message{}, fmt.Errorf("Failed to continue conversation: %v", err)
	}

	for _, m := range res {
		// Empty responses are tool call related.
		if m.String() == "" {
			log.Debug().Msg("Tool called")
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

func (c *openAIClient) convertHistory(history []Message) genai.Message {
	var sb strings.Builder
	sb.WriteString(
		`This is a text chat conversation between a customer that's a member of the public and a mosquito abatement district.
		The customer has reported a mosquito nuisance or mosquito breeding through the website report.mosquitoes.online.
		Messages from the customer are prefixed with 'customer:' and reponses from the service agent servicing the request are prefixed with 'agent:'.
		The agent wants to provide clear, confident, and succint information about the state of the customer's request. The agent also provides general information about how members of the public can help with controlling mosquitoes. For complex or highly specific requests, the agent will need to defer to the mosquito abatement district. This will take some time because contacting the district may take a few hours to get a response. When the agent needs to contact the district, the agent should tell the customer they are reaching out to the district and to expect a delay.
		Transcript starts:`,
	)
	for _, h := range history {
		if h.IsFromCustomer {
			sb.WriteString(fmt.Sprintf("\n\ncustomer (%s): %s\n", h.Content))
		} else {
			sb.WriteString(fmt.Sprintf("\n\nagent (%s): %s\n", h.Content))
		}
	}
	return genai.NewTextMessage(sb.String())
}

func (c *openAIClient) queryReportStatus(ctx context.Context, customer_phone string) (string, error) {
	return "Report is scheduled for work in 3 days at 2:00pm by the district", nil
}
