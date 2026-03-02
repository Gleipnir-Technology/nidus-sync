package sync

import (
	"context"
	"fmt"
	"net/http"
	"strconv"

	"github.com/Gleipnir-Technology/nidus-sync/db/enums"
	"github.com/Gleipnir-Technology/nidus-sync/db/models"
	"github.com/Gleipnir-Technology/nidus-sync/platform"
	"github.com/Gleipnir-Technology/nidus-sync/userfile"
	"github.com/go-chi/chi/v5"
	"github.com/rs/zerolog/log"
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

type contentUploadDetail struct {
	CSVFileID    int32
	Organization *models.Organization
	Upload       platform.UploadPoolDetail
}
type contentUploadPoolList struct {
	Uploads []platform.Upload
}
type contentUploadPool struct{}

func getUploadPool(ctx context.Context, r *http.Request, org *models.Organization, u *models.User) (*response[contentUploadPool], *errorWithStatus) {
	data := contentUploadPool{}
	return newResponse("sync/upload-csv-pool.html", data), nil
}

type contentUploadPoolFlyoverCreate struct{}

func getUploadPoolFlyoverCreate(ctx context.Context, r *http.Request, org *models.Organization, u *models.User) (*response[contentUploadPoolFlyoverCreate], *errorWithStatus) {
	data := contentUploadPoolFlyoverCreate{}
	return newResponse("sync/upload-csv-pool-flyover.html", data), nil
}

type contentUploadPoolCustomCreate struct{}

func getUploadPoolCustomCreate(ctx context.Context, r *http.Request, org *models.Organization, u *models.User) (*response[contentUploadPoolCustomCreate], *errorWithStatus) {
	data := contentUploadPoolCustomCreate{}
	return newResponse("sync/upload-csv-pool-custom.html", data), nil
}
func getUploadByID(ctx context.Context, r *http.Request, org *models.Organization, u *models.User) (*response[contentUploadDetail], *errorWithStatus) {
	test := newContentURLUpload()
	log.Info().Str("output", test.Discard(123)).Send()
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
	}
	return newResponse("sync/upload-by-id.html", data), nil
}

type FormUploadCommit struct{}

func postUploadCommit(ctx context.Context, r *http.Request, org *models.Organization, u *models.User, f FormUploadCommit) (string, *errorWithStatus) {
	file_id_str := chi.URLParam(r, "id")
	file_id_, err := strconv.ParseInt(file_id_str, 10, 32)
	if err != nil {
		return "", newError("Failed to parse file_id: %w", err)
	}
	err = platform.UploadCommit(ctx, org, int32(file_id_))
	if err != nil {
		return "", newError("Failed to mark discarded: %w", err)
	}
	return "/configuration/upload", nil
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
	return "/configuration/upload", nil
}

type FormUploadPool struct{}

func postUploadPoolFlyoverCreate(ctx context.Context, r *http.Request, org *models.Organization, u *models.User, f FormUploadPool) (string, *errorWithStatus) {
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
	saved_upload, err := platform.NewUpload(r.Context(), u, upload, enums.FileuploadCsvtypeFlyover)
	if err != nil {
		return "", newError("Failed to create new pool: %w", err)
	}
	return fmt.Sprintf("/configuration/upload/%d", saved_upload.ID), nil
}
func postUploadPoolCustomCreate(ctx context.Context, r *http.Request, org *models.Organization, u *models.User, f FormUploadPool) (string, *errorWithStatus) {
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
	pool_upload, err := platform.NewUpload(r.Context(), u, upload, enums.FileuploadCsvtypePoollist)
	if err != nil {
		return "", newError("Failed to create new pool: %w", err)
	}
	return fmt.Sprintf("/configuration/upload/%d", pool_upload.ID), nil
}
