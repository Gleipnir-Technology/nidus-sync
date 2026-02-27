package sync

import (
	"context"
	"net/http"

	"github.com/Gleipnir-Technology/nidus-sync/db/models"
)

type contentPlanningRoot struct{}

func getPlanningRoot(ctx context.Context, r *http.Request, org *models.Organization, user *models.User) (*response[contentPlanningRoot], *errorWithStatus) {
	return newResponse("sync/planning-root.html", contentPlanningRoot{}), nil
}
