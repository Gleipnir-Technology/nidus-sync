package sync

import (
	"context"
	"net/http"

	"github.com/Gleipnir-Technology/nidus-sync/db"
	"github.com/Gleipnir-Technology/nidus-sync/db/models"
	"github.com/Gleipnir-Technology/nidus-sync/html"
	nhttp "github.com/Gleipnir-Technology/nidus-sync/http"
)

type contentRadar struct {
	Organization *models.Organization
}

func getRadar(ctx context.Context, r *http.Request, org *models.Organization, user *models.User) (*html.Response[contentRadar], *nhttp.ErrorWithStatus) {
	org, err := user.Organization().One(ctx, db.PGInstance.BobDB)
	if err != nil {
		return nil, nhttp.NewError("get org: %w", err)
	}
	data := contentRadar{
		Organization: org,
	}
	return html.NewResponse("sync/radar.html", data), nil
}
