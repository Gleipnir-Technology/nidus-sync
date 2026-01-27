package platform

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/Gleipnir-Technology/nidus-sync/comms/text"
	"github.com/Gleipnir-Technology/nidus-sync/config"
	"github.com/Gleipnir-Technology/nidus-sync/db"
	"github.com/Gleipnir-Technology/nidus-sync/db/enums"
	"github.com/Gleipnir-Technology/nidus-sync/db/models"
	"github.com/Gleipnir-Technology/nidus-sync/db/sql"
	"github.com/Gleipnir-Technology/nidus-sync/llm"
	"github.com/aarondl/opt/omit"
	"github.com/aarondl/opt/omitnull"
	"github.com/nyaruka/phonenumbers"
	"github.com/rs/zerolog/log"
)

func HandleTextMessage(from string, to string, body string) {
	ctx := context.Background()
	type_, src := splitPhoneSource(from)
	dst, err := getDst(ctx, to)
	if err != nil {
		log.Error().Err(err).Str("to", to).Msg("Failed to get dst")
		return
	}

	_, err = insertTextLog(ctx, body, dst, src, enums.CommsTextoriginCustomer, false)
	if err != nil {
		log.Error().Err(err).Str("dst", dst).Msg("Failed to add text message log")
		return
	}
	subscribed, err := isSubscribed(ctx, src)
	if err != nil {
		log.Error().Err(err).Msg("Failed to handle message")
		return
	}
	// We don't know if they're subscribed or not.
	if subscribed == nil {
		body_l := strings.TrimSpace(strings.ToLower(body))
		switch body_l {
		case "stop":
			setSubscribed(ctx, src, false)
		case "yes":
			setSubscribed(ctx, src, true)
			handleWaitingTextJobs(ctx, src)
		default:
			content := "I have to start with either 'YES' or 'STOP' first, Which do you want?"
			/*err := insertTextLog(ctx, body, src, dst, enums.CommsTextoriginReiteration, false)
			if err != nil {
				log.Error().Err(err).Msg("Failed to add reiteration to the text log")
				return
			}*/
			err = sendText(ctx, src, dst, content, enums.CommsTextoriginReiteration, false)
			if err != nil {
				log.Error().Err(err).Msg("Failed to resend initial prompt.")
			}
		}
		return
	}
	previous_messages, err := loadPreviousMessages(ctx, dst, src)
	if err != nil {
		log.Error().Err(err).Str("dst", dst).Str("src", from).Msg("Failed to get previous messages")
		return
	}
	log.Info().Int("len", len(previous_messages)).Msg("passing")
	next_message, err := llm.GenerateNextMessage(ctx, previous_messages, src)
	if err != nil {
		log.Error().Err(err).Str("dst", dst).Str("src", from).Msg("Failed to generate next message")
		return
	}
	/*
		err = insertTextLog(ctx, next_message.Content, src, dst, enums.CommsTextoriginLLM, false)
		if err != nil {
			log.Error().Err(err).Str("dst", dst).Msg("Failed to insert new text message to the text log")
			return
		}
	*/
	err = sendText(ctx, dst, src, next_message.Content, enums.CommsTextoriginLLM, false)
	if err != nil {
		log.Error().Err(err).Str("src", src).Str("dst", dst).Str("content", next_message.Content).Msg("Failed to send response text")
		return
	}
	log.Info().Str("from", from).Str("from-type", type_).Str("to", to).Str("src", src).Str("dst", dst).Str("body", body).Str("reply", next_message.Content).Msg("Handled text message")
}

func TextStoreSources() error {
	ctx := context.TODO()
	src := phonenumbers.Format(&config.PhoneNumberReport, phonenumbers.E164)
	return ensureInDB(ctx, src)
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
func delayMessage(ctx context.Context, source string, destination string, content string, type_ enums.CommsTextjobtype) error {
	job, err := models.CommsTextJobs.Insert(&models.CommsTextJobSetter{
		Content:     omit.From(content),
		Created:     omit.From(time.Now()),
		Destination: omit.From(destination),
		//ID:
		Type: omit.From(type_),
	}).One(ctx, db.PGInstance.BobDB)
	if err != nil {
		return fmt.Errorf("Failed to add delayed text job: %w", err)
	}
	log.Info().Int32("id", job.ID).Msg("Created delayed text job")
	return nil
}

func ensureInitialText(ctx context.Context, src string, dst string) error {
	//
	origin := enums.CommsTextoriginWebsiteAction
	rows, err := models.CommsTextLogs.Query(
		models.SelectWhere.CommsTextLogs.Destination.EQ(dst),
		models.SelectWhere.CommsTextLogs.IsWelcome.EQ(true),
	).All(ctx, db.PGInstance.BobDB)
	if err != nil {
		return fmt.Errorf("Failed to query text logs: %w", err)
	}
	if len(rows) > 0 {
		return nil
	}
	content := "Welcome to Report Mosquitoes Online. We received your request and want to confirm text updates. Reply YES to continue. Reply STOP at any time to unsubscribe"
	err = sendText(ctx, src, dst, content, origin, true)
	if err != nil {
		return fmt.Errorf("Failed to send initial confirmation: %w", err)
	}
	return nil
}

func ensureInDB(ctx context.Context, destination string) (err error) {
	_, err = models.FindCommsPhone(ctx, db.PGInstance.BobDB, destination)
	if err != nil {
		// doesn't exist
		if err.Error() == "sql: no rows in result set" {
			_, err = models.CommsPhones.Insert(&models.CommsPhoneSetter{
				E164:         omit.From(destination),
				IsSubscribed: omitnull.FromPtr[bool](nil),
			}).One(ctx, db.PGInstance.BobDB)
			if err != nil {
				return fmt.Errorf("Failed to insert new phone contact: %w", err)
			}
			log.Info().Str("phone", destination).Msg("Added text to the comms database")
			return nil
		}
		return fmt.Errorf("Unexpected error searching for phone contact: %w", err)
	}
	return nil
}

// Translate from Twilio's representation of a RCS message sender to our concept of a phone number
// From: rcs:dev_report_mosquitoes_online_dosrvwxm_agent
// To: +16235525879
func getDst(ctx context.Context, to string) (string, error) {

	if to == config.TwilioRCSSenderRMO {
		return config.PhoneNumberReportStr, nil
	}
	/*
		phone, err := models.FindCommsPhone(ctx, db.PGInstance.BobDB, to)
		if err != nil {
			return "", fmt.Errorf("Failed to search for dest phone %s: %w", to, err)
		}
		return phone.E164, nil
	*/
	return "", fmt.Errorf("Cannot match phone number to '%s'", to)
}

func handleWaitingTextJobs(ctx context.Context, src string) {
	log.Info().Str("src", src).Msg("Pretend handle waiting jobs")

}

func insertTextLog(ctx context.Context, content string, destination string, source string, origin enums.CommsTextorigin, is_welcome bool) (log *models.CommsTextLog, err error) {
	log, err = models.CommsTextLogs.Insert(&models.CommsTextLogSetter{
		//ID:
		Content:      omit.From(content),
		Created:      omit.From(time.Now()),
		Destination:  omit.From(destination),
		IsWelcome:    omit.From(is_welcome),
		Origin:       omit.From(origin),
		Source:       omit.From(source),
		TwilioSid:    omitnull.FromPtr[string](nil),
		TwilioStatus: omit.From(""),
	}).One(ctx, db.PGInstance.BobDB)

	return log, err
}

func isSubscribed(ctx context.Context, src string) (*bool, error) {
	phone, err := models.FindCommsPhone(ctx, db.PGInstance.BobDB, src)
	if err != nil {
		return nil, fmt.Errorf("Failed to determine if '%s' is subscribed: %w", src, err)
	}
	if phone.IsSubscribed.IsNull() {
		return nil, nil
	}
	result := phone.IsSubscribed.MustGet()
	return &result, nil
}

func loadPreviousMessages(ctx context.Context, dst, src string) ([]llm.Message, error) {
	messages, err := sql.TextsBySenders(dst, src).All(ctx, db.PGInstance.BobDB)
	results := make([]llm.Message, 0)
	if err != nil {
		return results, fmt.Errorf("Failed to get message history for %s and %s: %w", dst, src, err)
	}
	log.Info().Int("count", len(messages)).Str("src", src).Str("dst", dst).Msg("Found previous messages")
	for _, m := range messages {
		is_from_customer := (m.Source == src)
		results = append(results, llm.Message{
			IsFromCustomer: is_from_customer,
			Content:        m.Content,
		})
	}
	return results, nil
}

func sendText(ctx context.Context, source string, destination string, message string, origin enums.CommsTextorigin, is_welcome bool) error {
	err := ensureInDB(ctx, destination)
	if err != nil {
		return fmt.Errorf("Failed to ensure text message destination is in the DB: %w", err)
	}
	log, err := insertTextLog(ctx, message, destination, source, origin, is_welcome)
	if err != nil {
		return fmt.Errorf("Failed to insert text message in the DB: %w", err)
	}
	sid, err := text.SendText(ctx, source, destination, message)
	if err != nil {
		return fmt.Errorf("Failed to send text message: %w", err)
	}
	err = log.Update(ctx, db.PGInstance.BobDB, &models.CommsTextLogSetter{
		TwilioSid:    omitnull.From(sid),
		TwilioStatus: omit.From("created"),
	})

	return nil
}

func splitPhoneSource(s string) (string, string) {
	parts := strings.Split(s, ":")
	switch len(parts) {
	case 0:
		return "this isn't", "possible"
	case 1:
		return "", s
	case 2:
		return parts[0], parts[1]
	default:
		log.Warn().Str("s", s).Msg("Got an incomprehensible number of parts of a phone number")
		return parts[0], parts[1]
	}

}

func setSubscribed(ctx context.Context, src string, is_subscribed bool) error {
	phone, err := models.FindCommsPhone(ctx, db.PGInstance.BobDB, src)
	if err != nil {
		return fmt.Errorf("Failed to determine if '%s' is subscribed: %w", src, err)
	}
	phone.Update(ctx, db.PGInstance.BobDB, &models.CommsPhoneSetter{
		IsSubscribed: omitnull.From(is_subscribed),
	})
	log.Info().Str("src", src).Bool("is_subscribed", is_subscribed).Msg("Set number subscribed")
	return nil
}
