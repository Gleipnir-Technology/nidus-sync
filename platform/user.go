package platform

import (
	"context"
	"fmt"

	"github.com/Gleipnir-Technology/nidus-sync/db"
	"github.com/Gleipnir-Technology/nidus-sync/db/models"
	"github.com/Gleipnir-Technology/nidus-sync/notification"
)

type User struct {
	DisplayName   string `json:"display_name"`
	Initials      string
	Notifications []notification.Notification
	Organization  Organization `json:"organization"`
	Role          string       `json:"role"`
	Username      string       `json:"username"`
}

func UsersByID(ctx context.Context, org *models.Organization) (map[int32]*User, error) {
	users, err := org.User().All(ctx, db.PGInstance.BobDB)
	if err != nil {
		return make(map[int32]*User, 0), fmt.Errorf("get all org users: %w", err)
	}
	organization := NewOrganization(org)
	results := make(map[int32]*User, len(users))
	for _, user := range users {
		results[user.ID] = &User{
			DisplayName:   user.DisplayName,
			Initials:      "",
			Notifications: []notification.Notification{},
			Organization:  organization,
			Role:          user.Role.String(),
			Username:      user.Username,
		}
	}
	return results, nil
}
