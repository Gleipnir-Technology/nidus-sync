package api

import (
	"context"
	"net/http"
	"strconv"

	"github.com/Gleipnir-Technology/nidus-sync/config"
	"github.com/Gleipnir-Technology/nidus-sync/html"
	nhttp "github.com/Gleipnir-Technology/nidus-sync/http"
	"github.com/Gleipnir-Technology/nidus-sync/platform"
	"github.com/gorilla/mux"
	"github.com/rs/zerolog/log"
)

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

func getUserSelf(ctx context.Context, r *http.Request, user platform.User, query queryParams) (*contentUserSelf, *nhttp.ErrorWithStatus) {
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

type responseListUser struct {
	Users []*platform.User `json:"users"`
}

func listUser(ctx context.Context, r *http.Request, user platform.User, query queryParams) (*responseListUser, *nhttp.ErrorWithStatus) {
	users, err := platform.UsersByOrg(ctx, user.Organization)
	if err != nil {
		return nil, nhttp.NewError("list users: %w", err)
	}
	results := make([]*platform.User, len(users))
	i := 0
	for _, v := range users {
		results[i] = v
		i++
	}
	return &responseListUser{
		Users: results,
	}, nil
}

type responseListUserSuggestion struct {
	Users []*platform.User `json:"users"`
}

func listUserSuggestion(ctx context.Context, r *http.Request, user platform.User, query queryParams) (*responseListUserSuggestion, *nhttp.ErrorWithStatus) {
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

func userPut(ctx context.Context, r *http.Request, user platform.User, updates platform.UserChangeRequest) (string, *nhttp.ErrorWithStatus) {
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
