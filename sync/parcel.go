package sync

import (
	"context"
	"net/http"

	"github.com/Gleipnir-Technology/nidus-sync/html"
	nhttp "github.com/Gleipnir-Technology/nidus-sync/http"
	"github.com/Gleipnir-Technology/nidus-sync/platform"
)

type contentParcel struct{}

func getParcel(ctx context.Context, r *http.Request, user platform.User) (*html.Response[contentParcel], *nhttp.ErrorWithStatus) {
	return html.NewResponse("sync/parcel.html", contentParcel{}), nil
}
