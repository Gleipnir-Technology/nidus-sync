package sync

import (
	"context"
	"net/http"

	"github.com/Gleipnir-Technology/nidus-sync/db/models"
	"github.com/Gleipnir-Technology/nidus-sync/html"
	nhttp "github.com/Gleipnir-Technology/nidus-sync/http"
)

type contentMessageList struct{}

func getMessageList(ctx context.Context, r *http.Request, org *models.Organization, user *models.User) (*html.Response[contentMessageList], *nhttp.ErrorWithStatus) {
	content := contentMessageList{}
	return html.NewResponse("sync/message-list.html", content), nil
}
