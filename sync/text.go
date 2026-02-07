package sync

import (
	"net/http"

	"github.com/Gleipnir-Technology/nidus-sync/db/models"
	"github.com/Gleipnir-Technology/nidus-sync/html"
)

type ContentTextMessages struct {
	User User
}

func getTextMessages(w http.ResponseWriter, r *http.Request, u *models.User) {
	userContent, err := contentForUser(r.Context(), u)
	if err != nil {
		respondError(w, "Failed to get user", err, http.StatusInternalServerError)
		return
	}
	content := ContentTextMessages{
		User: userContent,
	}
	html.RenderOrError(w, "sync/text-messages.html", content)
}
