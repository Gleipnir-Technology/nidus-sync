package sync

import (
	"context"
	"fmt"
	"net/http"

	"github.com/Gleipnir-Technology/nidus-sync/comms/email"
	"github.com/Gleipnir-Technology/nidus-sync/comms/text"
	"github.com/Gleipnir-Technology/nidus-sync/config"
	"github.com/Gleipnir-Technology/nidus-sync/db/enums"
	"github.com/Gleipnir-Technology/nidus-sync/db/models"
	"github.com/gorilla/schema"
	"github.com/rs/zerolog/log"
)

type contentSudo struct {
	ForwardEmailRMOAddress   string
	ForwardEmailNidusAddress string
}

func getSudo(ctx context.Context, user *models.User) (string, interface{}, *errorWithStatus) {
	if user.Role != enums.UserroleRoot {
		return "", nil, &errorWithStatus{
			Message: "You have to be a root user to access this",
			Status:  http.StatusForbidden,
		}
	}
	content := contentSudo{
		ForwardEmailRMOAddress:   config.ForwardEmailRMOAddress,
		ForwardEmailNidusAddress: config.ForwardEmailNidusAddress,
	}
	return "sync/sudo.html", content, nil
}

var decoder = schema.NewDecoder()

type FormEmail struct {
	Body    string `schema:"emailBody"`
	From    string `schema:"emailFrom"`
	Subject string `schema:"emailSubject"`
	To      string `schema:"emailTo"`
}

func postSudoEmail(ctx context.Context, u *models.User, e FormEmail) (string, *errorWithStatus) {
	if u.Role != enums.UserroleRoot {
		return "", &errorWithStatus{
			Message: "You must have sudo powers to do this",
			Status:  http.StatusForbidden,
		}
	}
	request := email.Request{
		From:    e.From,
		HTML:    fmt.Sprintf("<html><p>%s</p></html>", e.Body),
		Sender:  e.From,
		Subject: e.Subject,
		To:      e.To,
		Text:    e.Body,
	}
	resp, err := email.Send(ctx, request)
	if err != nil {
		log.Warn().Err(err).Msg("Failed to send email")
	} else {
		log.Info().Str("id", resp.ID).Msg("Sent Email...?")
	}
	return "/sudo", nil
}

type FormSMS struct {
	Message string `schema:"smsMessage"`
	Phone   string `schema:"smsPhone"`
}

func postSudoSMS(ctx context.Context, u *models.User, sms FormSMS) (string, *errorWithStatus) {
	if u.Role != enums.UserroleRoot {
		return "", &errorWithStatus{
			Message: "You must have sudo powers to do this",
			Status:  http.StatusForbidden,
		}
	}
	id, err := text.SendText(ctx, config.VoipMSNumber, sms.Phone, sms.Message)
	if err != nil {
		log.Warn().Err(err).Msg("Failed to send SMS")
	} else {
		log.Info().Str("id", id).Msg("Sent SMS")
	}
	return "/sudo", nil
}
