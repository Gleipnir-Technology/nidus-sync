package comms

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/Gleipnir-Technology/nidus-sync/config"
	"github.com/rs/zerolog/log"
)

var VOIP_MS_API = "https://voip.ms/api/v1/rest.php"

type SendSMSResponse struct {
	Status string `json:"status"`
	SMS int `json:"sms"`
}
func SendSMS(to string, content string) error {
	if len(content) > 160 {
		return errors.New("Message content is more than 160 characters")
	}
	params := url.Values{}
	params.Add("api_password", config.VoipMSPassword)
	params.Add("api_username", config.VoipMSUsername)
	params.Add("method", "sendSMS")
	params.Add("did", config.VoipMSNumber)
	params.Add("dst", to)
	params.Add("message", content)
	// Construct the URL with query parameters
	full_url := VOIP_MS_API + "?" + params.Encode()

	// Make the HTTP request
	resp, err := http.Get(full_url)
	if err != nil {
		log.Warn().Err(err).Str("url", full_url).Msg("Failed to make request to Voip.MS")
		return fmt.Errorf("Error making request: %w", err)
	}
	defer resp.Body.Close()

	// Read the response body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Warn().Err(err).Str("url", full_url).Msg("Failed to read Voip.MS response body")
		return fmt.Errorf("Failed to read response: %w", err)
	}
	log.Info().Str("response", string(body)).Msg("Response from Voip.MS")

	// Parse the JSON response
	var response SendSMSResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		return fmt.Errorf("Failed to unmarshal JSON response: %w", err)
	}
	return nil
}
