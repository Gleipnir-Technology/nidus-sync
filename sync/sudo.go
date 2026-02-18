package sync

import (
	"context"
	"net/http"

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

func postSudoSMS(w http.ResponseWriter, r *http.Request, u *models.User) {
	err := r.ParseForm()
	if err != nil {
		respondError(w, "Failed to parse form", err, http.StatusInternalServerError)
		return
	}

	var sms FormSMS

	err = decoder.Decode(&sms, r.PostForm)
	if err != nil {
		respondError(w, "Failed to decode form", err, http.StatusInternalServerError)
		return
	}
	log.Info().Str("msg", sms.Message).Str("phone", sms.Phone).Msg("Got SMS")
	http.Redirect(w, r, "/sudo", http.StatusFound)
}
