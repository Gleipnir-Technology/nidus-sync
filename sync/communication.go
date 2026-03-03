package sync

import (
	"context"
	"net/http"

	"github.com/Gleipnir-Technology/nidus-sync/db/models"
	"github.com/Gleipnir-Technology/nidus-sync/html"
	nhttp "github.com/Gleipnir-Technology/nidus-sync/http"
)

type contentCommunicationRoot struct{}

func getCommunicationRoot(ctx context.Context, r *http.Request, org *models.Organization, user *models.User) (*html.Response[contentCommunicationRoot], *nhttp.ErrorWithStatus) {
	return html.NewResponse("sync/communication-root.html", contentCommunicationRoot{}), nil
}
