package email

import (
	"context"
	"fmt"

	"github.com/Gleipnir-Technology/nidus-sync/config"
	"github.com/rs/zerolog/log"
	"resty.dev/v3"
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

type emailResponseError struct {
	StatusCode int    `json:"statusCode"`
	Error      string `json:"error"`
	Message    string `json:"message"`
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

var FORWARDEMAIL_EMAIL_POST_API = "https://api.forwardemail.net/v1/emails"

func Send(ctx context.Context, email Request) (result emailResponse, err error) {
	client := resty.New()

	var err_resp emailResponseError
	r, err := client.R().
		SetBasicAuth(config.ForwardEmailAPIToken, "").
		SetBody(email).
		SetContext(ctx).
		SetError(&err_resp).
		SetHeader("Content-Type", "application/json").
		SetResult(&result).
		Post(FORWARDEMAIL_EMAIL_POST_API)

	if err != nil {
		return result, fmt.Errorf("Failed to marshal email request: %w", err)
	}

	if r.IsError() {
		log.Error().Int("status", err_resp.StatusCode).Str("error", err_resp.Error).Str("msg", err_resp.Message).Msg("Email send error")
		return result, fmt.Errorf("Error response %d from email service: %s (%s)", err_resp.StatusCode, err_resp.Message, err_resp.Error)
	}
	return result, nil
}
