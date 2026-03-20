package api

import (
	"context"
	"net/http"

	nhttp "github.com/Gleipnir-Technology/nidus-sync/http"
	"github.com/Gleipnir-Technology/nidus-sync/platform"
)

func getUser(ctx context.Context, r *http.Request, user platform.User, query queryParams) (*platform.User, *nhttp.ErrorWithStatus) {
	counts, err := platform.NotificationCountsForUser(ctx, user)
	if err != nil {
		return nil, nhttp.NewError("get notifications: %w", err)
	}
	user.NotificationCounts = *counts
	return &user, nil
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
