package sync

import (
	"context"
	"net/http"

	"github.com/Gleipnir-Technology/nidus-sync/db/models"
)

type contentParcel struct{}

func getParcel(ctx context.Context, r *http.Request, org *models.Organization, user *models.User) (*response[contentParcel], *errorWithStatus) {
	return newResponse("sync/parcel.html", contentParcel{}), nil
}
