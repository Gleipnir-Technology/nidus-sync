package resource

import (
	"context"
	"fmt"
	"net/http"
	"strconv"

	"github.com/Gleipnir-Technology/nidus-sync/db/models"
	nhttp "github.com/Gleipnir-Technology/nidus-sync/http"
	"github.com/Gleipnir-Technology/nidus-sync/platform"
	"github.com/aarondl/opt/omitnull"
	"github.com/gorilla/mux"
	"github.com/rs/zerolog/log"
)

type user struct {
	Avatar           *string  `json:"avatar"`
	DisplayName      string   `json:"display_name"`
	ID               int      `json:"id"`
	Initials         string   `json:"initials"`
	IsActive         bool     `json:"is_active"`
	PasswordHash     string   `json:"-"`
	PasswordHashType string   `json:"-"`
	Role             string   `json:"role"`
	Tags             []string `json:"tags"`
	URI              string   `json:"uri"`
	Username         string   `json:"username"`
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
	return &user{
		Avatar:      avatar,
		DisplayName: u.DisplayName,
		ID:          int(u.ID),
		Initials:    u.Initials,
		IsActive:    u.Active,
		Role:        u.Role,
		Tags:        u.Tags,
		URI:         uri,
		Username:    u.Username,
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
		return "", nhttp.NewErrorStatus(http.StatusBadRequest, "user update: %w", err)
	}
	user_changes := &models.UserSetter{}
	if updates.Avatar != nil {
		avatar_uuid, err := res.router.UUIDFromURI("avatar.ByUUIDGet", *updates.Avatar)
		if err != nil {
			return "", nhttp.NewBadRequest("parse avatar uri: %w", err)
		}
		user_changes.Avatar = omitnull.FromPtr(avatar_uuid)
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
	log.Debug().Int("len", len(users)).Msg("building response")
	for i, v := range users {
		log.Debug().Int("i", i).Msg("making results")
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
