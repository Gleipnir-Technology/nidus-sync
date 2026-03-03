package sync

import (
	"context"
	"net/http"

	"github.com/Gleipnir-Technology/nidus-sync/db/models"
	"github.com/Gleipnir-Technology/nidus-sync/html"
	nhttp "github.com/Gleipnir-Technology/nidus-sync/http"
)

type contentParcel struct{}

func getParcel(ctx context.Context, r *http.Request, org *models.Organization, user *models.User) (*html.Response[contentParcel], *nhttp.ErrorWithStatus) {
	return html.NewResponse("sync/parcel.html", contentParcel{}), nil
}
