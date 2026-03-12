package sync

import (
	"context"
	"net/http"

	"github.com/Gleipnir-Technology/nidus-sync/html"
	nhttp "github.com/Gleipnir-Technology/nidus-sync/http"
	"github.com/Gleipnir-Technology/nidus-sync/platform"
)

type contentCommunicationRoot struct{}

func getCommunicationRoot(ctx context.Context, r *http.Request, user platform.User) (*html.Response[contentCommunicationRoot], *nhttp.ErrorWithStatus) {
	return html.NewResponse("sync/communication-root.html", contentCommunicationRoot{}), nil
}
