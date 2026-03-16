package email

import (
	"context"
	"fmt"

	"github.com/Gleipnir-Technology/nidus-sync/config"
	"github.com/Gleipnir-Technology/nidus-sync/db"
	"github.com/Gleipnir-Technology/nidus-sync/db/models"
	//"github.com/rs/zerolog/log"
)

type contentEmailInitial struct {
	Base         contentEmailBase
	Destination  string
	URLSubscribe string
}

func maybeSendInitialEmail(ctx context.Context, destination string) error {
	err := EnsureInDB(ctx, destination)
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
func urlEmailInBrowser(public_id string) string {
	return config.MakeURLReport("/email/render/%s", public_id)
}
func urlUnsubscribe(email string) string {
	return config.MakeURLReport("/email/unsubscribe?email=%s", email)
}
func sendEmailInitialContact(ctx context.Context, destination string) error {
	//data := pgtypes.HStore{}
	data := make(map[string]string, 0)
	source := config.ForwardEmailRMOAddress
	data["Destination"] = destination
	data["Source"] = source
	data["URLLogo"] = config.MakeURLReport("/static/img/nidus-logo-no-lettering-64.png")
	data["URLSubscribe"] = config.MakeURLReport("/email/confirm?email=%s", destination)
	data["URLUnsubscribe"] = urlUnsubscribe(destination)

	subject := "Welcome"
	err := sendEmailBegin(ctx, source, destination, templateInitialID, subject, data)
	if err != nil {
		return fmt.Errorf("Failed to send initial email to %s: %w", err)
	}
	return nil
}
