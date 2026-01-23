package comms

import (
	"bytes"
	"context"
	"embed"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/Gleipnir-Technology/nidus-sync/config"
	"github.com/Gleipnir-Technology/nidus-sync/db"
	"github.com/Gleipnir-Technology/nidus-sync/db/enums"
	"github.com/Gleipnir-Technology/nidus-sync/db/models"
	"github.com/Gleipnir-Technology/nidus-sync/htmlpage"
	"github.com/aarondl/opt/omit"
	"github.com/rs/zerolog/log"
)

func RenderEmailInitial(w http.ResponseWriter, destination string) {
	content := newContentEmailInitial(destination)
	renderOrError(w, initialT, content)
}

func RenderEmailReportConfirmation(w http.ResponseWriter, report_id string) {
	content := newContentEmailSubscriptionConfirmation(report_id)
	renderOrError(w, reportConfirmationT, content)
}

func SendEmailInitialContact(ctx context.Context, destination string) error {
	content := newContentEmailInitial(destination)
	text, html, err := renderEmailTemplates(reportConfirmationT, content)
	if err != nil {
		return fmt.Errorf("Failed to render email temlate: %w", err)
	}
	resp, err := sendEmail(ctx, emailRequest{
		From:    config.ForwardEmailReportAddress,
		HTML:    html,
		Subject: "Welcome",
		Text:    text,
		To:      destination,
	}, enums.CommsMessagetypeemailInitialContact)
	if err != nil {
		return fmt.Errorf("Failed to send email to %s: %w", err)
	}
	log.Info().Str("id", resp.ID).Str("to", destination).Msg("Sent initial contact email")
	return nil
}

func SendEmailReportConfirmation(ctx context.Context, to string, report_id string) error {
	report_id_str := publicReportID(report_id)
	content := newContentEmailSubscriptionConfirmation(report_id)
	text, html, err := renderEmailTemplates(reportConfirmationT, content)
	if err != nil {
		return fmt.Errorf("Failed to render template %s: %w", reportConfirmationT.name, err)
	}
	resp, err := sendEmail(ctx, emailRequest{
		From:    config.ForwardEmailReportAddress,
		HTML:    html,
		Subject: fmt.Sprintf("Mosquito Report Submission - %s", report_id_str),
		Text:    text,
		To:      to,
	}, enums.CommsMessagetypeemailReportSubscriptionConfirmation)
	if err != nil {
		return fmt.Errorf("Failed to send email report confirmation to %s for report %s: %w", to, report_id, err)
	}
	log.Info().Str("id", resp.ID).Str("to", to).Str("report_id", report_id).Msg("Sent report confirmation email")
	return nil
}

var (
	initialT            = buildTemplate("initial")
	reportConfirmationT = buildTemplate("report-subscription-confirmation")
)

//go:embed template/*
var embeddedFiles embed.FS

type attachmentRequest struct {
	Filename string `json:"filename"`
	Content  string `json:"content"`
}

type contentEmailBase struct {
	URLLogo          string
	URLUnsubscribe   string
	URLViewInBrowser string
}

type contentEmailReportConfirmation struct {
	Base            contentEmailBase
	URLReportStatus string
}
type contentEmailInitial struct {
	Base         contentEmailBase
	Destination  string
	URLSubscribe string
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

func newContentBase(b *contentEmailBase, url string) {
	b.URLLogo = config.MakeURLReport("/static/img/nidus-logo-no-lettering-64.png")
	b.URLUnsubscribe = config.MakeURLReport("/email/unsubscribe")
	b.URLViewInBrowser = url
}

func newContentEmailInitial(destination string) (result contentEmailInitial) {
	newContentBase(
		&result.Base,
		config.MakeURLReport("/email/initial"),
	)
	result.Destination = destination
	result.URLSubscribe = config.MakeURLReport("/email/subscribe?email=%s", destination)
	return result
}
func newContentEmailSubscriptionConfirmation(report_id string) (result contentEmailReportConfirmation) {
	newContentBase(
		&result.Base,
		config.MakeURLReport("/email/report/%s/subscription-confirmation", report_id),
	)
	result.URLReportStatus = config.MakeURLReport("/status/%s", report_id)
	return result
}

func publicReportID(s string) string {
	if len(s) != 12 {
		return s
	}
	return s[0:4] + "-" + s[4:8] + "-" + s[8:12]
}

func renderOrError(w http.ResponseWriter, template *builtTemplate, context interface{}) {
	buf := &bytes.Buffer{}
	err := template.executeTemplateHTML(buf, context)
	if err != nil {
		log.Error().Err(err).Str("name", template.name).Msg("Failed to render template")
		htmlpage.RespondError(w, "Failed to render template", err, http.StatusInternalServerError)
		return
	}
	buf.WriteTo(w)
}

var FORWARDEMAIL_API = "https://api.forwardemail.net/v1/emails"

func ensureInDB(ctx context.Context, destination string) (err error) {
	_, err = models.FindCommsEmail(ctx, db.PGInstance.BobDB, destination)
	if err != nil {
		// assume it exists
		log.Warn().Err(err).Msg("ElI, check what this error should look like")
		return nil
	}
	_, err = models.CommsEmails.Insert(&models.CommsEmailSetter{
		Address:      omit.From(destination),
		Confirmed:    omit.From(false),
		IsSubscribed: omit.From(false),
	}).One(ctx, db.PGInstance.BobDB)
	if err != nil {
		return fmt.Errorf("Failed to insert new email: %w", err)
	}
	log.Info().Str("email", destination).Msg("Added email to the comms database")
	return nil
}

func insertEmailLog(ctx context.Context, email emailRequest, t enums.CommsMessagetypeemail) (err error) {
	_, err = models.CommsEmailLogs.Insert(&models.CommsEmailLogSetter{
		Created:     omit.From(time.Now()),
		Destination: omit.From(email.To),
		Source:      omit.From(email.From),
		Type:        omit.From(t),
	}).One(ctx, db.PGInstance.BobDB)
	return err
}
func sendEmail(ctx context.Context, email emailRequest, t enums.CommsMessagetypeemail) (response emailResponse, err error) {
	ensureInDB(ctx, email.To)
	payload, err := json.Marshal(email)
	if err != nil {
		return response, fmt.Errorf("Failed to marshal email request: %w", err)
	}

	insertEmailLog(ctx, email, t)
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
