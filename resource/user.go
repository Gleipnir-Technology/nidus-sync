package resource

import (
	"context"
	"fmt"
	"net/http"
	"strconv"

	"github.com/Gleipnir-Technology/nidus-sync/config"
	"github.com/Gleipnir-Technology/nidus-sync/html"
	nhttp "github.com/Gleipnir-Technology/nidus-sync/http"
	"github.com/Gleipnir-Technology/nidus-sync/platform"
	"github.com/gorilla/mux"
	"github.com/rs/zerolog/log"
)

type userResponse struct {
	Avatar      string `json:"avatar"`
	DisplayName string `json:"display_name"`
	Initials    string `json:"initials"`
	IsActive    bool   `json:"is_active"`
	//Notifications      []Notification         `json:"notifications"`
	//NotificationCounts UserNotificationCounts `json:"notification_counts"`
	//Organization       Organization           `json:"organization"`
	PasswordHash     string   `json:"-"`
	PasswordHashType string   `json:"-"`
	Role             string   `json:"role"`
	Tags             []string `json:"tags"`
	URI              string   `json:"uri"`
	Username         string   `json:"username"`
}

func User(r *mux.Router) *userR {
	return &userR{
		router: r,
	}
}
func (res *userR) response(u *platform.User) (*userResponse, error) {
	if u == nil {
		return nil, fmt.Errorf("nil user")
	}
	log.Info().Int("id", u.ID).Msg("making response from user")
	i := strconv.FormatInt(int64(u.ID), 10)
	handler := res.router.Get("user.ByIDGet")
	if handler == nil {
		return nil, fmt.Errorf("nil handler")
	}
	uri, err := handler.URL("id", i)
	if err != nil {
		return nil, fmt.Errorf("build uri: %w", err)
	}
	return &userResponse{
		Avatar:      u.Avatar,
		DisplayName: u.DisplayName,
		Initials:    u.Initials,
		IsActive:    u.Active,
		Role:        u.Role,
		Tags:        u.Tags,
		URI:         uri.String(),
		Username:    u.Username,
	}, nil
}

type userR struct {
	router *mux.Router
}
type responseListUser struct {
	Users []*platform.User `json:"users"`
}
type contentURLAPI struct {
	Avatar              string `json:"avatar"`
	Communication       string `json:"communication"`
	PublicreportMessage string `json:"publicreport_message"`
	ReviewTask          string `json:"review_task"`
	Signal              string `json:"signal"`
	Upload              string `json:"upload"`
	User                string `json:"user"`
}
type contentURLs struct {
	API    contentURLAPI `json:"api"`
	Tegola string        `json:"tegola"`
	Tile   string        `json:"tile"`
}
type contentUserSelf struct {
	Self platform.User `json:"self"`
	URLs contentURLs   `json:"urls"`
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

func (res *userR) ByIDPut(ctx context.Context, r *http.Request, user platform.User, updates platform.UserChangeRequest) (string, *nhttp.ErrorWithStatus) {
	log.Info().Str("avatar", updates.Avatar).Msg("doing updates")
	vars := mux.Vars(r)
	user_id_str := vars["id"]
	user_id, err := strconv.Atoi(user_id_str)
	if err != nil {
		return "", nhttp.NewErrorStatus(http.StatusBadRequest, "user update: %w", err)
	}
	err = platform.UserUpdate(ctx, user, user_id, updates)
	if err != nil {
		return "", nhttp.NewError("user update: %w", err)
	}
	return "", nil
}

func (res *userR) SelfGet(ctx context.Context, r *http.Request, user platform.User, query QueryParams) (*contentUserSelf, *nhttp.ErrorWithStatus) {
	counts, err := platform.NotificationCountsForUser(ctx, user)
	if err != nil {
		return nil, nhttp.NewError("get notifications: %w", err)
	}
	org, err := platform.OrganizationByID(ctx, int(user.Organization.ID))
	if err != nil {
		return nil, nhttp.NewError("get org: %w", err)
	}
	user.Organization = *org
	user.NotificationCounts = *counts
	urls := html.NewContentURL()
	return &contentUserSelf{
		Self: user,
		URLs: contentURLs{
			API: contentURLAPI{
				Avatar:              config.MakeURLNidus("/api/avatar"),
				Communication:       urls.API.Communication,
				PublicreportMessage: urls.API.Publicreport.Message,
				ReviewTask:          config.MakeURLNidus("/api/review-task"),
				Signal:              config.MakeURLNidus("/api/signal"),
				Upload:              config.MakeURLNidus("/api/upload"),
				User:                config.MakeURLNidus("/api/user"),
			},
			Tegola: urls.Tegola,
			Tile:   config.MakeURLNidus("/api/tile/{z}/{y}/{x}"),
		},
	}, nil
}

func (res *userR) List(ctx context.Context, r *http.Request, user platform.User, query QueryParams) ([]*userResponse, *nhttp.ErrorWithStatus) {
	users, err := platform.UserList(ctx, user)
	if err != nil {
		return nil, nhttp.NewError("list users: %w", err)
	}
	results := make([]*userResponse, len(users))
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
