package sync

import (
	"fmt"
	"net/http"

	"github.com/Gleipnir-Technology/nidus-sync/db/models"
	"github.com/Gleipnir-Technology/nidus-sync/html"
	"github.com/Gleipnir-Technology/nidus-sync/platform"
	"github.com/Gleipnir-Technology/nidus-sync/userfile"
)

type ContentPoolList struct {
	Uploads []platform.PoolUpload
	URL     ContentURL
	User    User
}
type ContentPoolUpload struct {
	URL  ContentURL
	User User
}

func getPoolList(w http.ResponseWriter, r *http.Request, u *models.User) {
	ctx := r.Context()
	userContent, err := contentForUser(ctx, u)
	if err != nil {
		respondError(w, "Failed to get user", err, http.StatusInternalServerError)
		return
	}
	uploads, err := platform.PoolUploadList(ctx, u.OrganizationID)
	if err != nil {
		respondError(w, "Failed to get uploads", err, http.StatusInternalServerError)
		return
	}
	data := ContentPoolList{
		Uploads: uploads,
		URL:     newContentURL(),
		User:    userContent,
	}
	html.RenderOrError(w, "sync/pool-list.html", data)
}

func getPoolUpload(w http.ResponseWriter, r *http.Request, u *models.User) {
	userContent, err := contentForUser(r.Context(), u)
	if err != nil {
		respondError(w, "Failed to get user", err, http.StatusInternalServerError)
		return
	}
	data := ContentPoolUpload{
		URL:  newContentURL(),
		User: userContent,
	}
	html.RenderOrError(w, "sync/pool-csv-upload.html", data)
}
func getPoolUploadByID(w http.ResponseWriter, r *http.Request, u *models.User) {
	userContent, err := contentForUser(r.Context(), u)
	if err != nil {
		respondError(w, "Failed to get user", err, http.StatusInternalServerError)
		return
	}
	data := ContentPoolUpload{
		URL:  newContentURL(),
		User: userContent,
	}
	html.RenderOrError(w, "sync/pool-by-id.html", data)
}
func postPoolUpload(w http.ResponseWriter, r *http.Request, u *models.User) {
	err := r.ParseMultipartForm(32 << 10) // 32 MB buffer
	if err != nil {
		respondError(w, "Failed to parse form", err, http.StatusBadRequest)
		return
	}
	uploads, err := userfile.SaveFileUpload(r, "csvfile", userfile.CollectionCSV)
	if err != nil {
		respondError(w, "Failed to extract image uploads", err, http.StatusInternalServerError)
		return
	}
	if len(uploads) == 0 {
		respondError(w, "No upload found", nil, http.StatusBadRequest)
		return
	}
	if len(uploads) != 1 {
		respondError(w, "You must only submit one file at a time", nil, http.StatusBadRequest)
		return
	}
	upload := uploads[0]
	pool_upload, err := platform.NewPoolUpload(r.Context(), u, upload)
	if err != nil {
		respondError(w, "Failed to create new pool", err, http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, fmt.Sprintf("/pool/upload/%d", pool_upload.ID), http.StatusFound)
}
