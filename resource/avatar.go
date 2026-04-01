package resource

import (
	"context"
	"fmt"
	"net/http"

	nhttp "github.com/Gleipnir-Technology/nidus-sync/http"
	"github.com/Gleipnir-Technology/nidus-sync/platform"
	"github.com/Gleipnir-Technology/nidus-sync/platform/file"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

func Avatar(r *router) *avatarR {
	return &avatarR{
		router: r,
	}
}

type avatarR struct {
	router *router
}

func (res *avatarR) ByIDGet(ctx context.Context, r *http.Request, u platform.User) (file.Collection, uuid.UUID, *nhttp.ErrorWithStatus) {
	vars := mux.Vars(r)
	uid_str := vars["uuid"]
	uid, err := uuid.Parse(uid_str)
	if err != nil {
		return file.CollectionAvatar, uuid.UUID{}, nhttp.NewErrorStatus(http.StatusBadRequest, "parse uuid: %w", err)
	}
	return file.CollectionAvatar, uid, nil
}
func (res *avatarR) Create(ctx context.Context, r *http.Request, u platform.User, uploads []file.Upload) (string, *nhttp.ErrorWithStatus) {
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
