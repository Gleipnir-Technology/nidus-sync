package sync

import (
	"context"
	"net/http"

	"github.com/Gleipnir-Technology/nidus-sync/html"
	nhttp "github.com/Gleipnir-Technology/nidus-sync/http"
	"github.com/Gleipnir-Technology/nidus-sync/platform"
)

type contentMessageList struct{}

func getMessageList(ctx context.Context, r *http.Request, user platform.User) (*html.Response[contentMessageList], *nhttp.ErrorWithStatus) {
	content := contentMessageList{}
	return html.NewResponse("sync/message-list.html", content), nil
}
