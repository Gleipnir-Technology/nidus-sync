package sync

import (
	"context"
	"fmt"
	"net/http"
	"strconv"

	"github.com/Gleipnir-Technology/nidus-sync/db/enums"
	"github.com/Gleipnir-Technology/nidus-sync/db/models"
	"github.com/Gleipnir-Technology/nidus-sync/html"
	nhttp "github.com/Gleipnir-Technology/nidus-sync/http"
	"github.com/Gleipnir-Technology/nidus-sync/platform"
	"github.com/Gleipnir-Technology/nidus-sync/userfile"
	"github.com/go-chi/chi/v5"
	//"github.com/rs/zerolog/log"
)

type contentUploadList struct {
	RecentUploads []platform.UploadSummary
}
type contentUploadPlaceholder struct{}

func getUploadList(ctx context.Context, r *http.Request, org *models.Organization, user *models.User) (*html.Response[contentUploadList], *nhttp.ErrorWithStatus) {
	rows, err := platform.UploadSummaryList(ctx, org)
	return html.NewResponse("sync/upload-list.html", contentUploadList{
		RecentUploads: rows,
	}), nhttp.NewErrorMaybe("get upload list: %w", err)
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

func getUploadPool(ctx context.Context, r *http.Request, org *models.Organization, u *models.User) (*html.Response[contentUploadPool], *nhttp.ErrorWithStatus) {
	data := contentUploadPool{}
	return html.NewResponse("sync/upload-csv-pool.html", data), nil
}

type contentUploadPoolFlyoverCreate struct{}

func getUploadPoolFlyoverCreate(ctx context.Context, r *http.Request, org *models.Organization, u *models.User) (*html.Response[contentUploadPoolFlyoverCreate], *nhttp.ErrorWithStatus) {
	data := contentUploadPoolFlyoverCreate{}
	return html.NewResponse("sync/upload-csv-pool-flyover.html", data), nil
}

type contentUploadPoolCustomCreate struct{}

func getUploadPoolCustomCreate(ctx context.Context, r *http.Request, org *models.Organization, u *models.User) (*html.Response[contentUploadPoolCustomCreate], *nhttp.ErrorWithStatus) {
	data := contentUploadPoolCustomCreate{}
	return html.NewResponse("sync/upload-csv-pool-custom.html", data), nil
}
func getUploadByID(ctx context.Context, r *http.Request, org *models.Organization, u *models.User) (*html.Response[contentUploadDetail], *nhttp.ErrorWithStatus) {
	file_id_str := chi.URLParam(r, "id")
	file_id_, err := strconv.ParseInt(file_id_str, 10, 32)
	if err != nil {
		return nil, nhttp.NewError("Failed to parse file_id: %w", err)
	}
	file_id := int32(file_id_)
	detail, err := platform.GetUploadDetail(ctx, u.OrganizationID, file_id)
	if err != nil {
		return nil, nhttp.NewError("Failed to get pool: %w", err)
	}
	data := contentUploadDetail{
		CSVFileID:    file_id,
		Organization: org,
		Upload:       detail,
	}
	return html.NewResponse("sync/upload-by-id.html", data), nil
}

type FormUploadCommit struct{}

func postUploadCommit(ctx context.Context, r *http.Request, org *models.Organization, u *models.User, f FormUploadCommit) (string, *nhttp.ErrorWithStatus) {
	file_id_str := chi.URLParam(r, "id")
	file_id_, err := strconv.ParseInt(file_id_str, 10, 32)
	if err != nil {
		return "", nhttp.NewError("Failed to parse file_id: %w", err)
	}
	err = platform.UploadCommit(ctx, org, int32(file_id_))
	if err != nil {
		return "", nhttp.NewError("Failed to mark committed: %w", err)
	}
	return "/configuration/upload", nil
}

type FormUploadDiscard struct{}

func postUploadDiscard(ctx context.Context, r *http.Request, org *models.Organization, u *models.User, f FormUploadDiscard) (string, *nhttp.ErrorWithStatus) {
	file_id_str := chi.URLParam(r, "id")
	file_id_, err := strconv.ParseInt(file_id_str, 10, 32)
	if err != nil {
		return "", nhttp.NewError("Failed to parse file_id: %w", err)
	}
	err = platform.UploadDiscard(ctx, org, int32(file_id_))
	if err != nil {
		return "", nhttp.NewError("Failed to mark discarded: %w", err)
	}
	return "/configuration/upload", nil
}

type FormUploadPool struct{}

func postUploadPoolFlyoverCreate(ctx context.Context, r *http.Request, org *models.Organization, u *models.User, f FormUploadPool) (string, *nhttp.ErrorWithStatus) {
	uploads, err := userfile.SaveFileUpload(r, "csvfile", userfile.CollectionCSV)
	if err != nil {
		return "", nhttp.NewError("Failed to extract image uploads: %s", err)
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
	return fmt.Sprintf("/configuration/upload/%d", saved_upload.ID), nil
}
func postUploadPoolCustomCreate(ctx context.Context, r *http.Request, org *models.Organization, u *models.User, f FormUploadPool) (string, *nhttp.ErrorWithStatus) {
	uploads, err := userfile.SaveFileUpload(r, "csvfile", userfile.CollectionCSV)
	if err != nil {
		return "", nhttp.NewError("Failed to extract image uploads: %s", err)
	}
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
	return fmt.Sprintf("/configuration/upload/%d", pool_upload.ID), nil
}
