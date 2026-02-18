package sync

import (
	"context"
	"net/http"

	"github.com/Gleipnir-Technology/nidus-sync/comms/text"
	"github.com/Gleipnir-Technology/nidus-sync/config"
	"github.com/Gleipnir-Technology/nidus-sync/db/enums"
	"github.com/Gleipnir-Technology/nidus-sync/db/models"
	"github.com/gorilla/schema"
	"github.com/rs/zerolog/log"
)

type contentSudo struct{}

func getSudo(ctx context.Context, user *models.User) (string, interface{}, *errorWithStatus) {
	if user.Role != enums.UserroleRoot {
		return "", nil, &errorWithStatus{
			Message: "You have to be a root user to access this",
			Status:  http.StatusForbidden,
		}
	}
	content := contentAdminDash{}
	return "sync/sudo.html", content, nil
}

var decoder = schema.NewDecoder()

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
