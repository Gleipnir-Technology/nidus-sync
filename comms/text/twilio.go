package text

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/Gleipnir-Technology/nidus-sync/config"
	"github.com/rs/zerolog/log"
	"github.com/twilio/twilio-go"
	twilioApi "github.com/twilio/twilio-go/rest/api/v2010"
)

func sendTextTwilio(ctx context.Context, source string, destination string, message string) (string, error) {
	client := twilio.NewRestClient()

	params := &twilioApi.CreateMessageParams{}
	params.SetMessagingServiceSid(config.TwilioMessagingServiceSID)

	params.SetBody(message)
	params.SetTo(destination)
	resp, err := client.Api.CreateMessage(params)

	if err != nil {
		return "", fmt.Errorf("Failed to create message to %s: %w", destination, err)
	}
	if resp.Sid == nil {
		log.Warn().Str("src", source).Str("dst", destination).Msg("Text message sid is nil")
		return "", nil
	}
	log.Info().Str("src", source).Str("dst", destination).Str("message", message).Str("sid", *resp.Sid).Msg("Created text message")
	return *resp.Sid, nil
}

func sendSMSTwilio(destination, source, message string) error {
	client := twilio.NewRestClientWithParams(twilio.ClientParams{
		Username: config.TwilioAccountSID,
		Password: config.TwilioAuthToken,
	})
	params := &twilioApi.CreateMessageParams{}
	params.SetTo("+15558675309")
	params.SetFrom("+15017250604")
	params.SetBody("Hello from Go!")

	resp, err := client.Api.CreateMessage(params)
	if err != nil {
		return fmt.Errorf("Error sending SMS message: %w", err)
	}
	response, _ := json.Marshal(*resp)
	log.Debug().Str("response", string(response)).Msg("Send SMS")
	return nil
}
