package sync

import (
	"context"
	"fmt"
	"html/template"
	"net/http"

	"github.com/Gleipnir-Technology/nidus-sync/config"
	"github.com/Gleipnir-Technology/nidus-sync/html"
	nhttp "github.com/Gleipnir-Technology/nidus-sync/http"
	"github.com/Gleipnir-Technology/nidus-sync/platform"
	//"github.com/rs/zerolog/log"
)

type contentReviewPool struct {
	URLTiles template.HTMLAttr
}
type contentReviewRoot struct{}

func getReviewPool(ctx context.Context, r *http.Request, user platform.User) (*html.Response[contentReviewPool], *nhttp.ErrorWithStatus) {
	return html.NewResponse("sync/review/pool.html", contentReviewPool{
		URLTiles: template.HTMLAttr(fmt.Sprintf(`url-tiles="%s"`, config.MakeURLNidus("/api/tile/{z}/{y}/{x}"))),
	}), nil
}
func getReviewRoot(ctx context.Context, r *http.Request, user platform.User) (*html.Response[contentReviewRoot], *nhttp.ErrorWithStatus) {
	return html.NewResponse("sync/review/root.html", contentReviewRoot{}), nil
}
func getReviewSite(ctx context.Context, r *http.Request, user platform.User) (*html.Response[contentReviewRoot], *nhttp.ErrorWithStatus) {
	return html.NewResponse("sync/review/site.html", contentReviewRoot{}), nil
}
