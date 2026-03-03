package sync

import (
	"context"
	"net/http"

	"github.com/Gleipnir-Technology/nidus-sync/db/models"
	"github.com/Gleipnir-Technology/nidus-sync/html"
	nhttp "github.com/Gleipnir-Technology/nidus-sync/http"
)

type contentTextMessages struct{}

func getTextMessages(ctx context.Context, r *http.Request, org *models.Organization, u *models.User) (*html.Response[contentTextMessages], *nhttp.ErrorWithStatus) {
	content := contentTextMessages{}
	return html.NewResponse("sync/text-messages.html", content), nil
}
