package sync

import (
	"context"
	"net/http"

	"github.com/Gleipnir-Technology/nidus-sync/db/models"
)

type contentDownloadPlaceholder struct{}

func getDownloadList(ctx context.Context, r *http.Request, org *models.Organization, user *models.User) (*response[contentDownloadPlaceholder], *errorWithStatus) {
	content := contentDownloadPlaceholder{}
	return newResponse("sync/download-list.html", content), nil
}
