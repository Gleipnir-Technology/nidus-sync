package sync

import (
	"context"

	"github.com/Gleipnir-Technology/nidus-sync/db/models"
)

type contentMessageList struct{}

func getMessageList(ctx context.Context, user *models.User) (string, interface{}, *errorWithStatus) {
	content := contentMessageList{}
	return "sync/message-list.html", content, nil
}
