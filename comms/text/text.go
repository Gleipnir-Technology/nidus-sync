package text

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/Gleipnir-Technology/nidus-sync/config"
	"github.com/Gleipnir-Technology/nidus-sync/db/enums"
	"github.com/nyaruka/phonenumbers"
	"github.com/rs/zerolog/log"
	"github.com/twilio/twilio-go"
	twilioApi "github.com/twilio/twilio-go/rest/api/v2010"
)

type E164 = phonenumbers.PhoneNumber

func ParsePhoneNumber(input string) (*E164, error) {
	return phonenumbers.Parse(input, "US")
}

func sendText(ctx context.Context, source string, destination string, message string, origin enums.CommsTextorigin) error {
	err := ensureInDB(ctx, destination)
	if err != nil {
		return fmt.Errorf("Failed to ensure text message destination is in the DB: %w", err)
	}
	err = insertTextLog(ctx, message, destination, source, origin)
	if err != nil {
		return fmt.Errorf("Failed to insert text message in the DB: %w", err)
	}
	client := twilio.NewRestClient()

	params := &twilioApi.CreateMessageParams{}
	params.SetMessagingServiceSid(config.TwilioMessagingServiceSID)

	params.SetBody(message)
	params.SetTo(destination)
	resp, err := client.Api.CreateMessage(params)

	if err != nil {
		return fmt.Errorf("Failed to create message to %s: %w", destination, err)
	} else {
		if resp.Body != nil {
			log.Info().Str("dest", destination).Str("body", *resp.Body).Msg("Text message response")
		} else {
			log.Info().Str("dest", destination).Msg("Text message response is nil")
		}
	}
	return nil
}

func sendSMS(destination, source, message string) error {
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
