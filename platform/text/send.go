package text

import (
	"context"
	"fmt"
	"time"

	"github.com/Gleipnir-Technology/bob"
	"github.com/Gleipnir-Technology/nidus-sync/comms/text"
	"github.com/Gleipnir-Technology/nidus-sync/config"
	"github.com/aarondl/opt/omit"
	"github.com/aarondl/opt/omitnull"
	//"github.com/Gleipnir-Technology/nidus-sync/db"
	"github.com/Gleipnir-Technology/nidus-sync/db/enums"
	"github.com/Gleipnir-Technology/nidus-sync/db/models"
	"github.com/Gleipnir-Technology/nidus-sync/platform/background"
	"github.com/Gleipnir-Technology/nidus-sync/platform/types"
	"github.com/rs/zerolog/log"
)

func ensureInitialText(ctx context.Context, txn bob.Executor, dst types.E164) error {
	rows, err := models.CommsTextLogs.Query(
		models.SelectWhere.CommsTextLogs.Destination.EQ(dst.PhoneString()),
		models.SelectWhere.CommsTextLogs.IsWelcome.EQ(true),
	).All(ctx, txn)
	if err != nil {
		return fmt.Errorf("Failed to query text logs: %w", err)
	}
	if len(rows) > 0 {
		return nil
	}
	return sendInitialText(ctx, txn, dst)
}
func insertTextLog(ctx context.Context, txn bob.Executor, destination types.E164, source types.E164, origin enums.CommsTextorigin, content string, is_welcome bool, is_visible_to_llm bool) (l *models.CommsTextLog, err error) {
	l, err = models.CommsTextLogs.Insert(&models.CommsTextLogSetter{
		//ID:
		Content:        omit.From(content),
		Created:        omit.From(time.Now()),
		Destination:    omit.From(destination.PhoneString()),
		IsVisibleToLLM: omit.From(is_visible_to_llm),
		IsWelcome:      omit.From(is_welcome),
		Origin:         omit.From(origin),
		Source:         omit.From(source.PhoneString()),
		TwilioSid:      omitnull.FromPtr[string](nil),
		TwilioStatus:   omit.From(""),
	}).One(ctx, txn)

	return l, err
}
func resendInitialText(ctx context.Context, txn bob.Executor, dst types.E164) error {
	phone, err := models.FindCommsPhone(ctx, txn, dst.PhoneString())
	if err != nil {
		return fmt.Errorf("Failed to find phone %s: %w", dst, err)
	}
	err = phone.Update(ctx, txn, &models.CommsPhoneSetter{
		Status: omit.From(enums.CommsPhonestatustypeUnconfirmed),
	})
	if err != nil {
		return fmt.Errorf("Failed to clear subscription on phone %s: %w", dst, err)
	}
	return nil
}
func sendInitialText(ctx context.Context, txn bob.Executor, dst types.E164) error {
	content := "Welcome to Report Mosquitoes Online. We received your request and want to confirm text updates. Reply YES to continue. Reply STOP at any time to unsubscribe"
	_, err := sendTextDirect(ctx, txn, enums.CommsTextoriginWebsiteAction, dst.PhoneString(), content, false, true)
	if err != nil {
		return fmt.Errorf("send text: %w", err)
	}
	return nil
}

// Begin the process of sending the text message, but only get as far as adding it to
// the database, then let the backend finish sending.
func sendTextBegin(ctx context.Context, txn bob.Executor, user_id *int32, report_id *int32, destination types.E164, content string, type_ enums.CommsTextjobtype) (*int32, error) {
	err := EnsureInDB(ctx, txn, destination)
	if err != nil {
		return nil, fmt.Errorf("Failed to ensure text message destination is in the DB: %w", err)
	}
	job, err := models.CommsTextJobs.Insert(&models.CommsTextJobSetter{
		Content:     omit.From(content),
		CreatorID:   omitnull.FromPtr(user_id),
		Created:     omit.From(time.Now()),
		Destination: omit.From(destination.PhoneString()),
		//ID:
		ReportID: omitnull.FromPtr(report_id),
		Source:   omit.From(enums.CommsTextjobsourceRmo),
		Type:     omit.From(type_),
	}).One(ctx, txn)
	if err != nil {
		return nil, fmt.Errorf("Failed to add delayed text job: %w", err)
	}
	err = background.NewTextSend(ctx, txn, job.ID)
	if err != nil {
		return nil, fmt.Errorf("new background job: %w", err)
	}
	return &job.ID, nil
}
func sendTextCommandResponse(ctx context.Context, txn bob.Executor, dst types.E164, content string) error {
	_, err := sendTextDirect(ctx, txn, enums.CommsTextoriginCommandResponse, dst.PhoneString(), content, false, false)
	return err
}
func sendTextComplete(ctx context.Context, txn bob.Executor, job *models.CommsTextJob) error {
	dst, err := ParsePhoneNumber(job.Destination)
	if err != nil {
		return fmt.Errorf("parse phone: %w", err)
	}
	var origin enums.CommsTextorigin
	switch job.Type {
	case enums.CommsTextjobtypeReportConfirmation:
		origin = enums.CommsTextoriginWebsiteAction
	case enums.CommsTextjobtypeReportMessage:
		origin = enums.CommsTextoriginDistrict
	default:
		return fmt.Errorf("incomplete switch: %s", string(job.Type))
	}
	status, err := phoneStatus(ctx, *dst)
	if err != nil {
		return fmt.Errorf("Failed to check if subscribed: %w", err)
	}
	log.Debug().Str("phone status", string(status)).Str("destination", job.Destination).Send()
	switch status {
	case enums.CommsPhonestatustypeUnconfirmed:
		err := ensureInitialText(ctx, txn, *dst)
		if err != nil {
			return fmt.Errorf("Failed to ensure initial text has been sent: %w", err)
		}
		return nil
	case enums.CommsPhonestatustypeOkToSend:
		_, err = sendTextDirect(ctx, txn, origin, dst.PhoneString(), job.Content, true, false)
		if err != nil {
			return fmt.Errorf("Failed to send report subscription confirmation: %w", err)
		}
		return nil
	case enums.CommsPhonestatustypeStopped:
		resendInitialText(ctx, txn, *dst)
	}
	text_log, err := sendTextDirect(ctx, txn, origin, job.Destination, job.Content, true, false)
	if err != nil {
		return fmt.Errorf("send text direct: %w", err)
	}
	err = job.Update(ctx, txn, &models.CommsTextJobSetter{
		Completed: omitnull.From(time.Now()),
	})
	if err != nil {
		return fmt.Errorf("update job: %w", err)
	}
	if job.ReportID.IsValue() {
		_, err := models.ReportTexts.Insert(&models.ReportTextSetter{
			CreatorID: omit.From(job.CreatorID.MustGet()),
			ReportID:  omit.From(job.ReportID.MustGet()),
			TextLogID: omit.From(text_log.ID),
		}).One(ctx, txn)
		if err != nil {
			return fmt.Errorf("insert report_text: %w", err)
		}
	}
	return nil
}

// Send a text message and save the appropriate database records.
// Send immediately using the current goroutine
func sendTextDirect(ctx context.Context, txn bob.Executor, origin enums.CommsTextorigin, destination, content string, is_visible_to_llm, is_welcome bool) (*models.CommsTextLog, error) {
	text_log, err := models.CommsTextLogs.Insert(&models.CommsTextLogSetter{
		//ID:
		Content:        omit.From(content),
		Created:        omit.From(time.Now()),
		Destination:    omit.From(destination),
		IsVisibleToLLM: omit.From(is_visible_to_llm),
		IsWelcome:      omit.From(is_welcome),
		Origin:         omit.From(origin),
		Source:         omit.From(config.PhoneNumberReportStr),
		TwilioSid:      omitnull.FromPtr[string](nil),
		TwilioStatus:   omit.From(""),
	}).One(ctx, txn)
	if err != nil {
		return nil, fmt.Errorf("insert text log: %w", err)
	}
	pid, err := text.SendText(ctx, config.VoipMSNumber, destination, content)
	if err != nil {
		return nil, fmt.Errorf("send text: %w", err)
	}
	err = text_log.Update(ctx, txn, &models.CommsTextLogSetter{
		TwilioSid:    omitnull.From(pid),
		TwilioStatus: omit.From("created"),
	})
	if err != nil {
		return nil, fmt.Errorf("update  %w", err)
	}

	return text_log, nil
}
