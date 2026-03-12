package sync

import (
	"context"
	"net/http"

	"github.com/Gleipnir-Technology/nidus-sync/html"
	nhttp "github.com/Gleipnir-Technology/nidus-sync/http"
	"github.com/Gleipnir-Technology/nidus-sync/platform"
)

type contentPoolList struct{}

func getPoolList(ctx context.Context, r *http.Request, user platform.User) (*html.Response[contentPoolList], *nhttp.ErrorWithStatus) {
	return html.NewResponse("sync/pool-list.html", contentPoolList{}), nil
}
func getPoolCreate(ctx context.Context, r *http.Request, user platform.User) (*html.Response[contentPoolList], *nhttp.ErrorWithStatus) {
	return html.NewResponse("sync/pool-upload.html", contentPoolList{}), nil
}
func getPoolByID(ctx context.Context, r *http.Request, user platform.User) (*html.Response[contentPoolList], *nhttp.ErrorWithStatus) {
	return html.NewResponse("sync/pool-by-id.html", contentPoolList{}), nil
}
