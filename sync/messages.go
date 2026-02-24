package sync

import (
	"context"
	"net/http"

	"github.com/Gleipnir-Technology/nidus-sync/db/models"
)

type contentMessageList struct{}

func getMessageList(ctx context.Context, r *http.Request, org *models.Organization, user *models.User) (*response[contentMessageList], *errorWithStatus) {
	content := contentMessageList{}
	return newResponse("sync/message-list.html", content), nil
}
