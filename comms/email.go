package comms

import (
	"bytes"
	"embed"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/Gleipnir-Technology/nidus-sync/config"
	"github.com/rs/zerolog/log"
)

func SendEmailReportConfirmation(to string, report_id string) error {
	content := contentEmailReportConfirmation{
		URLLogo:              "https://dev-sync.nidus.cloud/static/img/nidus-logo-no-lettering-64.png",
		URLReportStatus:      fmt.Sprintf("https://dev-sync.nidus.cloud/report/%s", report_id),
		URLReportUnsubscribe: fmt.Sprintf("https://dev-sync.nidus.cloud/report/%s/unsubscribe", report_id),
		URLViewInBrowser:     fmt.Sprintf("https://dev-sync.nidus.cloud/report/%s/subscribe-confirmation", report_id),
	}
	text, html, err := renderEmailTemplates(reportConfirmationT, content)
	if err != nil {
		return fmt.Errorf("Failed to render template %s: %w", reportConfirmationT.name, err)
	}
	resp, err := sendEmail(emailRequest{
		From:    config.ForwardEmailReportAddress,
		HTML:    html,
		Subject: fmt.Sprintf("Mosquito Report Submission - %s", report_id),
		Text:    text,
		To:      to,
	})
	if err != nil {
		return fmt.Errorf("Failed to send email report confirmation to %s for report %s: %w", to, report_id, err)
	}
	log.Info().Str("id", resp.ID).Str("to", to).Str("report_id", report_id).Msg("Sent report confirmation email")
	return nil
}

var (
	reportConfirmationT = buildTemplate("report-subscription-confirmation")
)

//go:embed template/*
var embeddedFiles embed.FS

type attachmentRequest struct {
	Filename string `json:"filename"`
	Content  string `json:"content"`
}

type contentEmailReportConfirmation struct {
	URLLogo              string
	URLReportStatus      string
	URLReportUnsubscribe string
	URLViewInBrowser     string
}

type emailRequest struct {
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

func sendEmail(email emailRequest) (response emailResponse, err error) {
	url := "https://api.forwardemail.net/v1/emails"

	payload, err := json.Marshal(email)
	if err != nil {
		return response, fmt.Errorf("Failed to marshal email request: %w", err)
	}

	req, _ := http.NewRequest("POST", url, bytes.NewReader(payload))
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
