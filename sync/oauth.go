package sync

import (
	"net/http"
	"net/url"
	"strconv"

	"github.com/Gleipnir-Technology/nidus-sync/auth"
	"github.com/Gleipnir-Technology/nidus-sync/background"
	"github.com/Gleipnir-Technology/nidus-sync/config"
	"github.com/Gleipnir-Technology/nidus-sync/db/models"
	"github.com/Gleipnir-Technology/nidus-sync/html"
	"github.com/rs/zerolog/log"
)

type ContextOauthPrompt struct {
	User User
}

// Build the ArcGIS authorization URL with PKCE
func buildArcGISAuthURL(clientID string) string {
	baseURL := "https://www.arcgis.com/sharing/rest/oauth2/authorize/"

	params := url.Values{}
	params.Add("client_id", clientID)
	params.Add("redirect_uri", config.ArcGISOauthRedirectURL())
	params.Add("response_type", "code")
	//params.Add("code_challenge", generateCodeChallenge(codeVerifier))
	//params.Add("code_challenge_method", "S256")

	// See https://developers.arcgis.com/rest/users-groups-and-items/token/
	// expiration is defined in minutes
	var expiration int
	if config.IsProductionEnvironment() {
		// 2 weeks is the maximum allowed
		expiration = 20160
	} else {
		expiration = 20
	}
	params.Add("expiration", strconv.Itoa(expiration))

	return baseURL + "?" + params.Encode()
}

func getArcgisOauthBegin(w http.ResponseWriter, r *http.Request) {
	authURL := buildArcGISAuthURL(config.ClientID)
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
	http.Redirect(w, r, config.MakeURLNidus("/"), http.StatusFound)
}

func getOAuthRefresh(w http.ResponseWriter, r *http.Request) {
	user, err := auth.GetAuthenticatedUser(r)
	if err != nil {
		http.Redirect(w, r, "/?next=/oauth/refresh", http.StatusFound)
		return
	}
	oauthPrompt(w, r, user)
}

func oauthPrompt(w http.ResponseWriter, r *http.Request, user *models.User) {
	userContent, err := contentForUser(r.Context(), user)
	if err != nil {
		respondError(w, "Failed to get user content", err, http.StatusInternalServerError)
		return
	}
	data := ContextOauthPrompt{
		User: userContent,
	}
	html.RenderOrError(w, "sync/oauth-prompt.html", data)
}
