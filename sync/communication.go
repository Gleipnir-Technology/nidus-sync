package sync

import (
	"context"
	"net/http"

	"github.com/Gleipnir-Technology/nidus-sync/db/models"
)

type contentCommunicationRoot struct{}

func getCommunicationRoot(ctx context.Context, r *http.Request, org *models.Organization, user *models.User) (*response[contentCommunicationRoot], *errorWithStatus) {
	return newResponse("sync/communication-root.html", contentCommunicationRoot{}), nil
}
