package email

import (
	"context"
	"fmt"

	"github.com/Gleipnir-Technology/nidus-sync/config"
	//"github.com/rs/zerolog/log"
)

type contentEmailReportConfirmation struct {
	Base            contentEmailBase
	URLReportStatus string
}

func SendReportConfirmation(ctx context.Context, destination, report_id string) error {
	err := maybeSendInitialEmail(ctx, destination)
	if err != nil {
		return fmt.Errorf("Failed to handle initial email: %w", err)
	}
	data := make(map[string]string, 0)
	data["report_id"] = report_id
	report_id_str := publicReportID(report_id)
	data["ReportIDStr"] = report_id_str
	data["URLLogo"] = config.MakeURLReport("/static/img/nidus-logo-no-lettering-64.png")
	data["URLReportStatus"] = config.MakeURLReport("/status/%s", report_id)
	data["URLReportUnsubscribe"] = config.MakeURLReport("/email/unsubscribe/report/%s", report_id)
	data["URLUnsubscribe"] = urlUnsubscribe(destination)

	subject := fmt.Sprintf("Mosquito Report Submission - %s", report_id_str)
	return sendEmailBegin(ctx, config.ForwardEmailRMOAddress, destination, templateReportNotificationConfirmationID, subject, data)
}

func newContentEmailNotificationConfirmation(report_id string) (result contentEmailReportConfirmation) {
	result.URLReportStatus = config.MakeURLReport("/status/%s", report_id)
	return result
}
