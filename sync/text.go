package sync

import (
	"context"
	"net/http"

	"github.com/Gleipnir-Technology/nidus-sync/db/models"
)

type contentTextMessages struct{}

func getTextMessages(ctx context.Context, r *http.Request, org *models.Organization, u *models.User) (*response[contentTextMessages], *errorWithStatus) {
	content := contentTextMessages{}
	return newResponse("sync/text-messages.html", content), nil
}
