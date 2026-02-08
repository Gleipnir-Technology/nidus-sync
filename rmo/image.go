package rmo

import (
	"net/http"

	"github.com/Gleipnir-Technology/nidus-sync/userfile"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

// ServeImageByUUID reads an image with the given UUID from disk and writes it to the HTTP response
func getImageByUUID(w http.ResponseWriter, r *http.Request) {
	u := chi.URLParam(r, "uuid")
	if u == "" {
		http.NotFound(w, r)
		return
	}
	uid, err := uuid.Parse(u)
	if err != nil {
		http.Error(w, "Failed to parse uuid", http.StatusBadRequest)
		return
	}
	userfile.PublicImageFileToResponse(w, uid)
}
