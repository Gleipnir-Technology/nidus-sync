package sync

import (
	"net/http"

	"github.com/Gleipnir-Technology/nidus-sync/auth"
	"github.com/Gleipnir-Technology/nidus-sync/background"
	"github.com/Gleipnir-Technology/nidus-sync/config"
	"github.com/Gleipnir-Technology/nidus-sync/db/models"
	"github.com/Gleipnir-Technology/nidus-sync/htmlpage"
	"github.com/rs/zerolog/log"
)

var (
	oauthPromptT = buildTemplate("oauth-prompt", "authenticated")
)

func getArcgisOauthBegin(w http.ResponseWriter, r *http.Request) {
	authURL := config.BuildArcGISAuthURL(config.ClientID)
	http.Redirect(w, r, authURL, http.StatusFound)
}

func getArcgisOauthCallback(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")
	log.Info().Str("code", code).Msg("Handling oauth callback")
	if code == "" {
		respondError(w, "Access code is empty", nil, http.StatusBadRequest)
		return
	}
	user, err := auth.GetAuthenticatedUser(r)
	if err != nil {
		respondError(w, "You're not currently authenticated, which really shouldn't happen.", err, http.StatusUnauthorized)
		return
	}
	err = background.HandleOauthAccessCode(r.Context(), user, code)
	if err != nil {
		respondError(w, "Failed to handle access code", err, http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, config.MakeURLSync("/"), http.StatusFound)
}

func getOAuthRefresh(w http.ResponseWriter, r *http.Request) {
	user, err := auth.GetAuthenticatedUser(r)
	if err != nil {
		http.Redirect(w, r, "/?next=/oauth/refresh", http.StatusFound)
		return
	}
	oauthPrompt(w, user)
}

func oauthPrompt(w http.ResponseWriter, user *models.User) {
	dp := user.DisplayName
	data := ContentDashboard{
		User: User{
			DisplayName: dp,
			Initials:    extractInitials(dp),
			Username:    user.Username,
		},
	}
	htmlpage.RenderOrError(w, oauthPromptT, data)
}
