package sync

import (
	"context"
	"net/http"

	"github.com/Gleipnir-Technology/nidus-sync/html"
	nhttp "github.com/Gleipnir-Technology/nidus-sync/http"
	"github.com/Gleipnir-Technology/nidus-sync/platform"
)

type contentAdminDash struct{}

func getAdminDash(ctx context.Context, r *http.Request, user platform.User) (*html.Response[contentAdminDash], *nhttp.ErrorWithStatus) {
	content := contentAdminDash{}
	return html.NewResponse("sync/admin-dash.html", content), nil
}
