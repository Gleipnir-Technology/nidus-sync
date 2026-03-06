package sync

import (
	"context"
	"net/http"

	"github.com/Gleipnir-Technology/nidus-sync/db/models"
	"github.com/Gleipnir-Technology/nidus-sync/html"
	nhttp "github.com/Gleipnir-Technology/nidus-sync/http"
)

type contentReviewPool struct{}
type contentReviewRoot struct{}

func getReviewPool(ctx context.Context, r *http.Request, org *models.Organization, user *models.User) (*html.Response[contentReviewRoot], *nhttp.ErrorWithStatus) {
	return html.NewResponse("sync/review/pool.html", contentReviewRoot{}), nil
}
func getReviewRoot(ctx context.Context, r *http.Request, org *models.Organization, user *models.User) (*html.Response[contentReviewRoot], *nhttp.ErrorWithStatus) {
	return html.NewResponse("sync/review/root.html", contentReviewRoot{}), nil
}
