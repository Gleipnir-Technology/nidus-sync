package sync

import (
	"context"
	"net/http"

	"github.com/Gleipnir-Technology/nidus-sync/db/models"
)

type contentReviewRoot struct{}

func getReviewRoot(ctx context.Context, r *http.Request, org *models.Organization, user *models.User) (*response[contentReviewRoot], *errorWithStatus) {
	return newResponse("sync/review-root.html", contentReviewRoot{}), nil
}
