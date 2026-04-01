package api

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/Gleipnir-Technology/nidus-sync/db/models"
	"github.com/Gleipnir-Technology/nidus-sync/platform"
	"github.com/Gleipnir-Technology/nidus-sync/platform/file"
	"github.com/aarondl/opt/omit"
	"github.com/aarondl/opt/omitnull"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/rs/zerolog/log"
)

func apiImagePost(w http.ResponseWriter, r *http.Request, u platform.User) {
	vars := mux.Vars(r)
	id := vars["uuid"]
	noteUUID, err := uuid.Parse(id)
	if err != nil {
		http.Error(w, "Failed to decode the uuid", http.StatusBadRequest)
		return
	}

	var payload NoteImagePayload
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read the payload", http.StatusBadRequest)
		return
	}
	if err := json.Unmarshal(body, &payload); err != nil {
		//debugSaveRequest(body, err, "Image note POST JSON decode error")
		http.Error(w, "Failed to decode the payload", http.StatusBadRequest)
		return
	}
	ctx := r.Context()
	setter := models.NoteImageSetter{
		Created:        omit.From(payload.Created),
		CreatorID:      omit.From(int32(u.ID)),
		Deleted:        omitnull.FromPtr(payload.Deleted),
		DeletorID:      omitnull.FromPtr(payload.DeletorID),
		OrganizationID: omit.From(u.Organization.ID),
		Version:        omit.From(payload.Version),
		UUID:           omit.From(noteUUID),
	}
	err = platform.NoteImageCreate(ctx, u, setter)
	if err != nil {
		renderShim(w, r, errRender(err))
		return
	}
	w.WriteHeader(http.StatusAccepted)
}

func apiImageContentGet(w http.ResponseWriter, r *http.Request, u platform.User) {
	vars := mux.Vars(r)
	u_str := vars["uuid"]
	imageUUID, err := uuid.Parse(u_str)
	if err != nil {
		log.Error().Err(err).Msg("Failed to parse image UUID")
		http.Error(w, "Failed to parse image UUID", http.StatusBadRequest)
	}
	file.PublicImageFileToResponse(w, imageUUID)
	w.WriteHeader(http.StatusOK)
}
func apiImageContentPost(w http.ResponseWriter, r *http.Request, u platform.User) {
	vars := mux.Vars(r)
	u_str := vars["uuid"]
	imageUUID, err := uuid.Parse(u_str)
	if err != nil {
		log.Error().Err(err).Msg("Failed to parse image UUID")
		http.Error(w, "Failed to parse image UUID", http.StatusBadRequest)
	}
	err = file.ImageFileContentWrite(imageUUID, r.Body)
	if err != nil {
		renderShim(w, r, errRender(err))
		return
	}
	w.WriteHeader(http.StatusOK)
	log.Printf("Saved image file %s\n", imageUUID)
	fmt.Fprintf(w, "PNG uploaded successfully")
}
