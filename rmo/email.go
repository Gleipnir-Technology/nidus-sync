package rmo

import (
	"net/http"

	"github.com/Gleipnir-Technology/nidus-sync/comms/email"
	"github.com/Gleipnir-Technology/nidus-sync/db"
	"github.com/Gleipnir-Technology/nidus-sync/db/models"
	"github.com/Gleipnir-Technology/nidus-sync/html"
	"github.com/aarondl/opt/omit"
	"github.com/go-chi/chi/v5"
)

type ContentEmail struct {
	Email string
}

var (
	EmailConfirmT             = buildTemplate("email-confirm", "base")
	EmailConfirmCompleteT     = buildTemplate("email-confirm-complete", "base")
	EmailUnsubscribeT         = buildTemplate("email-unsubscribe", "base")
	EmailUnsubscribeCompleteT = buildTemplate("email-unsubscribe-complete", "base")
)

func getEmailByCode(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "code")
	//id := r.FormValue("id")
	if id == "" {
		http.Error(w, "You must specify an id", http.StatusBadRequest)
		return
	}
	ctx := r.Context()
	email_log, err := models.CommsEmailLogs.Query(
		models.SelectWhere.CommsEmailLogs.PublicID.EQ(id),
	).One(ctx, db.PGInstance.BobDB)
	if err != nil {
		respondError(w, "Failed to query email_log: %w", err, http.StatusInternalServerError)
		return
	}
	html, err := email.RenderHTML(email_log.TemplateID, email_log.TemplateData)
	if err != nil {
		respondError(w, "Failed to render email_log: %w", err, http.StatusInternalServerError)
		return
	}
	w.Write(html)
}
func getEmailReportUnsubscribe(w http.ResponseWriter, r *http.Request) {
	email := r.FormValue("email")
	html.RenderOrError(
		w,
		EmailConfirmT,
		ContentEmail{
			Email: email,
		},
	)
}
func getEmailConfirm(w http.ResponseWriter, r *http.Request) {
	email := r.FormValue("email")
	if email == "" {
		respondError(w, "Not sure what to do with an empty email", nil, http.StatusBadRequest)
		return
	}

	html.RenderOrError(
		w,
		EmailConfirmT,
		ContentEmail{
			Email: email,
		},
	)
}
func getEmailConfirmComplete(w http.ResponseWriter, r *http.Request) {
	html.RenderOrError(
		w,
		EmailConfirmCompleteT,
		map[string]string{},
	)
}
func getEmailUnsubscribe(w http.ResponseWriter, r *http.Request) {
	email := r.FormValue("email")
	html.RenderOrError(
		w,
		EmailUnsubscribeT,
		ContentEmail{
			Email: email,
		},
	)
}
func getEmailUnsubscribeComplete(w http.ResponseWriter, r *http.Request) {
	html.RenderOrError(
		w,
		EmailUnsubscribeCompleteT,
		map[string]string{},
	)
}
func postEmailConfirm(w http.ResponseWriter, r *http.Request) {
	email := r.PostFormValue("email")
	if email == "" {
		respondError(w, "Not sure what to do with an empty email", nil, http.StatusBadRequest)
		return
	}
	ctx := r.Context()
	email_contact, err := models.FindCommsEmailContact(ctx, db.PGInstance.BobDB, email)
	if err != nil {
		respondError(w, "Email not in the database", err, http.StatusNotFound)
		return
	}
	err = email_contact.Update(ctx, db.PGInstance.BobDB, &models.CommsEmailContactSetter{
		Confirmed: omit.From(true),
	})
	http.Redirect(w, r, "/email/confirm/complete", http.StatusFound)
}
func postEmailUnsubscribe(w http.ResponseWriter, r *http.Request) {
	email := r.PostFormValue("email")
	if email == "" {
		respondError(w, "Not sure what to do with an empty email", nil, http.StatusBadRequest)
		return
	}
	ctx := r.Context()
	email_contact, err := models.FindCommsEmailContact(ctx, db.PGInstance.BobDB, email)
	if err != nil {
		respondError(w, "Email not in the database", err, http.StatusNotFound)
		return
	}
	err = email_contact.Update(ctx, db.PGInstance.BobDB, &models.CommsEmailContactSetter{
		IsSubscribed: omit.From(false),
	})
	http.Redirect(w, r, "/email/unsubscribe/complete", http.StatusFound)
}
