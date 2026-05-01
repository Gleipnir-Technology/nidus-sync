package resource

import (
	"context"
	"fmt"
	"net/http"
	"strconv"

	"github.com/Gleipnir-Technology/nidus-sync/db/enums"
	nhttp "github.com/Gleipnir-Technology/nidus-sync/http"
	"github.com/Gleipnir-Technology/nidus-sync/platform"
	"github.com/Gleipnir-Technology/nidus-sync/platform/file"
	"github.com/gorilla/mux"
	"github.com/rs/zerolog/log"
)

type uploadR struct {
	router *mux.Router
}

func Upload(r *mux.Router) *uploadR {
	return &uploadR{
		router: r,
	}
}

func (res *uploadR) ByIDGet(ctx context.Context, r *http.Request, u platform.User, query QueryParams) (*platform.Upload, *nhttp.ErrorWithStatus) {
	vars := mux.Vars(r)
	file_id_str := vars["id"]
	file_id_, err := strconv.ParseInt(file_id_str, 10, 32)
	if err != nil {
		return nil, nhttp.NewError("Failed to parse file_id: %w", err)
	}
	file_id := int32(file_id_)
	detail, err := platform.GetUploadDetail(ctx, u.Organization.ID, file_id)
	if err != nil {
		return nil, nhttp.NewError("Failed to get pool: %w", err)
	}
	return detail, nil
}

type contentUploadList struct {
	RecentUploads []platform.Upload
}
type contentUploadPlaceholder struct{}

func (res *uploadR) List(ctx context.Context, r *http.Request, user platform.User, req QueryParams) (*contentUploadPoolList, *nhttp.ErrorWithStatus) {
	rows, err := platform.UploadList(ctx, user.Organization)
	if err != nil {
		return nil, nhttp.NewError("Get upload list: %w", err)
	}
	return &contentUploadPoolList{
		Uploads: rows,
	}, nil
}

type contentUploadDetail struct {
	CSVFileID    int32
	Organization platform.Organization
	Upload       platform.Upload
}
type contentUploadPoolList struct {
	Uploads []platform.Upload `json:"uploads"`
}

type FormUploadCommit struct{}

func (res *uploadR) Commit(ctx context.Context, r *http.Request, u platform.User, f FormUploadCommit) (string, *nhttp.ErrorWithStatus) {
	vars := mux.Vars(r)
	file_id_str := vars["id"]
	file_id_, err := strconv.ParseInt(file_id_str, 10, 32)
	if err != nil {
		return "", nhttp.NewError("Failed to parse file_id: %w", err)
	}
	err = platform.UploadCommit(ctx, u.Organization, int32(file_id_), u)
	if err != nil {
		return "", nhttp.NewError("Failed to mark committed: %w", err)
	}
	log.Debug().Int64("file_id", file_id_).Int("user_id", u.ID).Msg("Committed file")
	return "/configuration/upload", nil
}

type FormUploadDiscard struct{}

func (res *uploadR) Discard(ctx context.Context, r *http.Request, u platform.User, f FormUploadDiscard) (string, *nhttp.ErrorWithStatus) {
	vars := mux.Vars(r)
	file_id_str := vars["id"]
	file_id_, err := strconv.ParseInt(file_id_str, 10, 32)
	if err != nil {
		return "", nhttp.NewError("Failed to parse file_id: %w", err)
	}
	err = platform.UploadDiscard(ctx, u.Organization, int32(file_id_))
	if err != nil {
		return "", nhttp.NewError("Failed to mark discarded: %w", err)
	}
	return "/configuration/upload", nil
}

func (res *uploadR) PoolFlyoverCreate(ctx context.Context, r *http.Request, u platform.User, uploads []file.Upload) (string, *nhttp.ErrorWithStatus) {
	// If the organization we're uploading to doesn't have a service area, we can't process the upload correctly
	if !u.Organization.HasServiceArea() && !u.Organization.IsCatchall() {
		return "", nhttp.NewErrorStatus(http.StatusConflict, "Your organization does not yet have a service area")
	}
	if len(uploads) == 0 {
		return "", nhttp.NewErrorStatus(http.StatusBadRequest, "No upload found")
	}
	if len(uploads) != 1 {
		return "", nhttp.NewErrorStatus(http.StatusBadRequest, "You must only submit one file at a time")
	}
	upload := uploads[0]
	saved_upload, err := platform.NewUpload(r.Context(), u, upload, enums.FileuploadCsvtypeFlyover)
	if err != nil {
		return "", nhttp.NewError("Failed to create new pool: %w", err)
	}
	return fmt.Sprintf("/configuration/upload/%d", *saved_upload), nil
}
func (res *uploadR) PoolCustomCreate(ctx context.Context, r *http.Request, u platform.User, uploads []file.Upload) (string, *nhttp.ErrorWithStatus) {
	if len(uploads) == 0 {
		return "", nhttp.NewErrorStatus(http.StatusBadRequest, "No upload found")
	}
	if len(uploads) != 1 {
		return "", nhttp.NewErrorStatus(http.StatusBadRequest, "You must only submit one file at a time")
	}
	upload := uploads[0]
	pool_upload, err := platform.NewUpload(r.Context(), u, upload, enums.FileuploadCsvtypePoollist)
	if err != nil {
		return "", nhttp.NewError("Failed to create new pool: %w", err)
	}
	return fmt.Sprintf("/configuration/upload/%d", *pool_upload), nil
}
