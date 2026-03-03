package sync

import (
	"context"
	"net/http"

	"github.com/Gleipnir-Technology/nidus-sync/db/models"
	"github.com/Gleipnir-Technology/nidus-sync/html"
	nhttp "github.com/Gleipnir-Technology/nidus-sync/http"
)

type contentPoolList struct{}

func getPoolList(ctx context.Context, r *http.Request, org *models.Organization, user *models.User) (*html.Response[contentPoolList], *nhttp.ErrorWithStatus) {
	return html.NewResponse("sync/pool-list.html", contentPoolList{}), nil
}
func getPoolCreate(ctx context.Context, r *http.Request, org *models.Organization, user *models.User) (*html.Response[contentPoolList], *nhttp.ErrorWithStatus) {
	return html.NewResponse("sync/pool-upload.html", contentPoolList{}), nil
}
func getPoolByID(ctx context.Context, r *http.Request, org *models.Organization, user *models.User) (*html.Response[contentPoolList], *nhttp.ErrorWithStatus) {
	return html.NewResponse("sync/pool-by-id.html", contentPoolList{}), nil
}
