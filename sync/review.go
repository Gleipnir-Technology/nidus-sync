package sync

import (
	"context"
	"net/http"

	"github.com/Gleipnir-Technology/nidus-sync/db/models"
	"github.com/Gleipnir-Technology/nidus-sync/html"
	nhttp "github.com/Gleipnir-Technology/nidus-sync/http"
)

type contentReviewRoot struct{}

func getReviewRoot(ctx context.Context, r *http.Request, org *models.Organization, user *models.User) (*html.Response[contentReviewRoot], *nhttp.ErrorWithStatus) {
	return html.NewResponse("sync/review-root.html", contentReviewRoot{}), nil
}
