package sync

import (
	"github.com/Gleipnir-Technology/nidus-sync/db/models"
	"github.com/Gleipnir-Technology/nidus-sync/html"
	"github.com/Gleipnir-Technology/nidus-sync/userfile"
	"net/http"
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
	images, err := saveImageUploads(r.Context(), tx, uploads)
}
