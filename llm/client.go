package llm

import (
	"context"
	"fmt"
	"strings"

	"github.com/maruel/genai"
	"github.com/rs/zerolog/log"
)

type Message struct {
	Content        string
	IsFromCustomer bool
}

func GenerateNextMessage(ctx context.Context, history []Message, _handle_report_status func() (string, error), _handle_contact_district func(string), _handle_contact_supervisor func(string)) (Message, error) {
	msg := convertHistory(history)
	tools := genai.OptionsTools{
		Tools: []genai.ToolDef{
			{
				Name:        "contact_district",
				Description: "Reach out to the district to get answers for a customer about their operations or schedule.",
				Callback: func(ctx2 context.Context, input *ContactDistrictInput) (string, error) {
					_handle_contact_district(input.Reason)
					return "district has been contacted.", nil
				},
			}, {
				Name:        "contact_supervisor",
				Description: "Flag a conversation from a customer as abusive, concerning, or off-topic.",
				Callback: func(ctx2 context.Context, input *ContactSupervisorInput) (string, error) {
					_handle_contact_supervisor(input.Reason)
					return "supervisor has been notified", nil
				},
			}, {
				Name:        "query_report_status",
				Description: "This is used to answer any questions about the current state of the mosquito nuisance report.",
				Callback: func(ctx2 context.Context, input *QueryReportStatusInput) (string, error) {
					return _handle_report_status()
				},
			},
		},
	}
	next, err := client.continueConversation(ctx, tools, msg)
	if err != nil {
		return Message{}, fmt.Errorf("Failed to generate next message: %w", err)
	}
	trimmed, found := strings.CutPrefix(next.Content, "agent:")
	if !found {
		trimmed, found = strings.CutPrefix(next.Content, "Agent:")
		if !found {
			log.Warn().Str("content", next.Content).Msg("No 'agent:' prefix on next message")
		}
	}
	next.Content = trimmed

	return next, nil
}
func convertHistory(history []Message) genai.Message {
	var sb strings.Builder
	sb.WriteString(
		`This is a text chat conversation between a customer that's a member of the public and a mosquito abatement district.
		The customer has reported a mosquito nuisance or mosquito breeding through the website report.mosquitoes.online.
		Messages from the customer are prefixed with 'customer:' and reponses from the service agent servicing the request are prefixed with 'agent:'.
		The agent provides clear, confident, and succint information about the state of the customer's request.
		The agent answers just the questions that are asked, and prefers very short answers because the conversation is happening over SMS.
		The agent rarely asks questions, preferring to just answer direct queries.
		For complex or highly specific requests, the agent will need to defer to the mosquito abatement district. This will take some time because contacting the district may take a few hours to get a response. When the agent needs to contact the district, the agent should tell the customer they are reaching out to the district and to expect a delay.
		When conversations start to veer away from the agent's job they should contact a supervisor.
		Transcript starts:`,
	)
	for _, h := range history {
		if h.IsFromCustomer {
			sb.WriteString(fmt.Sprintf("\n\ncustomer: %s\n", h.Content))
		} else {
			sb.WriteString(fmt.Sprintf("\n\nagent: %s\n", h.Content))
		}
	}
	return genai.NewTextMessage(sb.String())
}
