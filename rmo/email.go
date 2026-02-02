package rmo

import (
	"net/http"

	"github.com/Gleipnir-Technology/nidus-sync/comms/email"
	"github.com/Gleipnir-Technology/nidus-sync/db"
	"github.com/Gleipnir-Technology/nidus-sync/db/models"
	"github.com/go-chi/chi/v5"
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
