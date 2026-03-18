package platform

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/Gleipnir-Technology/bob/dialect/psql"
	"github.com/Gleipnir-Technology/bob/dialect/psql/im"
	"github.com/Gleipnir-Technology/nidus-sync/db"
	enums "github.com/Gleipnir-Technology/nidus-sync/db/enums"
	"github.com/Gleipnir-Technology/nidus-sync/db/models"
	"github.com/Gleipnir-Technology/nidus-sync/debug"
	//"github.com/aarondl/opt/omit"
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
type UserNotificationCounts struct {
	Communications uint `json:"communication"`
	Home           uint `json:"home"`
}

// Clear all notifications for a given user with the given path
func ClearOauth(ctx context.Context, user *models.User) {
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

func NotifyOauthInvalid(ctx context.Context, user *models.User) {
	msg := "Oauth token invalidated"
	_, err := psql.Insert(
		im.Into("notification", "created", "id", "link", "message", "resolved_at", "type", "user_id"),
		im.Values(
			psql.Arg(time.Now()),
			psql.Raw("DEFAULT"),
			psql.Arg(NotificationPathOauthReset),
			psql.Arg(msg),
			psql.Raw("NULL"),
			psql.Arg(enums.NotificationtypeOauthTokenInvalidated),
			psql.Arg(user.ID),
		),
		//im.OnConflict("user_id", "link").DoNothing(),
		//im.OnConflictOnConstraint("unique_user_link_not_resolved").DoNothing(),
		im.OnConflict("user_id", "link").Where("resolved_at IS NULL").DoNothing(),
	).Exec(ctx, db.PGInstance.BobDB)
	/*
		notificationSetter := models.NotificationSetter{
			Created: omit.From(time.Now()),
			Message: omit.From(msg),
			Link:    omit.From(NotificationPathOauthReset),
			Type:    omit.From(enums.NotificationtypeOauthTokenInvalidated),
		}
		err := user.InsertUserNotifications(ctx, db.PGInstance.BobDB, &notificationSetter)
	*/
	if err != nil {
		if strings.Contains(err.Error(), "ERROR: duplicate key value violates unique constraint") {
			log.Info().Str("msg", msg).Int("user_id", int(user.ID)).Msg("Refusing to add another notification with the same type")
			return
		}
		debug.LogErrorTypeInfo(err)
		log.Error().Err(err).Msg("Failed to insert new notification. This is a programmer bug.")
		return
	}
}

func NotificationsForUser(ctx context.Context, u User) ([]Notification, error) {
	results := make([]Notification, 0)
	notifications, err := u.model.UserNotifications(
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
func NotificationCountsForUser(ctx context.Context, u User) (*UserNotificationCounts, error) {
	count_home, err := u.model.UserNotifications(
		models.SelectWhere.Notifications.ResolvedAt.IsNull(),
	).Count(ctx, db.PGInstance.BobDB)
	if err != nil {
		return nil, fmt.Errorf("Failed to get home notification count: %w", err)
	}
	count_reports, err := u.Organization.model.Reports(
		models.SelectWhere.PublicreportReports.Reviewed.IsNull(),
	).Count(ctx, db.PGInstance.BobDB)
	if err != nil {
		return nil, fmt.Errorf("Failed to get nuisance notification count: %w", err)
	}
	log.Debug().Int64("reports", count_reports).Int64("home", count_home).Int("user", u.ID).Msg("calculated notification counts")
	return &UserNotificationCounts{
		Communications: uint(count_reports),
		Home:           uint(count_home),
	}, nil
}

func notificationTypeName(t enums.Notificationtype) string {
	switch t {
	case enums.NotificationtypeOauthTokenInvalidated:
		return "oauth-token-invalid"
	default:
		return "unknown-type"
	}
}
