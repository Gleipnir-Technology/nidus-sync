package email

import (
	"context"
	"fmt"

	"github.com/Gleipnir-Technology/nidus-sync/config"
	"github.com/Gleipnir-Technology/nidus-sync/db/enums"
	"github.com/rs/zerolog/log"
)

func NewJobReportNotificationConfirmation(destination, report_id string) Job {
	return jobEmailReportNotificationConfirmation{
		dest:     destination,
		reportID: report_id,
	}
}

type jobEmailReportNotificationConfirmation struct {
	dest     string
	reportID string
}

func (job jobEmailReportNotificationConfirmation) destination() string {
	return job.dest
}
func (job jobEmailReportNotificationConfirmation) messageType() enums.CommsMessagetypeemail {
	return enums.CommsMessagetypeemailReportNotificationConfirmation
}
func (job jobEmailReportNotificationConfirmation) renderHTML() (string, error) {
	_ = newContentEmailNotificationConfirmation(job)
	return "", nil
}
func (job jobEmailReportNotificationConfirmation) renderTXT() (string, error) {
	return "fake txt", nil
}
func (job jobEmailReportNotificationConfirmation) subject() string {
	return ""
}

func sendEmailReportConfirmation(ctx context.Context, job Job) error {
	j, ok := job.(jobEmailReportNotificationConfirmation)
	if !ok {
		return fmt.Errorf("job is not for report subscription confirmation")
	}
	err := maybeSendInitialEmail(ctx, j.destination())
	if err != nil {
		return fmt.Errorf("Failed to handle initial email: %w", err)
	}
	data := make(map[string]string, 0)
	public_id := generatePublicId(enums.CommsMessagetypeemailInitialContact, data)
	data["report_id"] = j.reportID
	report_id_str := publicReportID(j.reportID)
	data["ReportIDStr"] = report_id_str
	data["URLLogo"] = config.MakeURLReport("/static/img/nidus-logo-no-lettering-64.png")
	data["URLReportStatus"] = config.MakeURLReport("/foo")
	data["URLReportUnsubscribe"] = config.MakeURLReport("/email/unsubscribe/report/%s", j.reportID)
	data["URLUnsubscribe"] = urlUnsubscribe(j.destination())
	data["URLViewInBrowser"] = urlEmailInBrowser(public_id)
	text, html, err := renderEmailTemplates(templateReportNotificationConfirmationID, data)
	if err != nil {
		return fmt.Errorf("Failed to render email report notification template: %w", err)
	}
	subject := fmt.Sprintf("Mosquito Report Submission - %s", report_id_str)
	err = insertEmailLog(ctx, data, j.destination(), public_id, config.ForwardEmailReportAddress, subject, templateReportNotificationConfirmationID)
	if err != nil {
		return fmt.Errorf("Failed to store email log: %w", err)
	}
	resp, err := sendEmail(ctx, emailRequest{
		From:    config.ForwardEmailReportAddress,
		HTML:    html,
		Subject: subject,
		Text:    text,
		To:      j.destination(),
	}, enums.CommsMessagetypeemailReportNotificationConfirmation)
	if err != nil {
		return fmt.Errorf("Failed to send email report confirmation to %s for report %s: %w", j.dest, j.reportID, err)
	}
	log.Info().Str("id", resp.ID).Str("dest", j.dest).Str("report_id", j.reportID).Msg("Sent report confirmation email")
	return nil
}

func newContentEmailNotificationConfirmation(job jobEmailReportNotificationConfirmation) (result contentEmailReportConfirmation) {
	result.URLReportStatus = config.MakeURLReport("/status/%s", job.reportID)
	return result
}
