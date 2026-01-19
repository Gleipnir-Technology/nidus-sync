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
	return sendEmail(emailRequest{
		From:    config.ForwardEmailReportAddress,
		HTML:    html,
		Subject: fmt.Sprintf("Mosquito Report %s Submission", report_id),
		Text:    text,
		To:      to,
	})
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

type emailResponse struct {
	Message string `json:"message"`
}

func sendEmail(email emailRequest) error {
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
