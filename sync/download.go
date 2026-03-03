package sync

import (
	"context"
	"net/http"

	"github.com/Gleipnir-Technology/nidus-sync/db/models"
	"github.com/Gleipnir-Technology/nidus-sync/html"
	nhttp "github.com/Gleipnir-Technology/nidus-sync/http"
)

type contentDownloadPlaceholder struct{}

func getDownloadList(ctx context.Context, r *http.Request, org *models.Organization, user *models.User) (*html.Response[contentDownloadPlaceholder], *nhttp.ErrorWithStatus) {
	content := contentDownloadPlaceholder{}
	return html.NewResponse("sync/download-list.html", content), nil
}
