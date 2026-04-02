package platform

import (
	"context"
	"fmt"
)

type session struct {
	Impersonating      *User
	NotificationCounts notificationCounts
}

func SessionCurrent(ctx context.Context, user User) (*session, error) {
	counts, err := NotificationCountsForUser(ctx, user)
	if err != nil {
		return nil, fmt.Errorf("get notifications: %w", err)
	}
	return &session{
		Impersonating:      nil,
		NotificationCounts: *counts,
	}, nil
}
