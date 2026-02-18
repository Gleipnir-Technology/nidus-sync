package email

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/Gleipnir-Technology/nidus-sync/config"
	"github.com/rs/zerolog/log"
)

type attachmentRequest struct {
	Filename string `json:"filename"`
	Content  string `json:"content"`
}

type Request struct {
	From        string              `json:"from"`
	To          string              `json:"to"`
	CC          []string            `json:"cc,omitempty"`
	BCC         []string            `json:"bcc,omitempty"`
	Subject     string              `json:"subject"`
	Text        string              `json:"text"`
	HTML        string              `json:"html,omitempty"`
	Attachments []attachmentRequest `json:"attachments,omitempty"`
	Sender      string              `json:"sender"`
	ReplyTo     string              `json:"replyTo,omitempty"`
	InReplyTo   string              `json:"inReplyTo,omitempty"`
	References  []string            `json:"references,omitempty"`
}

type emailEnvelope struct {
	From string   `json:"from"`
	To   []string `json:"to"`
}

type emailResponse struct {
	IsRedacted     bool              `json:"is_redacted"`
	CreatedAt      string            `json:"created_at"`
	HardBounces    []string          `json:"hard_bounces"`
	SoftBounces    []string          `json:"soft_bounces"`
	IsBounce       bool              `json:"is_bounce"`
	Alias          string            `json:"alias"`
	Domain         string            `json:"domain"`
	User           string            `json:"user"`
	Status         string            `json:"status"`
	IsLocked       bool              `json:"is_locked"`
	Envelope       emailEnvelope     `json:"envelope"`
	RequireTLS     bool              `json:"requireTLS"`
	MessageID      string            `json:"messageId"`
	Headers        map[string]string `json:"headers"`
	Date           string            `json:"date"`
	Subject        string            `json:"subject"`
	Accepted       []string          `json:"accepted"`
	Deliveries     []string          `json:"deliveries"`
	RejectedErrors []string          `json:"rejectedErrors"`
	ID             string            `json:"id"`
	Object         string            `json:"object"`
	UpdatedAt      string            `json:"updated_at"`
	Link           string            `json:"link"`
	Message        string            `json:"message"`
}

var FORWARDEMAIL_API = "https://api.forwardemail.net/v1/emails"

func Send(ctx context.Context, email Request) (response emailResponse, err error) {
	payload, err := json.Marshal(email)
	if err != nil {
		return response, fmt.Errorf("Failed to marshal email request: %w", err)
	}

	req, _ := http.NewRequest("POST", FORWARDEMAIL_API, bytes.NewReader(payload))
	req.SetBasicAuth(config.ForwardEmailAPIToken, "")
	req.Header.Add("Content-Type", "application/json")

	res, _ := http.DefaultClient.Do(req)

	defer res.Body.Close()
	body, _ := io.ReadAll(res.Body)

	// Parse the JSON response
	err = json.Unmarshal(body, &response)
	if err != nil {
		log.Warn().Str("status", res.Status).Str("response_body", string(body)).Msg("Attempted to send email but couldn't parse the resulting JSON")
		return response, fmt.Errorf("Failed to unmarshal JSON response: %w", err)
	}
	return response, nil
}
