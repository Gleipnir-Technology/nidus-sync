package sync

import (
	"context"
	"net/http"

	"github.com/Gleipnir-Technology/nidus-sync/db/models"
)

type contentPoolList struct{}

func getPoolList(ctx context.Context, r *http.Request, org *models.Organization, user *models.User) (*response[contentPoolList], *errorWithStatus) {
	return newResponse("sync/pool-list.html", contentPoolList{}), nil
}
func getPoolCreate(ctx context.Context, r *http.Request, org *models.Organization, user *models.User) (*response[contentPoolList], *errorWithStatus) {
	return newResponse("sync/pool-upload.html", contentPoolList{}), nil
}
func getPoolByID(ctx context.Context, r *http.Request, org *models.Organization, user *models.User) (*response[contentPoolList], *errorWithStatus) {
	return newResponse("sync/pool-by-id.html", contentPoolList{}), nil
}
