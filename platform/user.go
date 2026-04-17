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
	//"github.com/aarondl/opt/omitnull"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
)

type NoUserError struct{}

func (e NoUserError) Error() string { return "That user does not exist" }

type User struct {
	Active           bool
	Avatar           *uuid.UUID
	DisplayName      string
	ID               int
	Initials         string
	IsDronePilot     bool
	IsWarrant        bool
	Organization     Organization
	PasswordHash     string
	PasswordHashType string
	Role             string
	Username         string

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
func (u User) IsAccountOwner() bool {
	return u.model.Role == enums.UserroleAccountOwner
}
func newUser(ctx context.Context, org Organization, user *models.User) User {
	avatar := user.Avatar.Ptr()
	u := User{
		Active:           true,
		Avatar:           avatar,
		DisplayName:      user.DisplayName,
		ID:               int(user.ID),
		Initials:         extractInitials(user.DisplayName),
		IsDronePilot:     user.IsDronePilot,
		IsWarrant:        user.IsWarrant,
		Organization:     org,
		PasswordHash:     user.PasswordHash,
		PasswordHashType: string(user.PasswordHashType),
		Role:             user.Role.String(),
		Username:         user.Username,

		model: user,
	}
	return u
}

func CreateUser(ctx context.Context, username string, name string, password_hash string) (*User, error) {
	o_setter := models.OrganizationSetter{
		IsCatchall: omit.From(false),
		Name:       omit.From(fmt.Sprintf("%s's organization", username)),
	}
	o, err := models.Organizations.Insert(&o_setter).One(ctx, db.PGInstance.BobDB)
	if err != nil {
		return nil, fmt.Errorf("Failed to create organization: %w", err)
	}
	log.Info().Int32("id", o.ID).Msg("Created organization")
	u_setter := models.UserSetter{
		DisplayName:      omit.From(name),
		IsActive:         omit.From(true),
		IsDronePilot:     omit.From(false),
		IsWarrant:        omit.From(false),
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
func UserList(ctx context.Context, user User) ([]*User, error) {
	var query models.UsersQuery
	var orgByID map[int32]*Organization
	if user.HasRoot() {
		query = models.Users.Query()
		orgs, err := OrganizationList(ctx)
		if err != nil {
			return nil, fmt.Errorf("org list: %w", err)
		}
		orgByID = make(map[int32]*Organization, len(orgs))
		for _, org := range orgs {
			orgByID[org.ID] = org
		}
	} else {
		query = user.Organization.model.User()
		orgByID = make(map[int32]*Organization, 1)
		orgByID[user.model.OrganizationID] = &user.Organization
	}
	rows, err := query.All(ctx, db.PGInstance.BobDB)
	results := make([]*User, len(rows))
	if err != nil {
		return nil, fmt.Errorf("query users: %w", err)
	}
	for i, row := range rows {
		org, ok := orgByID[row.OrganizationID]
		if !ok {
			return nil, fmt.Errorf("get org %d", row.OrganizationID)
		}
		new_user := newUser(ctx, *org, row)
		results[i] = &new_user
	}
	return results, nil
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
func UserSuggestion(ctx context.Context, user User, query string) ([]*User, error) {
	query_arg := "%" + query + "%"
	if user.HasRoot() {
		return userSuggestionRoot(ctx, user, query_arg)
	} else {
		return userSuggestionNonRoot(ctx, user, query_arg)
	}
}
func userSuggestionNonRoot(ctx context.Context, user User, query_arg string) ([]*User, error) {
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
	results := make([]*User, len(users))
	for i, user := range users {
		u := toUser(user)
		results[i] = &u
	}
	return results, nil
}

func UserUpdate(ctx context.Context, user User, user_id int, updates *models.UserSetter) error {
	target_user, err := models.FindUser(ctx, db.PGInstance.BobDB, int32(user_id))
	if err != nil {
		return fmt.Errorf("find user: %w", err)
	}
	if user.model.Role != enums.UserroleRoot && target_user.OrganizationID != target_user.OrganizationID {
		return fmt.Errorf("Current user (%d) isn't allowed to change this user (%d)", user.ID, target_user.ID)
	}
	err = target_user.Update(ctx, db.PGInstance.BobDB, updates)
	return err
}
func userSuggestionRoot(ctx context.Context, user User, query_arg string) ([]*User, error) {
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
	results := make([]*User, len(users))
	for i, user := range users {
		u := toUser(user)
		org := org_map[user.OrganizationID]
		u.Organization = Organization{
			model: org,
		}
		results[i] = &u
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
		} else if err.Error() == "context canceled" {
			return nil, err
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
		DisplayName:      user.DisplayName,
		ID:               int(user.ID),
		Initials:         extractInitials(user.DisplayName),
		Organization:     Organization{},
		PasswordHash:     user.PasswordHash,
		PasswordHashType: string(user.PasswordHashType),
		Role:             user.Role.String(),
		Username:         user.Username,

		model: user,
	}
}
