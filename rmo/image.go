package rmo

import (
	"net/http"

	"github.com/Gleipnir-Technology/nidus-sync/platform/file"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

// ServeImageByUUID reads an image with the given UUID from disk and writes it to the HTTP response
func getImageByUUID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	u := vars["uuid"]
	if u == "" {
		http.NotFound(w, r)
		return
	}
	uid, err := uuid.Parse(u)
	if err != nil {
		http.Error(w, "Failed to parse uuid", http.StatusBadRequest)
		return
	}
	file.ImageFileToWriter(file.CollectionPublicImage, uid, w)
}
