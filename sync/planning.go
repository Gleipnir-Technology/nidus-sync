package sync

import (
	"context"
	"net/http"

	"github.com/Gleipnir-Technology/nidus-sync/db/models"
	"github.com/Gleipnir-Technology/nidus-sync/html"
	nhttp "github.com/Gleipnir-Technology/nidus-sync/http"
	"github.com/Gleipnir-Technology/nidus-sync/platform"
	"github.com/rs/zerolog/log"
)

type contentPlanningRoot struct {
	ArcgisAccessToken string
}

func getPlanningRoot(ctx context.Context, r *http.Request, user platform.User) (*html.Response[contentPlanningRoot], *nhttp.ErrorWithStatus) {
	var oauth_token *models.ArcgisOauthToken
	var err error
	var access_token string
	oauth_token, err = platform.GetOAuthForOrg(ctx, user.Organization)
	if err != nil {
		log.Warn().Err(err).Msg("Failed to get oauth")
		oauth_token = nil
		access_token = ""
	} else {
		access_token = oauth_token.AccessToken
	}
	return html.NewResponse("sync/planning-root.html", contentPlanningRoot{
		ArcgisAccessToken: access_token,
	}), nil
}
