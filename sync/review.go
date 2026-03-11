package sync

import (
	"context"
	"fmt"
	"html/template"
	"net/http"

	"github.com/Gleipnir-Technology/nidus-sync/background"
	"github.com/Gleipnir-Technology/nidus-sync/config"
	"github.com/Gleipnir-Technology/nidus-sync/db/models"
	"github.com/Gleipnir-Technology/nidus-sync/html"
	nhttp "github.com/Gleipnir-Technology/nidus-sync/http"
	"github.com/rs/zerolog/log"
)

type contentReviewPool struct {
	ArcgisAccessToken string
	URLTiles          template.HTMLAttr
}
type contentReviewRoot struct{}

func getReviewPool(ctx context.Context, r *http.Request, org *models.Organization, user *models.User) (*html.Response[contentReviewPool], *nhttp.ErrorWithStatus) {
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
	return html.NewResponse("sync/review/pool.html", contentReviewPool{
		ArcgisAccessToken: access_token,
		URLTiles:          template.HTMLAttr(fmt.Sprintf(`url-tiles="%s"`, config.MakeURLNidus("/api/tile/{z}/{y}/{x}"))),
	}), nil
}
func getReviewRoot(ctx context.Context, r *http.Request, org *models.Organization, user *models.User) (*html.Response[contentReviewRoot], *nhttp.ErrorWithStatus) {
	return html.NewResponse("sync/review/root.html", contentReviewRoot{}), nil
}
func getReviewSite(ctx context.Context, r *http.Request, org *models.Organization, user *models.User) (*html.Response[contentReviewRoot], *nhttp.ErrorWithStatus) {
	return html.NewResponse("sync/review/site.html", contentReviewRoot{}), nil
}
