package sync

import (
	"context"
	"fmt"
	"net/http"
	"strconv"

	"github.com/Gleipnir-Technology/nidus-sync/db"
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

type contentPoolDetail struct {
	CSVFileID    int32
	Organization *models.Organization
	Upload       platform.UploadPoolDetail
}
type contentUploadPoolList struct {
	Uploads []platform.PoolUpload
}
type contentUploadPoolCreate struct{}

func getUploadPoolList(ctx context.Context, r *http.Request, org *models.Organization, u *models.User) (*response[contentUploadPoolList], *errorWithStatus) {
	uploads, err := platform.PoolUploadList(ctx, u.OrganizationID)
	if err != nil {
		return nil, newError("Failed to get uploads: %w", err)
	}
	data := contentUploadPoolList{
		Uploads: uploads,
	}
	return newResponse("sync/pool-list.html", data), nil
}

func getUploadPoolCreate(ctx context.Context, r *http.Request, org *models.Organization, u *models.User) (*response[contentUploadPoolCreate], *errorWithStatus) {
	data := contentUploadPoolCreate{}
	return newResponse("sync/pool-csv-upload.html", data), nil
}
func getUploadByID(ctx context.Context, r *http.Request, org *models.Organization, u *models.User) (*response[contentPoolDetail], *errorWithStatus) {
	org, err := u.Organization().One(ctx, db.PGInstance.BobDB)
	if err != nil {
		return nil, newError("Failed to get organization: %w", err)
	}
	file_id_str := chi.URLParam(r, "id")
	file_id, err := strconv.ParseInt(file_id_str, 10, 32)
	if err != nil {
		return nil, newError("Failed to parse file_id: %w", err)
	}
	detail, err := platform.GetUploadPoolDetail(ctx, u.OrganizationID, int32(file_id))
	if err != nil {
		return nil, newError("Failed to get pool: %w", err)
	}
	data := contentPoolDetail{
		CSVFileID:    int32(file_id),
		Organization: org,
		Upload:       detail,
	}
	return newResponse("sync/pool-by-id.html", data), nil
}

type FormUploadPool struct{}

func postUploadPoolCreate(ctx context.Context, r *http.Request, u *models.User, f FormUploadPool) (string, *errorWithStatus) {
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
	return fmt.Sprintf("/pool/upload/%d", pool_upload.ID), nil
}
