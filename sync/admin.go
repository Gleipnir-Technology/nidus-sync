package sync

import (
	"context"
	"net/http"

	"github.com/Gleipnir-Technology/nidus-sync/db/models"
)

type contentAdminDash struct{}

func getAdminDash(ctx context.Context, r *http.Request, org *models.Organization, user *models.User) (*response[contentAdminDash], *errorWithStatus) {
	content := contentAdminDash{}
	return newResponse("sync/admin-dash.html", content), nil
}
