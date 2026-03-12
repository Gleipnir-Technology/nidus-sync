package sync

import (
	"context"
	"fmt"
	"net/http"

	"github.com/Gleipnir-Technology/nidus-sync/comms/email"
	"github.com/Gleipnir-Technology/nidus-sync/comms/text"
	"github.com/Gleipnir-Technology/nidus-sync/config"
	"github.com/Gleipnir-Technology/nidus-sync/html"
	nhttp "github.com/Gleipnir-Technology/nidus-sync/http"
	"github.com/Gleipnir-Technology/nidus-sync/platform"
	"github.com/rs/zerolog/log"
)

type contentSudo struct {
	ForwardEmailRMOAddress   string
	ForwardEmailNidusAddress string
}

func getSudo(ctx context.Context, r *http.Request, user platform.User) (*html.Response[contentSudo], *nhttp.ErrorWithStatus) {
	if !user.HasRoot() {
		return nil, &nhttp.ErrorWithStatus{
			Message: "You have to be a root user to access this",
			Status:  http.StatusForbidden,
		}
	}
	content := contentSudo{
		ForwardEmailRMOAddress:   config.ForwardEmailRMOAddress,
		ForwardEmailNidusAddress: config.ForwardEmailNidusAddress,
	}
	return html.NewResponse("sync/sudo.html", content), nil
}

type FormEmail struct {
	Body    string `schema:"emailBody"`
	From    string `schema:"emailFrom"`
	Subject string `schema:"emailSubject"`
	To      string `schema:"emailTo"`
}

func postSudoEmail(ctx context.Context, r *http.Request, u platform.User, e FormEmail) (string, *nhttp.ErrorWithStatus) {
	if !u.HasRoot() {
		return "", &nhttp.ErrorWithStatus{
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
		log.Info().Str("id", resp.ID).Str("to", e.To).Msg("Sent Email")
	}
	return "/sudo", nil
}

type FormSMS struct {
	Message string `schema:"smsMessage"`
	Phone   string `schema:"smsPhone"`
}

func postSudoSMS(ctx context.Context, r *http.Request, u platform.User, sms FormSMS) (string, *nhttp.ErrorWithStatus) {
	if !u.HasRoot() {
		return "", &nhttp.ErrorWithStatus{
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
