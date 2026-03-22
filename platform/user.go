package platform

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/Gleipnir-Technology/bob/dialect/psql"
	"github.com/Gleipnir-Technology/bob/dialect/psql/dialect"
	"github.com/Gleipnir-Technology/bob/dialect/psql/sm"
	"github.com/Gleipnir-Technology/bob/mods"
	"github.com/Gleipnir-Technology/nidus-sync/db"
	"github.com/Gleipnir-Technology/nidus-sync/db/enums"
	"github.com/Gleipnir-Technology/nidus-sync/db/models"
	"github.com/Gleipnir-Technology/nidus-sync/debug"
	"github.com/aarondl/opt/omit"
	"github.com/rs/zerolog/log"
)

type NoUserError struct{}

func (e NoUserError) Error() string { return "That user does not exist" }

type User struct {
	DisplayName        string                 `json:"display_name"`
	ID                 int                    `json:"-"`
	Initials           string                 `json:"initials"`
	Notifications      []Notification         `json:"notifications"`
	NotificationCounts UserNotificationCounts `json:"notification_counts"`
	Organization       Organization           `json:"organization"`
	PasswordHash       string                 `json:"-"`
	PasswordHashType   string                 `json:"-"`
	Role               string                 `json:"role"`
	Username           string                 `json:"username"`

	model *models.User
}

func (u User) AsJSON() string {
	content, err := json.Marshal(u)
	if err != nil {
		return fmt.Sprintf("{error: \"%s\"}", err.Error())
	}
	return string(content)
}
func (u User) HasRoot() bool {
	return u.model.Role == enums.UserroleRoot
}
func newUser(ctx context.Context, org Organization, user *models.User) User {
	u := User{
		DisplayName:        user.DisplayName,
		ID:                 int(user.ID),
		Initials:           extractInitials(user.DisplayName),
		Notifications:      []Notification{},
		NotificationCounts: UserNotificationCounts{},
		Organization:       org,
		PasswordHash:       user.PasswordHash,
		PasswordHashType:   string(user.PasswordHashType),
		Role:               user.Role.String(),
		Username:           user.Username,

		model: user,
	}
	counts, err := NotificationCountsForUser(ctx, u)
	if err != nil {
		log.Error().Err(err).Int32("id", user.ID).Msg("failed to get notification counts for user")
	}
	u.NotificationCounts = *counts
	return u
}

func CreateUser(ctx context.Context, username string, name string, password_hash string) (*User, error) {
	o_setter := models.OrganizationSetter{
		Name: omit.From(fmt.Sprintf("%s's organization", username)),
	}
	o, err := models.Organizations.Insert(&o_setter).One(ctx, db.PGInstance.BobDB)
	if err != nil {
		return nil, fmt.Errorf("Failed to create organization: %w", err)
	}
	log.Info().Int32("id", o.ID).Msg("Created organization")
	u_setter := models.UserSetter{
		DisplayName:      omit.From(name),
		OrganizationID:   omit.From(o.ID),
		PasswordHash:     omit.From(password_hash),
		PasswordHashType: omit.From(enums.HashtypeBcrypt14),
		Role:             omit.From(enums.UserroleAccountOwner),
		Username:         omit.From(username),
	}
	user, err := models.Users.Insert(&u_setter).One(ctx, db.PGInstance.BobDB)
	if err != nil {
		return nil, fmt.Errorf("Failed to create user: %w", err)
	}
	log.Info().Int32("id", user.ID).Str("username", user.Username).Msg("Created user")
	u := newUser(ctx, newOrganization(o), user)
	return &u, nil
}
func UserByID(ctx context.Context, user_id int32) (*User, error) {
	return getUser(ctx, models.SelectWhere.Users.ID.EQ(user_id))
}
func UserByUsername(ctx context.Context, username string) (*User, error) {
	return getUser(ctx, models.SelectWhere.Users.Username.EQ(username))
}
func UsersByOrg(ctx context.Context, org Organization) (map[int32]*User, error) {
	users, err := org.model.User().All(ctx, db.PGInstance.BobDB)
	if err != nil {
		return make(map[int32]*User, 0), fmt.Errorf("get all org users: %w", err)
	}
	results := make(map[int32]*User, len(users))
	for _, user := range users {
		u := newUser(ctx, org, user)
		results[user.ID] = &u
	}
	return results, nil
}
func UserSuggestion(ctx context.Context, user User, query string) ([]User, error) {
	query_arg := "%" + query + "%"
	if user.HasRoot() {
		return userSuggestionRoot(ctx, user, query_arg)
	} else {
		return userSuggestionNonRoot(ctx, user, query_arg)
	}
}
func userSuggestionNonRoot(ctx context.Context, user User, query_arg string) ([]User, error) {
	users, err := models.Users.Query(
		sm.Where(
			psql.Or(
				psql.Quote("username").ILike(psql.Arg(query_arg)),
				psql.Quote("display_name").ILike(psql.Arg(query_arg)),
			),
		),
		sm.Where(
			psql.Quote("organization_id").EQ(psql.Arg(user.Organization.ID)),
		),
	).All(ctx, db.PGInstance.BobDB)
	if err != nil {
		return nil, fmt.Errorf("query users: %w", err)
	}
	results := make([]User, len(users))
	for i, user := range users {
		results[i] = toUser(user)
	}
	return results, nil
}
func userSuggestionRoot(ctx context.Context, user User, query_arg string) ([]User, error) {
	users, err := models.Users.Query(
		sm.Where(
			psql.Or(
				psql.Quote("username").ILike(psql.Arg(query_arg)),
				psql.Quote("display_name").ILike(psql.Arg(query_arg)),
			),
		),
	).All(ctx, db.PGInstance.BobDB)
	if err != nil {
		return nil, fmt.Errorf("query users: %w", err)
	}
	organization_ids := make([]int32, 0)
	for _, user := range users {
		organization_ids = append(organization_ids, user.OrganizationID)
	}
	orgs, err := models.Organizations.Query(
		sm.Where(
			psql.Quote("id").EQ(psql.Any(organization_ids)),
		),
	).All(ctx, db.PGInstance.BobDB)
	if err != nil {
		return nil, fmt.Errorf("query orgs: %w", err)
	}
	org_map := make(map[int32]*models.Organization, len(orgs))
	for _, org := range orgs {
		org_map[org.ID] = org
	}
	results := make([]User, len(users))
	for i, user := range users {
		u := toUser(user)
		org := org_map[user.OrganizationID]
		u.Organization = Organization{
			model: org,
		}
		results[i] = u
	}
	return results, nil
}
func getUser(ctx context.Context, where mods.Where[*dialect.SelectQuery]) (*User, error) {
	user, err := models.Users.Query(
		models.Preload.User.Organization(),
		where,
	).One(ctx, db.PGInstance.BobDB)
	if err != nil {
		log.Debug().Err(err).Msg("getUser failed")
		if err.Error() == "No such user" || err.Error() == "sql: no rows in result set" {
			return nil, &NoUserError{}
		} else {
			debug.LogErrorTypeInfo(err)
			log.Error().Err(err).Msg("Unrecognized error. This should be updated in the findUser code")
			return nil, err
		}
	}
	org := newOrganization(user.R.Organization)

	u := newUser(ctx, org, user)
	return &u, nil
}
func extractInitials(name string) string {
	parts := strings.Fields(name)
	var initials strings.Builder

	for _, part := range parts {
		if len(part) > 0 {
			initials.WriteString(strings.ToUpper(string(part[0])))
		}
	}

	return initials.String()
}
func toUser(user *models.User) User {
	return User{
		DisplayName:        user.DisplayName,
		ID:                 int(user.ID),
		Initials:           extractInitials(user.DisplayName),
		Notifications:      []Notification{},
		NotificationCounts: UserNotificationCounts{},
		Organization:       Organization{},
		PasswordHash:       user.PasswordHash,
		PasswordHashType:   string(user.PasswordHashType),
		Role:               user.Role.String(),
		Username:           user.Username,

		model: user,
	}
}
