package resource

import (
	"context"
	"fmt"
	"net/http"
	"strconv"

	"github.com/Gleipnir-Technology/nidus-sync/db/enums"
	"github.com/Gleipnir-Technology/nidus-sync/db/models"
	nhttp "github.com/Gleipnir-Technology/nidus-sync/http"
	"github.com/Gleipnir-Technology/nidus-sync/platform"
	"github.com/aarondl/opt/omit"
	"github.com/aarondl/opt/omitnull"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	//"github.com/rs/zerolog/log"
)

type user struct {
	Avatar           omitnull.Val[string] `json:"avatar"`
	DisplayName      omit.Val[string]     `json:"display_name"`
	ID               omit.Val[int]        `json:"id"`
	Initials         omit.Val[string]     `json:"initials"`
	IsActive         omit.Val[bool]       `json:"is_active"`
	PasswordHash     omit.Val[string]     `json:"-"`
	PasswordHashType omit.Val[string]     `json:"-"`
	Role             omit.Val[string]     `json:"role"`
	Tags             omit.Val[[]string]   `json:"tags"`
	URI              omit.Val[string]     `json:"uri"`
	Username         omit.Val[string]     `json:"username"`
}

func User(r *router) *userR {
	return &userR{
		router: r,
	}
}
func (res *userR) response(u *platform.User) (*user, error) {
	if u == nil {
		return nil, fmt.Errorf("nil user")
	}
	avatar, err := res.router.UUIDToURI("avatar.ByUUIDGet", u.Avatar)
	if err != nil {
		return nil, fmt.Errorf("id to uri: %w", err)
	}
	uri, err := res.router.IDToURI("user.ByIDGet", u.ID)
	if err != nil {
		return nil, fmt.Errorf("id to uri: %w", err)
	}
	tags := make([]string, 0)
	if u.IsDronePilot {
		tags = append(tags, "drone pilot")
	}
	if u.IsWarrant {
		tags = append(tags, "warrant")
	}
	return &user{
		Avatar:      omitnull.FromPtr(avatar),
		DisplayName: omit.From(u.DisplayName),
		ID:          omit.From(int(u.ID)),
		Initials:    omit.From(u.Initials),
		IsActive:    omit.From(u.Active),
		Role:        omit.From(u.Role),
		Tags:        omit.From(tags),
		URI:         omit.From(uri),
		Username:    omit.From(u.Username),
	}, nil
}

type userR struct {
	router *router
}
type responseListUser struct {
	Users []*platform.User `json:"users"`
}

func (res *userR) ByIDGet(ctx context.Context, r *http.Request, user platform.User, query QueryParams) (*platform.User, *nhttp.ErrorWithStatus) {
	vars := mux.Vars(r)
	user_id_str := vars["id"]
	user_id, err := strconv.Atoi(user_id_str)
	u, err := platform.UserByID(ctx, int32(user_id))
	if err != nil {
		return nil, nhttp.NewError("get user: %w", err)
	}
	return u, nil
}

func (res *userR) ByIDPut(ctx context.Context, r *http.Request, user platform.User, updates user) (string, *nhttp.ErrorWithStatus) {
	vars := mux.Vars(r)
	user_id_str := vars["id"]
	user_id, err := strconv.Atoi(user_id_str)
	if err != nil {
		return "", nhttp.NewErrorStatus(http.StatusBadRequest, "user id conversion: %w", err)
	}
	user_changes := &models.UserSetter{}
	if !user.HasRoot() && !user.IsAccountOwner() && user.ID != user_id {
		return "", nhttp.NewForbidden("Only account owners can change other users")
	}
	if updates.Avatar.IsValue() {
		avatar_uuid, err := res.router.UUIDFromURI("avatar.ByUUIDGet", updates.Avatar.MustGet())
		if err != nil {
			return "", nhttp.NewBadRequest("parse avatar uri: %w", err)
		}
		user_changes.Avatar = omitnull.FromPtr(avatar_uuid)
	} else if updates.Avatar.IsNull() {
		user_changes.Avatar = omitnull.FromPtr[uuid.UUID](nil)
	}
	if updates.DisplayName.IsValue() {
		user_changes.DisplayName = updates.DisplayName
	}
	if updates.Role.IsValue() {
		// Don't allow privilege escalation
		if user.HasRoot() || user.IsAccountOwner() {
			var role enums.Userrole
			v := updates.Role.MustGet()
			err := role.Scan(v)
			if err != nil {
				return "", nhttp.NewBadRequest("invalid role %s: %w", v, err)
			}
			user_changes.Role = omit.From(role)
		} else {
			return "", nhttp.NewBadRequest("you aren't allowed to change roles")
		}
	}
	if updates.Tags.IsValue() {
		for i, v := range updates.Tags.MustGet() {
			user_changes.IsDronePilot = omit.From(false)
			user_changes.IsWarrant = omit.From(false)
			switch v {
			case "drone pilot":
				user_changes.IsDronePilot = omit.From(true)
			case "warrant":
				user_changes.IsWarrant = omit.From(true)
			default:
				return "", nhttp.NewBadRequest("'%s' (item %d) is not a valid tag", v, i)
			}
		}
	}

	err = platform.UserUpdate(ctx, user, user_id, user_changes)
	if err != nil {
		return "", nhttp.NewError("user update: %w", err)
	}
	return "", nil
}

func (res *userR) SelfGet(ctx context.Context, r *http.Request, user platform.User, query QueryParams) (*user, *nhttp.ErrorWithStatus) {
	resp, err := res.response(&user)
	if err != nil {
		return nil, nhttp.NewError("create response: %w", err)
	}
	return resp, nil
}

func (res *userR) List(ctx context.Context, r *http.Request, u platform.User, query QueryParams) ([]*user, *nhttp.ErrorWithStatus) {
	users, err := platform.UserList(ctx, u)
	if err != nil {
		return nil, nhttp.NewError("list users: %w", err)
	}
	results := make([]*user, len(users))
	//log.Debug().Int("len", len(users)).Msg("building response")
	for i, v := range users {
		//log.Debug().Int("i", i).Msg("making results")
		resp, err := res.response(v)
		if err != nil {
			return nil, nhttp.NewError("create response: %w", err)
		}
		results[i] = resp
	}
	return results, nil
}

type responseListUserSuggestion struct {
	Users []*platform.User `json:"users"`
}

func (res *userR) SuggestionGet(ctx context.Context, r *http.Request, user platform.User, query QueryParams) (*responseListUserSuggestion, *nhttp.ErrorWithStatus) {
	if query.Query == nil {
		return nil, nhttp.NewErrorStatus(http.StatusBadRequest, "you need to include a query")
	}
	users, err := platform.UserSuggestion(ctx, user, *query.Query)
	if err != nil {
		return nil, nhttp.NewError("query suggestions: %w", err)
	}
	return &responseListUserSuggestion{
		Users: users,
	}, nil
}
