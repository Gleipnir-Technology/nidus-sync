package sync

import (
	"context"

	"github.com/Gleipnir-Technology/nidus-sync/db/models"
)

type contentDownloadPlaceholder struct{}

func getDownloadList(ctx context.Context, user *models.User) (string, interface{}, *errorWithStatus) {
	content := contentSettingPlaceholder{}
	return "sync/download-list.html", content, nil
}
