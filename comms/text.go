package comms

import (
	"encoding/json"
	"fmt"

	"github.com/Gleipnir-Technology/nidus-sync/config"
	"github.com/nyaruka/phonenumbers"
	"github.com/rs/zerolog/log"
	"github.com/twilio/twilio-go"
	twilioApi "github.com/twilio/twilio-go/rest/api/v2010"
)

type E164 = phonenumbers.PhoneNumber

func ParsePhoneNumber(input string) (*E164, error) {
	return phonenumbers.Parse(input, "US")
}

func SendText(source E164, destination E164, message string) error {
	log.Info().Msg("Sending text message...")
	// Make sure TWILIO_ACCOUNT_SID and TWILIO_AUTH_TOKEN exists in your environment
	client := twilio.NewRestClient()

	params := &twilioApi.CreateMessageParams{}
	params.SetMessagingServiceSid(config.TwilioMessagingServiceSID)

	params.SetBody(message)
	dest := phonenumbers.Format(&destination, phonenumbers.E164)
	params.SetTo(dest)
	resp, err := client.Api.CreateMessage(params)

	if err != nil {
		return fmt.Errorf("Failed to create message to %s: %w", dest, err)
	} else {
		if resp.Body != nil {
			log.Info().Str("dest", dest).Str("body", *resp.Body).Msg("Text message response")
		} else {
			log.Info().Str("dest", dest).Msg("Text message response is nil")
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
