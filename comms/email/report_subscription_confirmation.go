package email

import (
	"context"
	"fmt"

	"github.com/Gleipnir-Technology/nidus-sync/config"
	"github.com/Gleipnir-Technology/nidus-sync/db/enums"
	//"github.com/rs/zerolog/log"
)

func NewJobReportSubscriptionConfirmation(destination, report_id string) Job {
	return jobEmailReportSubscriptionConfirmation{
		dest:     destination,
		reportID: report_id,
	}
}

type jobEmailReportSubscriptionConfirmation struct {
	dest     string
	reportID string
}

func (job jobEmailReportSubscriptionConfirmation) destination() string {
	return job.dest
}
func (job jobEmailReportSubscriptionConfirmation) messageType() enums.CommsMessagetypeemail {
	return enums.CommsMessagetypeemailReportSubscriptionConfirmation
}
func (job jobEmailReportSubscriptionConfirmation) renderHTML() (string, error) {
	_ = newContentEmailSubscriptionConfirmation(job)
	return "", nil
}
func (job jobEmailReportSubscriptionConfirmation) renderTXT() (string, error) {
	return "fake txt", nil
}
func (job jobEmailReportSubscriptionConfirmation) subject() string {
	return ""
}

func sendEmailReportConfirmation(ctx context.Context, job Job) error {
	j, ok := job.(jobEmailReportSubscriptionConfirmation)
	if !ok {
		return fmt.Errorf("job is not for report subscription confirmation")
	}
	err := maybeSendInitialEmail(ctx, j.destination())
	if err != nil {
		return fmt.Errorf("Failed to handle initial email: %w", err)
	}
	return nil
	/*
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
	*/
}

func newContentEmailSubscriptionConfirmation(job jobEmailReportSubscriptionConfirmation) (result contentEmailReportConfirmation) {
	/*newContentBase(
		&result.Base,
		config.MakeURLReport("/email/report/%s/subscription-confirmation", job.reportID),
	)*/
	result.URLReportStatus = config.MakeURLReport("/status/%s", job.reportID)
	return result
}
