package sync

import (
	//"context"
	//"fmt"
	"net/http"
	//"strings"
	//"time"

	"github.com/Gleipnir-Technology/nidus-sync/db/models"
	"github.com/Gleipnir-Technology/nidus-sync/html"
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
	User          User
}

func getNotificationList(w http.ResponseWriter, r *http.Request, u *models.User) {
	userContent, err := contentForUser(r.Context(), u)
	if err != nil {
		respondError(w, "Failed to get user", err, http.StatusInternalServerError)
		return
	}
	ctx := r.Context()
	notifications, err := notification.ForUser(ctx, u)
	if err != nil {
		respondError(w, "Failed to get notifications", err, http.StatusInternalServerError)
	}
	html.RenderOrError(w, "sync/notification-list.html", &contentNotificationList{
		Notifications: notifications,
		User:          userContent,
	})
}
