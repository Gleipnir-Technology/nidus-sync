package sync

import (
	"context"
	"net/http"

	"github.com/Gleipnir-Technology/nidus-sync/db/enums"
	"github.com/Gleipnir-Technology/nidus-sync/db/models"
)

type contentSudo struct{}

func getSudo(ctx context.Context, user *models.User) (string, interface{}, *errorWithStatus) {
	if user.Role != enums.UserroleRoot {
		return "", nil, &errorWithStatus{
			Message: "You have to be a root user to access this",
			Status:  http.StatusForbidden,
		}
	}
	content := contentAdminDash{}
	return "sync/sudo.html", content, nil
}
