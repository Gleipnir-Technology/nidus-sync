package main

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	enums "github.com/Gleipnir-Technology/nidus-sync/enums"
	"github.com/Gleipnir-Technology/nidus-sync/models"
	"github.com/aarondl/opt/omit"
	"github.com/aarondl/opt/omitnull"
)

var (
	NotificationPathOauthReset string = "/oauth/refresh"
)

type Notification struct {
	Link    string
	Message string
	Time    time.Time
	Type    string
}

// Clear all notifications for a given user with the given path
func clearNotificationsOauth(ctx context.Context, user *models.User) {
	setter := models.NotificationSetter{
		ResolvedAt: omitnull.From(time.Now()),
	}
	updater := models.Notifications.Update(
		//models.SelectWhere.Notifications.Link.EQ(NotificationPathOauthReset),
		models.UpdateWhere.Notifications.Link.EQ(NotificationPathOauthReset),
		models.UpdateWhere.Notifications.UserID.EQ(user.ID),
		setter.UpdateMod(),
	)
	updater.Exec(ctx, PGInstance.BobDB)
	//user.UserNotifications(
	//models.SelectWhere.Notifications.Link.EQ(NotificationPathOauthReset),
	//).UpdateAll()
}

func notifyOauthInvalid(ctx context.Context, user *models.User) {
	notificationSetter := models.NotificationSetter{
		Created: omit.From(time.Now()),
		Message: omit.From("Oauth token invalidated"),
		Link:    omit.From(NotificationPathOauthReset),
		Type:    omit.From(enums.NotificationtypeOauthTokenInvalidated),
	}
	err := user.InsertUserNotifications(ctx, PGInstance.BobDB, &notificationSetter)
	if err != nil {
		LogErrorTypeInfo(err)
		slog.Error("Failed to insert new notification. Update this clause to detect duplicate inserts.", slog.String("err", err.Error()))
		return
	}
}

func notificationsForUser(ctx context.Context, u *models.User) ([]Notification, error) {
	results := make([]Notification, 0)
	notifications, err := u.UserNotifications(
		models.SelectWhere.Notifications.ResolvedAt.IsNull(),
	).All(ctx, PGInstance.BobDB)
	if err != nil {
		return results, fmt.Errorf("Failed to get notifications: %v", err)
	}
	for _, n := range notifications {
		results = append(results, Notification{
			Link:    n.Link,
			Message: n.Message,
			Time:    n.Created,
			Type:    notificationTypeName(n.Type),
		})
	}
	return results, nil
}

func notificationTypeName(t enums.Notificationtype) string {
	switch t {
	case enums.NotificationtypeOauthTokenInvalidated:
		return "alert"
	default:
		return "unknown-type"
	}
}
