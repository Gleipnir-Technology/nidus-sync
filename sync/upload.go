package sync

import (
	"context"

	"github.com/Gleipnir-Technology/nidus-sync/db/models"
)

type contentUploadPlaceholder struct{}

func getUploadList(ctx context.Context, user *models.User) (string, interface{}, *errorWithStatus) {
	content := contentUploadPlaceholder{}
	return "sync/upload-list.html", content, nil
}
