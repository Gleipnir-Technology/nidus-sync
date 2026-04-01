package resource

import (
	"context"
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
type avatar struct {
	URI string `json:"uri"`
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
func (res *avatarR) Create(ctx context.Context, r *http.Request, u platform.User, uploads []file.Upload) (*avatar, *nhttp.ErrorWithStatus) {
	if len(uploads) == 0 {
		return nil, nhttp.NewErrorStatus(http.StatusBadRequest, "No upload found")
	}
	if len(uploads) != 1 {
		return nil, nhttp.NewErrorStatus(http.StatusBadRequest, "You must only submit one file at a time")
	}
	upload := uploads[0]
	err := platform.AvatarCreate(r.Context(), u, upload)
	if err != nil {
		return nil, nhttp.NewErrorStatus(http.StatusBadRequest, "Create avatar: %w", err)
	}
	uri, err := res.router.UUIDToURI("avatar.ByIDGet", &upload.UUID)
	if err != nil {
		return nil, nhttp.NewError("create uri: %w", err)
	}
	return &avatar{
		URI: *uri,
	}, nil
}
