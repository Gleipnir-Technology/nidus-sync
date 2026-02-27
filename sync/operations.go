package sync

import (
	"context"
	"net/http"

	"github.com/Gleipnir-Technology/nidus-sync/db/models"
)

type contentOperationsRoot struct{}

func getOperationsRoot(ctx context.Context, r *http.Request, org *models.Organization, user *models.User) (*response[contentOperationsRoot], *errorWithStatus) {
	return newResponse("sync/operations-root.html", contentOperationsRoot{}), nil
}
