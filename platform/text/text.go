package text

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/Gleipnir-Technology/nidus-sync/config"
	"github.com/Gleipnir-Technology/nidus-sync/db"
	"github.com/Gleipnir-Technology/nidus-sync/db/enums"
	"github.com/Gleipnir-Technology/nidus-sync/db/models"
	"github.com/Gleipnir-Technology/nidus-sync/platform/background"
	"github.com/Gleipnir-Technology/nidus-sync/platform/event"
	"github.com/Gleipnir-Technology/nidus-sync/platform/types"
	"github.com/aarondl/opt/omit"
	"github.com/aarondl/opt/omitnull"
	"github.com/nyaruka/phonenumbers"
	"github.com/rs/zerolog/log"
)

func HandleTextMessage(ctx context.Context, source string, destination string, content string) error {
	src, err := ParsePhoneNumber(source)
	if err != nil {
		return fmt.Errorf("parse source '%s': %w", source, err)
	}
	dst, err := ParsePhoneNumber(destination)
	if err != nil {
		return fmt.Errorf("parse destination '%s': %w", destination, err)
	}
	txn, err := db.PGInstance.BobDB.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("start txn: %w", err)
	}
	defer txn.Rollback(ctx)

	status, err := phoneStatus(ctx, *src)
	if err != nil {
		return fmt.Errorf("Failed to get phone status")
	}
	is_visible_to_llm := status != enums.CommsPhonestatustypeUnconfirmed

	l, err := models.CommsTextLogs.Insert(&models.CommsTextLogSetter{
		//ID:
		Content:        omit.From(content),
		Created:        omit.From(time.Now()),
		Destination:    omit.From(dst.PhoneString()),
		IsVisibleToLLM: omit.From(is_visible_to_llm),
		IsWelcome:      omit.From(false),
		Origin:         omit.From(enums.CommsTextoriginCustomer),
		Source:         omit.From(src.PhoneString()),
		TwilioSid:      omitnull.FromPtr[string](nil),
		TwilioStatus:   omit.From(""),
	}).One(ctx, txn)
	if err != nil {
		return fmt.Errorf("insert text log: %w", err)
	}
	log.Debug().Int32("id", l.ID).Msg("insert comms text log")
	err = background.NewTextRespond(ctx, txn, l.ID)
	if err != nil {
		return fmt.Errorf("text respond: %w", err)
	}
	txn.Commit(ctx)
	return err
}

func respondText(ctx context.Context, log_id int32) error {
	txn, err := db.PGInstance.BobDB.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin tx: %w", err)
	}
	defer txn.Rollback(ctx)
	l, err := models.FindCommsTextLog(ctx, txn, log_id)
	if err != nil {
		return fmt.Errorf("find comms: %w", err)
	}
	src, err := ParsePhoneNumber(l.Source)
	if err != nil {
		return fmt.Errorf("parse source: %w", err)
	}

	status, err := phoneStatus(ctx, *src)
	if err != nil {
		return fmt.Errorf("Failed to get phone status")
	}

	body_l := strings.TrimSpace(strings.ToLower(l.Content))
	// If the user isn't confirmed for sending regular texts ensure they get a reprompt
	if status == enums.CommsPhonestatustypeUnconfirmed {
		switch body_l {
		case "yes":
			err = setPhoneStatus(ctx, txn, *src, enums.CommsPhonestatustypeOkToSend)
			if err != nil {
				return fmt.Errorf("set phone status: %w", err)
			}
			content := "Thanks, we've confirmed your phone number. You can text STOP at any time if you change your mind"
			err = sendTextCommandResponse(ctx, txn, *src, content)
			if err != nil {
				return fmt.Errorf("send response: %w", err)
			}
			handleWaitingTextJobs(ctx, *src)
		// We don't handle 'stop' here because we allow them to say 'stop' at any time, regardless of
		// phone status.
		//case "stop":
		default:
			content := "I have to start with either 'YES' or 'STOP' first, Which do you want?"
			err = sendTextCommandResponse(ctx, txn, *src, content)
			if err != nil {
				log.Error().Err(err).Msg("Failed to resend initial prompt.")
			}
		}
		return nil
	}
	switch body_l {
	case "stop":
		content := "You have successfully been unsubscribed. You will not receive any more messages from this number. Reply START to resubscribe."
		err = sendTextCommandResponse(ctx, txn, *src, content)
		if err != nil {
			log.Error().Err(err).Msg("Failed to send unsubscribe acknowledgement.")
		}
		setPhoneStatus(ctx, txn, *src, enums.CommsPhonestatustypeStopped)
		return nil
	case "reset conversation":
		err = handleResetConversation(ctx, txn, *src)
		if err != nil {
			log.Error().Err(err).Msg("Failed to wipe memory")
			content := "Failed to wipe memory"
			sendTextCommandResponse(ctx, txn, *src, content)
			return fmt.Errorf("reset conversation: %w", err)
		}
		return nil
	}
	// If we've got an open public report from this phone number then we'll let the district respond
	reports, err := reportsForTextRecipient(ctx, txn, *src)
	if err != nil {
		return fmt.Errorf("has open report: %w", err)
	}
	for _, report := range reports {
		_, err = models.PublicreportReportLogs.Insert(&models.PublicreportReportLogSetter{
			Created:    omit.From(time.Now()),
			EmailLogID: omitnull.FromPtr[int32](nil),
			// ID
			ReportID:  omit.From(report.ID),
			TextLogID: omitnull.From(log_id),
			Type:      omit.From(enums.PublicreportReportlogtypeMessageText),
			UserID:    omitnull.FromPtr[int32](nil),
		}).One(ctx, txn)
		if err != nil {
			return fmt.Errorf("insert report log: %w", err)
		}
		event.Updated(event.TypeRMOPublicReport, report.OrganizationID, report.PublicID)
	}
	// If humans are involved, wait for them.
	if len(reports) > 0 {
		return nil
	}
	// Otherwise let the LLM handle the response
	return respondTextLLM(ctx, *src)
}

func respondTextLLM(ctx context.Context, src types.E164) error {
	previous_messages, err := loadPreviousMessagesForLLM(ctx, src)
	if err != nil {
		return fmt.Errorf("Failed to get previous messages: %w", err)
	}
	log.Info().Int("len", len(previous_messages)).Msg("passing")
	next_message, err := generateNextMessage(ctx, previous_messages, src)
	if err != nil {
		return fmt.Errorf("Failed to generate next message: %w", err)
	}
	txn, err := db.PGInstance.BobDB.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("start txn: %w", err)
	}
	defer txn.Rollback(ctx)
	_, err = sendTextDirect(ctx, txn, enums.CommsTextoriginLLM, src.PhoneString(), next_message.Content, true, false)
	if err != nil {
		return fmt.Errorf("Failed to send response text: %w", err)
	}
	txn.Commit(ctx)
	return nil
}

func ParsePhoneNumber(input string) (*types.E164, error) {
	n, err := phonenumbers.Parse(input, "US")
	if err != nil {
		return nil, err
	}
	return types.NewE164(n), nil
}

func StoreSources() error {
	ctx := context.TODO()
	for _, n := range []string{config.PhoneNumberReportStr, config.PhoneNumberSupportStr, config.VoipMSNumber} {
		var err error
		// Deal with Voip.ms not expecting API calls with the prefixed +1
		if !strings.HasPrefix(n, "+1") {
			dest, err := ParsePhoneNumber("+1" + n)
			if err != nil {
				return fmt.Errorf("Failed to parse +1'%s' as phone number: %w", n, err)
			}
			err = EnsureInDB(ctx, db.PGInstance.BobDB, *dest)
		} else {
			dest, err := ParsePhoneNumber(n)
			if err != nil {
				return fmt.Errorf("Failed to parse '%s' as phone number: %w", n, err)
			}
			err = EnsureInDB(ctx, db.PGInstance.BobDB, *dest)
		}
		if err != nil {
			return fmt.Errorf("Failed to add number '%s' to DB: %w", n, err)
		}
	}
	return nil
}

func UpdateMessageStatus(twilio_sid string, status string) {
	ctx := context.TODO()
	l, err := models.CommsTextLogs.Query(
		models.SelectWhere.CommsTextLogs.TwilioSid.EQ(twilio_sid),
	).One(ctx, db.PGInstance.BobDB)
	if err != nil {
		log.Error().Err(err).Str("twilio_sid", twilio_sid).Str("status", status).Msg("Failed to update message status query failed")
		return
	}
	err = l.Update(ctx, db.PGInstance.BobDB, &models.CommsTextLogSetter{
		TwilioStatus: omit.From(status),
	})
	if err != nil {
		log.Error().Err(err).Str("twilio_sid", twilio_sid).Str("status", status).Msg("Failed to update message status update failed")
		return
	}
}
