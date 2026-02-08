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
	URL  ContentURL
	User User
}
type ContentPoolUpload struct {
	URL  ContentURL
	User User
}

func getPoolList(w http.ResponseWriter, r *http.Request, u *models.User) {
	userContent, err := contentForUser(r.Context(), u)
	if err != nil {
		respondError(w, "Failed to get user", err, http.StatusInternalServerError)
		return
	}
	data := ContentPoolList{
		URL:  newContentURL(),
		User: userContent,
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
func postPoolUpload(w http.ResponseWriter, r *http.Request, u *models.User) {
	err := r.ParseMultipartForm(32 << 10) // 32 MB buffer
	if err != nil {
		respondError(w, "Failed to parse form", err, http.StatusBadRequest)
		return
	}
	uploads, err := userfile.SaveFileUpload(r, "csvfile", "pool", "csv")
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
	http.Redirect(w, r, fmt.Sprintf("/pool/upload/%d", pool_upload.ID), http.StatusFound)
}
