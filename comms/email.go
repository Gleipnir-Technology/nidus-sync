package comms

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/Gleipnir-Technology/nidus-sync/config"
	"github.com/rs/zerolog/log"
)

type AttachmentRequest struct {
	Filename string `json:"filename"`
	Content  string `json:"content"`
}

type EmailRequest struct {
	From        string              `json:"from"`
	To          string              `json:"to"`
	CC          []string            `json:"cc,omitempty"`
	BCC         []string            `json:"bcc,omitempty"`
	Subject     string              `json:"subject"`
	Text        string              `json:"text"`
	HTML        string              `json:"html,omitempty"`
	Attachments []AttachmentRequest `json:"attachments,omitempty"`
	Sender      string              `json:"sender"`
	ReplyTo     string              `json:"replyTo,omitempty"`
	InReplyTo   string              `json:"inReplyTo,omitempty"`
	References  []string            `json:"references,omitempty"`
}

type EmailResponse struct {
	Message string `json:"message"`
}

func SendEmail(email EmailRequest) error {
	url := "https://api.forwardemail.net/v1/emails"

	payload, err := json.Marshal(email)
	if err != nil {
		return fmt.Errorf("Failed to marshal email request: %w", err)
	}

	req, _ := http.NewRequest("POST", url, bytes.NewReader(payload))
	req.SetBasicAuth(config.ForwardEmailAPIToken, "")
	req.Header.Add("Content-Type", "application/json")

	res, _ := http.DefaultClient.Do(req)

	defer res.Body.Close()
	body, _ := io.ReadAll(res.Body)

	log.Info().Str("status", res.Status).Str("response_body", string(body)).Msg("Attempted to send email")
	return nil
}
