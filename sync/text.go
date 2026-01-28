package sync

import (
	"net/http"

	"github.com/Gleipnir-Technology/nidus-sync/db/models"
	"github.com/Gleipnir-Technology/nidus-sync/htmlpage"
)

type ContentTextMessages struct {
	User User
}

var (
	textMessagesT = buildTemplate("text-messages", "authenticated")
)

func getTextMessages(w http.ResponseWriter, r *http.Request, u *models.User) {
	userContent, err := contentForUser(r.Context(), u)
	if err != nil {
		respondError(w, "Failed to get user", err, http.StatusInternalServerError)
		return
	}
	content := ContentTextMessages{
		User: userContent,
	}
	htmlpage.RenderOrError(w, textMessagesT, content)
}
