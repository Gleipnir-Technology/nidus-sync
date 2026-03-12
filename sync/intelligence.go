package sync

import (
	"context"
	"net/http"

	"github.com/Gleipnir-Technology/nidus-sync/html"
	nhttp "github.com/Gleipnir-Technology/nidus-sync/http"
	"github.com/Gleipnir-Technology/nidus-sync/platform"
)

type contentIntelligenceRoot struct{}

func getIntelligenceRoot(ctx context.Context, r *http.Request, user platform.User) (*html.Response[contentIntelligenceRoot], *nhttp.ErrorWithStatus) {
	return html.NewResponse("sync/intelligence-root.html", contentIntelligenceRoot{}), nil
}
