package sync

import (
	"net/http"

	"github.com/Gleipnir-Technology/nidus-sync/db/models"
	"github.com/Gleipnir-Technology/nidus-sync/html"
)

type contentRadar struct {
	URL  ContentURL
	User User
}

func getRadar(w http.ResponseWriter, r *http.Request, u *models.User) {
	ctx := r.Context()
	userContent, err := contentForUser(ctx, u)
	if err != nil {
		respondError(w, "Failed to get user", err, http.StatusInternalServerError)
		return
	}
	data := contentRadar{
		URL:  newContentURL(),
		User: userContent,
	}
	html.RenderOrError(w, "sync/radar.html", data)
}
