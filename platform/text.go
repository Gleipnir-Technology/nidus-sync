package platform

import (
	"context"
	"fmt"
	"strings"

	"github.com/Gleipnir-Technology/nidus-sync/comms/text"
	"github.com/Gleipnir-Technology/nidus-sync/config"
	"github.com/Gleipnir-Technology/nidus-sync/db"
	"github.com/Gleipnir-Technology/nidus-sync/db/models"
	"github.com/Gleipnir-Technology/nidus-sync/db/sql"
	"github.com/Gleipnir-Technology/nidus-sync/llm"
	"github.com/rs/zerolog/log"
)

// Translate from Twilio's representation of a RCS message sender to our concept of a phone number
// From: rcs:dev_report_mosquitoes_online_dosrvwxm_agent
// To: +16235525879
func getDst(ctx context.Context, to string) (string, error) {

	if to == config.TwilioRCSSenderRMO {
		return config.PhoneNumberReportStr, nil
	}
	/*
		phone, err := models.FindCommsPhone(ctx, db.PGInstance.BobDB, to)
		if err != nil {
			return "", fmt.Errorf("Failed to search for dest phone %s: %w", to, err)
		}
		return phone.E164, nil
	*/
	return "", fmt.Errorf("Cannot match phone number to '%s'", to)
}

func loadPreviousMessages(ctx context.Context, dst, src string) ([]llm.Message, error) {
	messages, err := sql.TextsBySenders(dst, src).All(ctx, db.PGInstance.BobDB)
	results := make([]llm.Message, 0)
	if err != nil {
		return results, fmt.Errorf("Failed to get message history for %s and %s: %w", dst, src, err)
	}
	log.Info().Int("count", len(messages)).Str("src", src).Str("dst", dst).Msg("Found previous messages")
	for _, m := range messages {
		is_from_customer := (m.Source == src)
		results = append(results, llm.Message{
			IsFromCustomer: is_from_customer,
			Content:        m.Content,
		})
	}
	return results, nil
}

func splitPhoneSource(s string) (string, string) {
	parts := strings.Split(s, ":")
	switch len(parts) {
	case 0:
		return "this isn't", "possible"
	case 1:
		return "", s
	case 2:
		return parts[0], parts[1]
	default:
		log.Warn().Str("s", s).Msg("Got an incomprehensible number of parts of a phone number")
		return parts[0], parts[1]
	}

}

func isSubscribed(ctx context.Context, src string) (bool, error) {
	phone, err := models.FindCommsPhone(ctx, db.PGInstance.BobDB, src)
	if err != nil {
		return false, fmt.Errorf("Failed to determine if '%s' is subscribed: %w", src, err)
	}
	return phone.IsSubscribed, nil
}

func HandleTextMessage(from string, to string, body string) {
	ctx := context.Background()
	type_, src := splitPhoneSource(from)
	dst, err := getDst(ctx, to)
	if err != nil {
		log.Error().Err(err).Str("to", to).Msg("Failed to get dst")
		return
	}
	subscribed, err := isSubscribed(ctx, from)
	if err != nil {
		log.Error().Err(err).Msg("Failed to handle message")
		return
	}
	if !subscribed {
		err = text.SendInitialReprompt(ctx, dst, src)
		if err != nil {
			log.Error().Err(err).Msg("Failed to resend initial prompt.")
		}
		return
	}
	previous_messages, err := loadPreviousMessages(ctx, dst, src)
	if err != nil {
		log.Error().Err(err).Str("dst", dst).Str("src", from).Msg("Failed to get previous messages")
		return
	}
	current := llm.Message{
		Content:        body,
		IsFromCustomer: true,
	}
	log.Info().Int("len", len(previous_messages)).Msg("passing")
	next_message, err := llm.GenerateNextMessage(previous_messages, current)
	if err != nil {
		log.Error().Err(err).Str("dst", dst).Str("src", from).Msg("Failed to generate next message")
		return
	}
	text.SendTextFromLLM(next_message.Content)
	log.Info().Str("from", from).Str("from-type", type_).Str("to", to).Str("src", src).Str("dst", dst).Str("body", body).Str("reply", next_message.Content).Msg("Handling text message")
}
