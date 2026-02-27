package sync

import (
	"context"
	"net/http"

	"github.com/Gleipnir-Technology/nidus-sync/db/models"
)

type contentIntelligenceRoot struct{}

func getIntelligenceRoot(ctx context.Context, r *http.Request, org *models.Organization, user *models.User) (*response[contentIntelligenceRoot], *errorWithStatus) {
	return newResponse("sync/intelligence-root.html", contentIntelligenceRoot{}), nil
}
