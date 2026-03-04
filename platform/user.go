package platform

import (
	"github.com/Gleipnir-Technology/nidus-sync/notification"
)

type User struct {
	DisplayName   string
	Initials      string
	Notifications []notification.Notification
	Organization  Organization
	Role          string
	Username      string
}
