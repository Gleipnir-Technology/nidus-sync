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
