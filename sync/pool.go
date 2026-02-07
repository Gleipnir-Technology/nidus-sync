package sync

import (
	"github.com/Gleipnir-Technology/nidus-sync/db/models"
	"github.com/Gleipnir-Technology/nidus-sync/html"
	"net/http"
)

type ContentPoolList struct {
	URL  ContentURL
	User User
}

func getPoolList(w http.ResponseWriter, r *http.Request, u *models.User) {
	userContent, err := contentForUser(r.Context(), u)
	if err != nil {
		respondError(w, "Failed to get user", err, http.StatusInternalServerError)
		return
	}
	data := ContentPoolList{
		URL:  newContentURL(),
		User: userContent,
	}
	html.RenderOrError(w, "sync/pool-list.html", data)
}
