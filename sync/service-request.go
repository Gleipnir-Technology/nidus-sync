package sync

import (
	"context"

	"github.com/Gleipnir-Technology/nidus-sync/db/models"
)

type contentServiceRequestPlaceholder struct{}

func getServiceRequestList(ctx context.Context, user *models.User) (string, interface{}, *errorWithStatus) {
	content := contentServiceRequestPlaceholder{}
	return "sync/service-request-list.html", content, nil
}
