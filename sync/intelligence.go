package sync

import (
	"context"
	"net/http"

	"github.com/Gleipnir-Technology/nidus-sync/background"
	"github.com/Gleipnir-Technology/nidus-sync/db/models"
	"github.com/Gleipnir-Technology/nidus-sync/html"
	nhttp "github.com/Gleipnir-Technology/nidus-sync/http"
	"github.com/rs/zerolog/log"
)

type contentIntelligenceRoot struct {
	ArcgisAccessToken string
}

func getIntelligenceRoot(ctx context.Context, r *http.Request, org *models.Organization, user *models.User) (*html.Response[contentIntelligenceRoot], *nhttp.ErrorWithStatus) {
	var oauth_token *models.ArcgisOauthToken
	var err error
	var access_token string
	oauth_token, err = background.GetOAuthForOrg(ctx, org)
	if err != nil {
		log.Warn().Err(err).Msg("Failed to get oauth")
		oauth_token = nil
		access_token = ""
	} else {
		access_token = oauth_token.AccessToken
	}
	return html.NewResponse("sync/intelligence-root.html", contentIntelligenceRoot{
		ArcgisAccessToken: access_token,
	}), nil
}
