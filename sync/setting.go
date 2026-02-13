package sync

import (
	"net/http"

	"github.com/Gleipnir-Technology/nidus-sync/db/models"
	"github.com/Gleipnir-Technology/nidus-sync/html"
)

func getSettings(w http.ResponseWriter, r *http.Request, u *models.User) {
	userContent, err := contentForUser(r.Context(), u)
	if err != nil {
		respondError(w, "Failed to get user content", err, http.StatusInternalServerError)
		return
	}
	data := ContentAuthenticatedPlaceholder{
		URL:  newContentURL(),
		User: userContent,
	}
	html.RenderOrError(w, "sync/settings.html", data)
}
