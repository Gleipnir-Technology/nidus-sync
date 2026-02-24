package sync

import (
	"context"
	"net/http"

	"github.com/Gleipnir-Technology/nidus-sync/db"
	"github.com/Gleipnir-Technology/nidus-sync/db/models"
)

type contentRadar struct {
	Organization *models.Organization
}

func getRadar(ctx context.Context, r *http.Request, org *models.Organization, user *models.User) (*response[contentRadar], *errorWithStatus) {
	org, err := user.Organization().One(ctx, db.PGInstance.BobDB)
	if err != nil {
		return nil, newError("get org: %w", err)
	}
	data := contentRadar{
		Organization: org,
	}
	return newResponse("sync/radar.html", data), nil
}
