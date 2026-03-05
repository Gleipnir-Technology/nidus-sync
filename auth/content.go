package auth

import (
	"context"
	"strings"

	"github.com/Gleipnir-Technology/nidus-sync/db/models"
	"github.com/Gleipnir-Technology/nidus-sync/notification"
	"github.com/Gleipnir-Technology/nidus-sync/platform"
)

func ContentForUser(ctx context.Context, user *models.User) (platform.User, error) {
	notifications, err := notification.ForUser(ctx, user)
	if err != nil {
		return platform.User{}, err
	}
	org := user.R.Organization
	var organization platform.Organization
	if org != nil {
		organization.ID = int32(org.ID)
		organization.Name = org.Name
	}
	return platform.User{
		DisplayName:   user.DisplayName,
		Initials:      extractInitials(user.DisplayName),
		Notifications: notifications,
		Organization:  organization,
		Role:          user.Role.String(),
		Username:      user.Username,
	}, nil

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
