package sync

import (
	"context"
	"fmt"
	"net/http"
	"strconv"

	"github.com/Gleipnir-Technology/nidus-sync/config"
	"github.com/Gleipnir-Technology/nidus-sync/db/models"
	"github.com/Gleipnir-Technology/nidus-sync/platform"
	"github.com/Gleipnir-Technology/nidus-sync/userfile"
	"github.com/go-chi/chi/v5"
)

type contentUploadList struct {
	RecentUploads []platform.UploadSummary
}
type contentUploadPlaceholder struct{}

func getUploadList(ctx context.Context, r *http.Request, org *models.Organization, user *models.User) (*response[contentUploadList], *errorWithStatus) {
	rows, err := platform.UploadSummaryList(ctx, org)
	return newResponse("sync/upload-list.html", contentUploadList{
		RecentUploads: rows,
	}), newErrorMaybe("get upload list: %w", err)
}

type contentUploadURL struct {
	Discard string // URL for discarding the upload
}

func newContentUploadURL(id int32) contentUploadURL {
	id_str := strconv.FormatInt(int64(id), 10)
	return contentUploadURL{
		Discard: config.MakeURLNidus("/upload/%s/discard", id_str),
	}
}

type contentUploadDetail struct {
	CSVFileID    int32
	Organization *models.Organization
	Upload       platform.UploadPoolDetail
	URL          contentUploadURL
}
type contentUploadPoolList struct {
	Uploads []platform.PoolUpload
}
type contentUploadPoolCreate struct{}

func getUploadPoolCreate(ctx context.Context, r *http.Request, org *models.Organization, u *models.User) (*response[contentUploadPoolCreate], *errorWithStatus) {
	data := contentUploadPoolCreate{}
	return newResponse("sync/upload-csv-pool.html", data), nil
}
func getUploadByID(ctx context.Context, r *http.Request, org *models.Organization, u *models.User) (*response[contentUploadDetail], *errorWithStatus) {
	file_id_str := chi.URLParam(r, "id")
	file_id_, err := strconv.ParseInt(file_id_str, 10, 32)
	if err != nil {
		return nil, newError("Failed to parse file_id: %w", err)
	}
	file_id := int32(file_id_)
	detail, err := platform.GetUploadPoolDetail(ctx, u.OrganizationID, file_id)
	if err != nil {
		return nil, newError("Failed to get pool: %w", err)
	}
	data := contentUploadDetail{
		CSVFileID:    file_id,
		Organization: org,
		Upload:       detail,
		URL:          newContentUploadURL(file_id),
	}
	return newResponse("sync/upload-by-id.html", data), nil
}

type FormUploadDiscard struct{}

func postUploadDiscard(ctx context.Context, r *http.Request, org *models.Organization, u *models.User, f FormUploadDiscard) (string, *errorWithStatus) {
	file_id_str := chi.URLParam(r, "id")
	file_id_, err := strconv.ParseInt(file_id_str, 10, 32)
	if err != nil {
		return "", newError("Failed to parse file_id: %w", err)
	}
	err = platform.UploadDiscard(ctx, org, int32(file_id_))
	if err != nil {
		return "", newError("Failed to mark discarded: %w", err)
	}
	return "/upload", nil
}

type FormUploadPool struct{}

func postUploadPoolCreate(ctx context.Context, r *http.Request, org *models.Organization, u *models.User, f FormUploadPool) (string, *errorWithStatus) {
	uploads, err := userfile.SaveFileUpload(r, "csvfile", userfile.CollectionCSV)
	if err != nil {
		return "", newError("Failed to extract image uploads: %s", err)
	}
	if len(uploads) == 0 {
		return "", newErrorStatus(http.StatusBadRequest, "No upload found")
	}
	if len(uploads) != 1 {
		return "", newErrorStatus(http.StatusBadRequest, "You must only submit one file at a time")
	}
	upload := uploads[0]
	pool_upload, err := platform.NewPoolUpload(r.Context(), u, upload)
	if err != nil {
		return "", newError("Failed to create new pool: %w", err)
	}
	return fmt.Sprintf("/upload/%d", pool_upload.ID), nil
}
