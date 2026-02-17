package sync

import (
	"context"

	"github.com/Gleipnir-Technology/nidus-sync/db"
	"github.com/Gleipnir-Technology/nidus-sync/db/models"
)

type contentRadar struct {
	Organization *models.Organization
}

func getRadar(ctx context.Context, user *models.User) (string, contentRadar, *errorWithStatus) {
	org, err := user.Organization().One(ctx, db.PGInstance.BobDB)
	if err != nil {
		return "", contentRadar{}, newError("get org: %w", err)
	}
	data := contentRadar{
		Organization: org,
	}
	return "sync/radar.html", data, nil
}
