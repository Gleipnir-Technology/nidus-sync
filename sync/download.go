package sync

import (
	"context"
	"net/http"

	"github.com/Gleipnir-Technology/nidus-sync/html"
	nhttp "github.com/Gleipnir-Technology/nidus-sync/http"
	"github.com/Gleipnir-Technology/nidus-sync/platform"
)

type contentDownloadPlaceholder struct{}

func getDownloadList(ctx context.Context, r *http.Request, user platform.User) (*html.Response[contentDownloadPlaceholder], *nhttp.ErrorWithStatus) {
	content := contentDownloadPlaceholder{}
	return html.NewResponse("sync/download-list.html", content), nil
}
