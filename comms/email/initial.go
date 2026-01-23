package email

import (
	"context"
	"fmt"

	"github.com/Gleipnir-Technology/nidus-sync/config"
	"github.com/Gleipnir-Technology/nidus-sync/db"
	"github.com/Gleipnir-Technology/nidus-sync/db/enums"
	"github.com/Gleipnir-Technology/nidus-sync/db/models"
	"github.com/rs/zerolog/log"
)

type jobInitial struct {
	base jobEmailBase
}

func (job jobInitial) Destination() string {
	return job.base.destination
}

func maybeSendInitialEmail(ctx context.Context, destination string) error {
	err := ensureInDB(ctx, destination)
	if err != nil {
		return fmt.Errorf("Failed to add email recipient to database: %w", err)
	}
	rows, err := models.CommsEmailLogs.Query(
		models.SelectWhere.CommsEmailLogs.Destination.EQ(destination),
		models.SelectWhere.CommsEmailLogs.TemplateID.EQ(templateInitialID),
	).All(ctx, db.PGInstance.BobDB)

	// We already sent an initial email
	if len(rows) > 0 {
		return nil
	}

	return sendEmailInitialContact(ctx, destination)
}
func sendEmailInitialContact(ctx context.Context, destination string) error {
	//data := pgtypes.HStore{}
	data := make(map[string]string, 0)
	source := config.ForwardEmailReportAddress
	data["destination"] = destination
	data["source"] = source
	data["url_logo"] = config.MakeURLReport("/static/img/nidus-logo-no-lettering-64.png")
	data["url_subscribe"] = config.MakeURLReport("/email/subscribe?email=%s", destination)
	data["url_unsubscribe"] = config.MakeURLReport("/email/unsubscribe")
	public_id := generatePublicId(enums.CommsMessagetypeemailInitialContact, data)
	data["url_browser"] = config.MakeURLReport("/email?id=%s", public_id)

	text, html, err := renderEmailTemplates(templateInitialID, data)
	if err != nil {
		return fmt.Errorf("Failed to render email temlates: %w", err)
	}

	subject := "Welcome"
	err = insertEmailLog(ctx, data, destination, public_id, source, subject, templateInitialID)
	if err != nil {
		return fmt.Errorf("Failed to store email log: %w", err)
	}
	resp, err := sendEmail(ctx, emailRequest{
		From:    source,
		HTML:    html,
		Subject: subject,
		Text:    text,
		To:      destination,
	}, enums.CommsMessagetypeemailInitialContact)

	if err != nil {
		return fmt.Errorf("Failed to send email to %s: %w", err)
	}
	log.Info().Str("id", resp.ID).Str("to", destination).Msg("Sent initial contact email")
	return nil
}
