package text

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/Gleipnir-Technology/bob/dialect/psql"
	"github.com/Gleipnir-Technology/bob/dialect/psql/um"
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

type E164 = phonenumbers.PhoneNumber

func HandleTextMessage(src string, dst string, body string) {
	ctx := context.Background()

	_, err := insertTextLog(ctx, body, dst, src, enums.CommsTextoriginCustomer, false, true)
	if err != nil {
		log.Error().Err(err).Str("dst", dst).Msg("Failed to add text message log")
		return
	}
	subscribed, err := isSubscribed(ctx, src)
	if err != nil {
		log.Error().Err(err).Msg("Failed to handle message")
		return
	}
	body_l := strings.TrimSpace(strings.ToLower(body))
	// We don't know if they're subscribed or not.
	if subscribed == nil {
		switch body_l {
		case "yes":
			setSubscribed(ctx, src, true)
			handleWaitingTextJobs(ctx, src)
		default:
			content := "I have to start with either 'YES' or 'STOP' first, Which do you want?"
			err = sendText(ctx, dst, src, content, enums.CommsTextoriginReiteration, false, false)
			if err != nil {
				log.Error().Err(err).Msg("Failed to resend initial prompt.")
			}
		}
		return
	}
	switch body_l {
	case "stop":
		content := "You have successfully been unsubscribed. You will not receive any more messages from this number. Reply START to resubscribe."
		err = sendText(ctx, dst, src, content, enums.CommsTextoriginCommandResponse, false, false)
		if err != nil {
			log.Error().Err(err).Msg("Failed to send unsubscribe acknowledgement.")
		}
		setSubscribed(ctx, src, false)
		return
	case "reset conversation":
		handleResetConversation(ctx, src, dst)
		return
	default:
	}
	previous_messages, err := loadPreviousMessagesForLLM(ctx, dst, src)
	if err != nil {
		log.Error().Err(err).Str("dst", dst).Str("src", src).Msg("Failed to get previous messages")
		return
	}
	log.Info().Int("len", len(previous_messages)).Msg("passing")
	next_message, err := generateNextMessage(ctx, previous_messages, src)
	if err != nil {
		log.Error().Err(err).Str("dst", dst).Str("src", src).Msg("Failed to generate next message")
		return
	}
	err = sendText(ctx, dst, src, next_message.Content, enums.CommsTextoriginLLM, false, true)
	if err != nil {
		log.Error().Err(err).Str("src", src).Str("dst", dst).Str("content", next_message.Content).Msg("Failed to send response text")
		return
	}
	log.Info().Str("src", src).Str("dst", dst).Str("body", body).Str("reply", next_message.Content).Msg("Handled text message")
}

func ParsePhoneNumber(input string) (*E164, error) {
	return phonenumbers.Parse(input, "US")
}

func StoreSources() error {
	ctx := context.TODO()
	for _, n := range []string{config.PhoneNumberReportStr, config.PhoneNumberSupportStr, config.VoipMSNumber} {
		var err error
		// Deal with Voip.ms not expecting API calls with the prefixed +1
		if !strings.HasPrefix(n, "+1") {
			err = ensureInDB(ctx, "+1"+n)
		} else {
			err = ensureInDB(ctx, n)
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
func delayMessage(ctx context.Context, source enums.CommsTextjobsource, destination string, content string, type_ enums.CommsTextjobtype) error {
	job, err := models.CommsTextJobs.Insert(&models.CommsTextJobSetter{
		Content:     omit.From(content),
		Created:     omit.From(time.Now()),
		Destination: omit.From(destination),
		//ID:
		Source: omit.From(source),
		Type:   omit.From(type_),
	}).One(ctx, db.PGInstance.BobDB)
	if err != nil {
		return fmt.Errorf("Failed to add delayed text job: %w", err)
	}
	log.Info().Int32("id", job.ID).Msg("Created delayed text job")
	return nil
}

func resendInitialText(ctx context.Context, src string, dst string) error {
	phone, err := models.FindCommsPhone(ctx, db.PGInstance.BobDB, dst)
	if err != nil {
		return fmt.Errorf("Failed to find phone %s: %w", dst, err)
	}
	err = phone.Update(ctx, db.PGInstance.BobDB, &models.CommsPhoneSetter{
		IsSubscribed: omitnull.FromPtr[bool](nil),
	})
	if err != nil {
		return fmt.Errorf("Failed to clear subscription on phone %s: %w", dst, err)
	}
	return nil
}

func sendInitialText(ctx context.Context, src string, dst string) error {
	content := "Welcome to Report Mosquitoes Online. We received your request and want to confirm text updates. Reply YES to continue. Reply STOP at any time to unsubscribe"
	origin := enums.CommsTextoriginWebsiteAction
	err := sendText(ctx, src, dst, content, origin, true, true)
	if err != nil {
		return fmt.Errorf("Failed to send initial confirmation: %w", err)
	}
	return nil
}

func ensureInitialText(ctx context.Context, src string, dst string) error {
	//
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
	return sendInitialText(ctx, src, dst)
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

func generateNextMessage(ctx context.Context, history []llm.Message, customer_phone string) (llm.Message, error) {
	_handle_report_status := func() (string, error) {
		return "Report: ABCD-1234-5678, District: Delta MVCD, Status: scheduled, Appointment: Wednesday 3:30pm", nil
	}
	_handle_contact_district := func(reason string) {
		log.Warn().Str("reason", reason).Msg("Contacting district")
	}
	_handle_contact_supervisor := func(reason string) {
		log.Warn().Str("reason", reason).Msg("Contacting supervisor")
	}
	return llm.GenerateNextMessage(ctx, history, _handle_report_status, _handle_contact_district, _handle_contact_supervisor)
}

func handleWaitingTextJobs(ctx context.Context, src string) {
	jobs, err := models.CommsTextJobs.Query(
		models.SelectWhere.CommsTextJobs.Destination.EQ(src),
		models.SelectWhere.CommsTextJobs.Completed.IsNull(),
	).All(ctx, db.PGInstance.BobDB)
	if err != nil {
		log.Error().Err(err).Msg("Failed to query for jobs")
		return
	}
	for _, job := range jobs {
		var src string
		switch job.Source {
		case enums.CommsTextjobsourceRmo:
			src = config.PhoneNumberReportStr
		//case enums.CommsTextJobsourcenidus:
		//src := config.PhoneNumebrNidusStr
		default:
			log.Error().Str("source", job.Source.String()).Msg("Can't support background text job.")
		}
		err = sendText(ctx, src, job.Destination, job.Content, enums.CommsTextoriginWebsiteAction, false, true)
		if err != nil {
			log.Error().Err(err).Int32("id", job.ID).Msg("Failed to send delayed text job.")
			continue
		}
		err := job.Update(ctx, db.PGInstance.BobDB, &models.CommsTextJobSetter{
			Completed: omitnull.From(time.Now()),
		})
		if err != nil {
			log.Error().Err(err).Int32("id", job.ID).Msg("Failed to update delayed text job.")
			continue
		}
	}
}

func handleResetConversation(ctx context.Context, src string, dst string) {
	err := wipeLLMMemory(ctx, src, dst)
	if err != nil {
		log.Error().Err(err).Str("src", src).Str("dst", dst).Msg("Failed to wipe memory")
		content := "Failed to wip memory"
		err = sendText(ctx, dst, src, content, enums.CommsTextoriginCommandResponse, false, false)
		if err != nil {
			log.Error().Err(err).Str("src", src).Str("dst", dst).Msg("Failed to indicated memory wipe failure.")
		}
		return
	}
	content := "LLM memory wiped"
	err = sendText(ctx, dst, src, content, enums.CommsTextoriginCommandResponse, false, false)
	if err != nil {
		log.Error().Err(err).Str("src", src).Str("dst", dst).Msg("Failed to indicated memory wiped.")
		return
	}
	log.Info().Err(err).Str("src", src).Str("dst", dst).Msg("Wiped LLM memory")
}

func insertTextLog(ctx context.Context, content string, destination string, source string, origin enums.CommsTextorigin, is_welcome bool, is_visible_to_llm bool) (log *models.CommsTextLog, err error) {
	log, err = models.CommsTextLogs.Insert(&models.CommsTextLogSetter{
		//ID:
		Content:        omit.From(content),
		Created:        omit.From(time.Now()),
		Destination:    omit.From(destination),
		IsVisibleToLLM: omit.From(is_visible_to_llm),
		IsWelcome:      omit.From(is_welcome),
		Origin:         omit.From(origin),
		Source:         omit.From(source),
		TwilioSid:      omitnull.FromPtr[string](nil),
		TwilioStatus:   omit.From(""),
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

func loadPreviousMessagesForLLM(ctx context.Context, dst, src string) ([]llm.Message, error) {
	messages, err := sql.TextsBySenders(dst, src).All(ctx, db.PGInstance.BobDB)
	results := make([]llm.Message, 0)
	if err != nil {
		return results, fmt.Errorf("Failed to get message history for %s and %s: %w", dst, src, err)
	}
	for _, m := range messages {
		if m.IsVisibleToLLM {
			is_from_customer := (m.Source == src)
			results = append(results, llm.Message{
				IsFromCustomer: is_from_customer,
				Content:        m.Content,
			})
		}
	}
	return results, nil
}

func sendText(ctx context.Context, source string, destination string, message string, origin enums.CommsTextorigin, is_welcome bool, is_visible_to_llm bool) error {
	err := ensureInDB(ctx, destination)
	if err != nil {
		return fmt.Errorf("Failed to ensure text message destination is in the DB: %w", err)
	}
	l, err := insertTextLog(ctx, message, destination, source, origin, is_welcome, is_visible_to_llm)
	if err != nil {
		return fmt.Errorf("Failed to insert text message in the DB: %w", err)
	}
	sid, err := text.SendText(ctx, source, destination, message)
	if err != nil {
		return fmt.Errorf("Failed to send text message: %w", err)
	}
	err = l.Update(ctx, db.PGInstance.BobDB, &models.CommsTextLogSetter{
		TwilioSid:    omitnull.From(sid),
		TwilioStatus: omit.From("created"),
	})
	if err != nil {
		return fmt.Errorf("Failed to update text Twilio status: %w", err)
	}
	log.Info().Int32("id", l.ID).Bool("is_visible_to_llm", is_visible_to_llm).Str("message", message).Msg("inserted text log")

	return nil
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

func wipeLLMMemory(ctx context.Context, src string, dst string) error {
	rows, err := sql.TextsBySenders(dst, src).All(ctx, db.PGInstance.BobDB)
	if err != nil {
		return fmt.Errorf("Failed to query for texts: %w", err)
	}
	ids := make([]int32, 0)
	for _, r := range rows {
		ids = append(ids, r.ID)
	}
	_, err = models.CommsTextLogs.Update(
		um.Where(
			models.CommsTextLogs.Columns.ID.EQ(psql.Any(ids)),
		),
		um.SetCol("is_visible_to_llm").ToArg(false),
	).Exec(ctx, db.PGInstance.BobDB)
	if err != nil {
		return fmt.Errorf("Failed to update texts: %w", err)
	}

	return nil
}
