package api

import (
	"context"
	"net/http"

	"github.com/Gleipnir-Technology/nidus-sync/db/models"
	nhttp "github.com/Gleipnir-Technology/nidus-sync/http"
)

type contentListSignal struct{}

func listSignal(ctx context.Context, r *http.Request, org *models.Organization, user *models.User) (*contentListSignal, *nhttp.ErrorWithStatus) {
	return nil, nil
}
