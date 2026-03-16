package email

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/Gleipnir-Technology/bob"
	"github.com/Gleipnir-Technology/nidus-sync/comms/email"
	"github.com/Gleipnir-Technology/nidus-sync/config"
	"github.com/Gleipnir-Technology/nidus-sync/db"
	"github.com/Gleipnir-Technology/nidus-sync/db/enums"
	"github.com/Gleipnir-Technology/nidus-sync/db/models"
	"github.com/Gleipnir-Technology/nidus-sync/platform/background"
	"github.com/aarondl/opt/omit"
	"github.com/aarondl/opt/omitnull"
	"github.com/rs/zerolog/log"
)

func EnsureInDB(ctx context.Context, destination string) (err error) {
	_, err = models.FindCommsEmailContact(ctx, db.PGInstance.BobDB, destination)
	if err != nil {
		// doesn't exist
		if err.Error() == "sql: no rows in result set" {
			public_id := fmt.Sprintf("%x", sha256.Sum256([]byte(destination)))
			_, err = models.CommsEmailContacts.Insert(&models.CommsEmailContactSetter{
				Address:      omit.From(destination),
				Confirmed:    omit.From(false),
				IsSubscribed: omit.From(false),
				PublicID:     omit.From(public_id),
			}).One(ctx, db.PGInstance.BobDB)
			if err != nil {
				return fmt.Errorf("Failed to insert new email: %w", err)
			}
			log.Info().Str("email", destination).Msg("Added email to the comms database")
			return nil
		}
		return fmt.Errorf("Unexpected error searching for contact: %w", err)
	}
	return nil
}

func insertEmailLog(ctx context.Context, data map[string]string, destination string, public_id string, source string, subject string, template_id int32) (email_id *int32, err error) {
	data_for_insert := db.ConvertToPGData(data)
	var type_ enums.CommsMessagetypeemail
	switch template_id {
	case templateReportNotificationConfirmationID:
		type_ = enums.CommsMessagetypeemailReportNotificationConfirmation
	case templateInitialID:
		type_ = enums.CommsMessagetypeemailInitialContact
	default:
		return nil, fmt.Errorf("Unrecognized template ID %d", template_id)
	}
	e, err := models.CommsEmailLogs.Insert(&models.CommsEmailLogSetter{
		//ID:
		Created:        omit.From(time.Now()),
		DeliveryStatus: omit.From("initial"),
		Destination:    omit.From(destination),
		PublicID:       omit.From(public_id),
		SentAt:         omitnull.FromPtr[time.Time](nil),
		Source:         omit.From(source),
		Subject:        omit.From(subject),
		TemplateID:     omit.From(template_id),
		TemplateData:   omit.From(data_for_insert),
		Type:           omit.From(type_),
	}).One(ctx, db.PGInstance.BobDB)
	if err != nil {
		return nil, fmt.Errorf("insern email log: %w", err)
	}
	return &e.ID, nil
}
func generatePublicId(template int32, m map[string]string) string {
	if m == nil || len(m) == 0 {
		// Return hash of empty string for empty maps
		emptyHash := sha256.Sum256([]byte(""))
		return hex.EncodeToString(emptyHash[:])
	}

	// Get and sort keys for deterministic ordering
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	// Build a string with all key-value pairs
	var sb strings.Builder
	// Add type first
	sb.WriteString(fmt.Sprintf("template:%d,", template))
	for _, k := range keys {
		sb.WriteString(k)
		sb.WriteString(":") // Separator between key and value
		sb.WriteString(m[k])
		sb.WriteString(",") // Separator between pairs
	}

	// Compute SHA-256 hash
	hasher := sha256.New()
	hasher.Write([]byte(sb.String()))
	hashBytes := hasher.Sum(nil)

	// Convert to hex string and return
	return hex.EncodeToString(hashBytes)
}
func sendEmailBegin(ctx context.Context, source string, destination string, template int32, subject string, data map[string]string) error {
	public_id := generatePublicId(template, data)
	data["URLViewInBrowser"] = urlEmailInBrowser(public_id)

	e, err := insertEmailLog(ctx, data, destination, public_id, config.ForwardEmailRMOAddress, subject, template)
	if err != nil {
		return fmt.Errorf("Failed to store email log: %w", err)
	}
	return background.NewEmailSend(ctx, db.PGInstance.BobDB, *e)
}
func sendEmailComplete(ctx context.Context, txn bob.Executor, email_id int32) error {
	email_log, err := models.FindCommsEmailLog(ctx, txn, email_id)
	if err != nil {
		return fmt.Errorf("find email: %w", err)
	}
	data := db.ConvertFromPGData(email_log.TemplateData)
	text, html, err := renderEmailTemplates(email_log.TemplateID, data)
	if err != nil {
		return fmt.Errorf("Failed to render email report notification template: %w", err)
	}
	resp, err := email.Send(ctx, email.Request{
		From:    config.ForwardEmailRMOAddress,
		HTML:    html,
		Subject: email_log.Subject,
		Text:    text,
		To:      email_log.Destination,
	})
	if err != nil {
		return fmt.Errorf("Failed to send email %d: %w", email_log.ID, err)
	}
	log.Info().Str("response id", resp.ID).Int32("email id", email_log.ID).Msg("Sent email")
	return nil
}
