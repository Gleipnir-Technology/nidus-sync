package htmlpage

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/Gleipnir-Technology/nidus-sync/db"
	enums "github.com/Gleipnir-Technology/nidus-sync/db/enums"
	"github.com/Gleipnir-Technology/nidus-sync/db/models"
	"github.com/Gleipnir-Technology/nidus-sync/debug"
	"github.com/aarondl/opt/omit"
	"github.com/aarondl/opt/omitnull"
	"github.com/rs/zerolog/log"
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
	updater.Exec(ctx, db.PGInstance.BobDB)
	//user.UserNotifications(
	//models.SelectWhere.Notifications.Link.EQ(NotificationPathOauthReset),
	//).UpdateAll()
}

func notifyOauthInvalid(ctx context.Context, user *models.User) {
	msg := "Oauth token invalidated"
	notificationSetter := models.NotificationSetter{
		Created: omit.From(time.Now()),
		Message: omit.From(msg),
		Link:    omit.From(NotificationPathOauthReset),
		Type:    omit.From(enums.NotificationtypeOauthTokenInvalidated),
	}
	err := user.InsertUserNotifications(ctx, db.PGInstance.BobDB, &notificationSetter)
	if err != nil {
		if strings.HasPrefix(err.Error(), "ERROR: duplicate key value violates unique constraint") {
			log.Info().Str("msg", msg).Int("user_id", int(user.ID)).Msg("Refusing to add another notification with the same type")
			return
		}
		debug.LogErrorTypeInfo(err)
		log.Error().Err(err).Msg("Failed to insert new notification. This is a programmer bug.")
		return
	}
}

func notificationsForUser(ctx context.Context, u *models.User) ([]Notification, error) {
	results := make([]Notification, 0)
	notifications, err := u.UserNotifications(
		models.SelectWhere.Notifications.ResolvedAt.IsNull(),
	).All(ctx, db.PGInstance.BobDB)
	if err != nil {
		return results, fmt.Errorf("Failed to get notifications: %w", err)
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
