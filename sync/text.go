package sync

import (
	"context"
	"net/http"

	"github.com/Gleipnir-Technology/nidus-sync/html"
	nhttp "github.com/Gleipnir-Technology/nidus-sync/http"
	"github.com/Gleipnir-Technology/nidus-sync/platform"
)

type contentTextMessages struct{}

func getTextMessages(ctx context.Context, r *http.Request, u platform.User) (*html.Response[contentTextMessages], *nhttp.ErrorWithStatus) {
	content := contentTextMessages{}
	return html.NewResponse("sync/text-messages.html", content), nil
}
