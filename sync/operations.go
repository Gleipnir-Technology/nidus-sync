package sync

import (
	"context"
	"net/http"

	"github.com/Gleipnir-Technology/nidus-sync/db/models"
	"github.com/Gleipnir-Technology/nidus-sync/html"
	nhttp "github.com/Gleipnir-Technology/nidus-sync/http"
)

type contentOperationsRoot struct{}

func getOperationsRoot(ctx context.Context, r *http.Request, org *models.Organization, user *models.User) (*html.Response[contentOperationsRoot], *nhttp.ErrorWithStatus) {
	return html.NewResponse("sync/operations-root.html", contentOperationsRoot{}), nil
}
