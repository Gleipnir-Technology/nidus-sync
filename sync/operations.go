package sync

import (
	"context"
	"net/http"

	"github.com/Gleipnir-Technology/nidus-sync/html"
	nhttp "github.com/Gleipnir-Technology/nidus-sync/http"
	"github.com/Gleipnir-Technology/nidus-sync/platform"
)

type contentOperationsRoot struct{}

func getOperationsRoot(ctx context.Context, r *http.Request, user platform.User) (*html.Response[contentOperationsRoot], *nhttp.ErrorWithStatus) {
	return html.NewResponse("sync/operations-root.html", contentOperationsRoot{}), nil
}
