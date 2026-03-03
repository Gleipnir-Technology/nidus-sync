package sync

import (
	"context"
	//"fmt"
	"net/http"
	//"strings"
	//"time"

	"github.com/Gleipnir-Technology/nidus-sync/db/models"
	"github.com/Gleipnir-Technology/nidus-sync/html"
	nhttp "github.com/Gleipnir-Technology/nidus-sync/http"
	"github.com/Gleipnir-Technology/nidus-sync/notification"
	//"github.com/Gleipnir-Technology/bob"
	//"github.com/Gleipnir-Technology/bob/dialect/psql"
	//"github.com/Gleipnir-Technology/bob/dialect/psql/sm"
	//"github.com/Gleipnir-Technology/nidus-sync/db"
	//"github.com/Gleipnir-Technology/nidus-sync/db/sql"
	//"github.com/google/uuid"
	//"github.com/uber/h3-go/v4"
)

type contentNotificationList struct {
	Notifications []notification.Notification
}

func getNotificationList(ctx context.Context, r *http.Request, org *models.Organization, u *models.User) (*html.Response[contentNotificationList], *nhttp.ErrorWithStatus) {
	notifications, err := notification.ForUser(ctx, u)
	if err != nil {
		return nil, nhttp.NewError("Failed to get notifications: %w", err)
	}
	return html.NewResponse("sync/notification-list.html", contentNotificationList{
		Notifications: notifications,
	}), nil
}
