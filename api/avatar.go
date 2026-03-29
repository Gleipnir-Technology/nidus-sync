package api

import (
	"context"
	"fmt"
	"net/http"

	nhttp "github.com/Gleipnir-Technology/nidus-sync/http"
	"github.com/Gleipnir-Technology/nidus-sync/platform"
	"github.com/Gleipnir-Technology/nidus-sync/platform/file"
)

func avatarPost(ctx context.Context, r *http.Request, u platform.User, uploads []file.Upload) (string, *nhttp.ErrorWithStatus) {
	if len(uploads) == 0 {
		return "", nhttp.NewErrorStatus(http.StatusBadRequest, "No upload found")
	}
	if len(uploads) != 1 {
		return "", nhttp.NewErrorStatus(http.StatusBadRequest, "You must only submit one file at a time")
	}
	upload := uploads[0]
	err := platform.AvatarCreate(r.Context(), u, upload)
	if err != nil {
		return "", nhttp.NewErrorStatus(http.StatusBadRequest, "Create avatar: %w", err)
	}
	return fmt.Sprintf("/avatar/%s", upload.UUID.String()), nil
}
