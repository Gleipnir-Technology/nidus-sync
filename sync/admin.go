package sync

import (
	"context"

	"github.com/Gleipnir-Technology/nidus-sync/db/models"
)

type contentAdminDash struct{}

func getAdminDash(ctx context.Context, user *models.User) (string, interface{}, *errorWithStatus) {
	content := contentAdminDash{}
	return "sync/admin-dash.html", content, nil
}
