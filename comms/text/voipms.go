package text

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"

	"github.com/Gleipnir-Technology/nidus-sync/config"
	"github.com/rs/zerolog/log"
)

var VOIP_MS_API = "https://voip.ms/api/v1/rest.php"

type VoipMSResponse struct {
	MMS     int    `json:"mms"`
	Message string `json:"message"`
	Status  string `json:"status"`
	SMS     int    `json:"sms"`
}

func sendTextVoipms(ctx context.Context, to string, content string, media ...string) (string, error) {
	if len(content) > 2048 {
		return "", errors.New("Message content is more than 160 characters")
	}
	params := url.Values{}
	params.Add("api_password", config.VoipMSPassword)
	params.Add("api_username", config.VoipMSUsername)
	params.Add("method", "sendMMS")
	params.Add("did", config.VoipMSNumber)
	params.Add("dst", to)
	params.Add("message", content)
	/*
	   for i, med := range media {
	           // These should be one of:
	           //  1. A full URL that the service cat GET
	           //  2. A base64-encoded image starting with "data:image/png;base64,iVBORw0KGgoAAAANSUh..."
	           params.Add(fmt.Sprintf("media%d", i+1), med)
	   }
	   params.Add(fmt.Sprintf("media%d", len(media)+1), "")
	*/

	response, err := makeVoipMSRequest(params)
	if err != nil {
		return "", fmt.Errorf("Failed to send MMS: %w", err)
	}
	log.Info().Str("status", response.Status).Int("mms", response.MMS).Msg("Sent MMS message")
	return strconv.Itoa(response.MMS), nil
}

func sendSMSVoipms(to string, content string) (string, error) {
	if len(content) > 160 {
		return "", errors.New("Message content is more than 160 characters")
	}
	params := url.Values{}
	params.Add("api_password", config.VoipMSPassword)
	params.Add("api_username", config.VoipMSUsername)
	params.Add("method", "sendSMS")
	params.Add("did", config.VoipMSNumber)
	params.Add("dst", to)
	params.Add("message", content)

	response, err := makeVoipMSRequest(params)
	if err != nil {
		return "", fmt.Errorf("Failed to send SMS: %w", err)
	}
	log.Info().Str("status", response.Status).Int("sms", response.SMS).Msg("Sent MMS message")
	return strconv.Itoa(response.SMS), nil
}

func makeVoipMSRequest(params url.Values) (VoipMSResponse, error) {
	result := VoipMSResponse{}
	// Construct the URL with query parameters
	full_url := VOIP_MS_API + "?" + params.Encode()

	// Make the HTTP request
	log.Debug().Str("full_url", full_url).Msg("Sending command to VoIP.ms")
	resp, err := http.Get(full_url)
	if err != nil {
		log.Warn().Err(err).Str("url", full_url).Msg("Failed to make request to Voip.MS")
		return result, fmt.Errorf("Error making request: %w", err)
	}
	defer resp.Body.Close()

	// Read the response body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Warn().Err(err).Str("url", full_url).Msg("Failed to read Voip.MS response body")
		return result, fmt.Errorf("Failed to read response: %w", err)
	}
	log.Info().Str("response", string(body)).Msg("Response from Voip.MS")

	// Parse the JSON response
	var response VoipMSResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		return result, fmt.Errorf("Failed to unmarshal JSON response: %w", err)
	}
	return response, nil
}
