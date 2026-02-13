package sync

import (
	"net/http"

	"github.com/Gleipnir-Technology/nidus-sync/arcgis"
	"github.com/Gleipnir-Technology/nidus-sync/db/models"
	"github.com/Gleipnir-Technology/nidus-sync/html"
)

type ContentSettingIntegration struct {
	ArcGISOAuth *models.OauthToken
	URL         ContentURL
	User        User
}

func getSetting(w http.ResponseWriter, r *http.Request, u *models.User) {
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
func getSettingIntegration(w http.ResponseWriter, r *http.Request, u *models.User) {
	ctx := r.Context()
	userContent, err := contentForUser(ctx, u)
	if err != nil {
		respondError(w, "Failed to get user content", err, http.StatusInternalServerError)
		return
	}
	oauth, err := arcgis.GetOAuthForUser(ctx, u)
	if err != nil {
		respondError(w, "Failed to get oauth", err, http.StatusInternalServerError)
		return
	}
	data := ContentSettingIntegration{
		ArcGISOAuth: oauth,
		URL:         newContentURL(),
		User:        userContent,
	}
	html.RenderOrError(w, "sync/setting-integration.html", data)
}
