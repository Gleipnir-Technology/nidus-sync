package rmo

import (
	"net/http"

	"github.com/Gleipnir-Technology/nidus-sync/userfile"
	"github.com/go-chi/chi/v5"
)

// ServeImageByUUID reads an image with the given UUID from disk and writes it to the HTTP response
func getImageByUUID(w http.ResponseWriter, r *http.Request) {
	uid := chi.URLParam(r, "uuid")
	if uid == "" {
		http.NotFound(w, r)
		return
	}
	userfile.PublicImageFileToResponse(w, uid)
}
