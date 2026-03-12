package platform

import (
	"context"
	"fmt"
	"strings"

	"github.com/aarondl/opt/omit"
	//"github.com/Gleipnir-Technology/bob"
	"github.com/Gleipnir-Technology/bob/dialect/psql/dialect"
	"github.com/Gleipnir-Technology/bob/mods"
	"github.com/Gleipnir-Technology/nidus-sync/db"
	"github.com/Gleipnir-Technology/nidus-sync/db/enums"
	"github.com/Gleipnir-Technology/nidus-sync/db/models"
	"github.com/Gleipnir-Technology/nidus-sync/debug"
	"github.com/rs/zerolog/log"
)

type NoUserError struct{}

func (e NoUserError) Error() string { return "That user does not exist" }

type User struct {
	DisplayName      string         `json:"display_name"`
	ID               int            `json:"-"`
	Initials         string         `json:"initials"`
	Notifications    []Notification `json:"-"`
	Organization     Organization   `json:"organization"`
	PasswordHash     string         `json:"-"`
	PasswordHashType string         `json:"-"`
	Role             string         `json:"role"`
	Username         string         `json:"username"`

	model *models.User
}

func (u User) HasRoot() bool {
	return u.model.Role != enums.UserroleRoot
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
	return &User{
		DisplayName:   user.DisplayName,
		Initials:      extractInitials(user.DisplayName),
		Notifications: []Notification{},
		Organization:  newOrganization(o),
		Role:          user.Role.String(),
		Username:      user.Username,

		model: user,
	}, nil
}
func UserByID(ctx context.Context, user_id int) (*User, error) {
	return getUser(ctx, models.SelectWhere.Users.ID.EQ(int32(user_id)))
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
		results[user.ID] = &User{
			DisplayName:   user.DisplayName,
			Initials:      "",
			Notifications: []Notification{},
			Organization:  org,
			Role:          user.Role.String(),
			Username:      user.Username,
			model:         user,
		}
	}
	return results, nil
}
func getUser(ctx context.Context, where mods.Where[*dialect.SelectQuery]) (*User, error) {
	user, err := models.Users.Query(
		models.Preload.User.Organization(),
		where,
	).One(ctx, db.PGInstance.BobDB)
	if err != nil {
		if err.Error() == "No such user" || err.Error() == "sql: no rows in result set" {
			return nil, &NoUserError{}
		} else {
			debug.LogErrorTypeInfo(err)
			log.Error().Err(err).Msg("Unrecognized error. This should be updated in the findUser code")
			return nil, err
		}
	}
	org := newOrganization(user.R.Organization)

	return &User{
		DisplayName:   user.DisplayName,
		Initials:      extractInitials(user.DisplayName),
		Notifications: []Notification{},
		Organization:  org,
		Role:          user.Role.String(),
		Username:      user.Username,
	}, nil
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
