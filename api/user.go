package api

import (
	"context"
	"net/http"

	"github.com/Gleipnir-Technology/nidus-sync/html"
	nhttp "github.com/Gleipnir-Technology/nidus-sync/http"
	"github.com/Gleipnir-Technology/nidus-sync/platform"
)

type contentURLAPI struct {
	Communication string `json:"communication"`
}
type contentURLs struct {
	API    contentURLAPI `json:"api"`
	Tegola string        `json:"tegola"`
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
	user.NotificationCounts = *counts
	urls := html.NewContentURL()
	return &contentUserSelf{
		Self: user,
		URLs: contentURLs{
			API: contentURLAPI{
				Communication: urls.API.Communication,
			},
			Tegola: urls.Tegola,
		},
	}, nil
}

type responseListUser struct {
	Users []platform.User `json:"users"`
}

func listUser(ctx context.Context, r *http.Request, user platform.User, query queryParams) (*responseListUser, *nhttp.ErrorWithStatus) {
	return &responseListUser{
		Users: []platform.User{},
	}, nil
}

type responseListUserSuggestion struct {
	Users []platform.User `json:"users"`
}

func listUserSuggestion(ctx context.Context, r *http.Request, user platform.User, query queryParams) (*responseListUser, *nhttp.ErrorWithStatus) {
	if query.Query == nil {
		return nil, nhttp.NewErrorStatus(http.StatusBadRequest, "you need to include a query")
	}
	users, err := platform.UserSuggestion(ctx, user, *query.Query)
	if err != nil {
		return nil, nhttp.NewError("query suggestions: %w", err)
	}
	return &responseListUser{
		Users: users,
	}, nil
}
