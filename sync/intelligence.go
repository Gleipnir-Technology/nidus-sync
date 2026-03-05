package sync

import (
	"context"
	"net/http"

	"github.com/Gleipnir-Technology/nidus-sync/db/models"
	"github.com/Gleipnir-Technology/nidus-sync/html"
	nhttp "github.com/Gleipnir-Technology/nidus-sync/http"
)

type contentIntelligenceRoot struct{}

func getIntelligenceRoot(ctx context.Context, r *http.Request, org *models.Organization, user *models.User) (*html.Response[contentIntelligenceRoot], *nhttp.ErrorWithStatus) {
	return html.NewResponse("sync/intelligence-root.html", contentIntelligenceRoot{}), nil
}
