package sync

import (
	"context"
	"net/http"

	"github.com/Gleipnir-Technology/nidus-sync/db/models"
	"github.com/Gleipnir-Technology/nidus-sync/html"
	nhttp "github.com/Gleipnir-Technology/nidus-sync/http"
)

type contentPlanningRoot struct{}

func getPlanningRoot(ctx context.Context, r *http.Request, org *models.Organization, user *models.User) (*html.Response[contentPlanningRoot], *nhttp.ErrorWithStatus) {
	return html.NewResponse("sync/planning-root.html", contentPlanningRoot{}), nil
}
