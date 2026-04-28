package api

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/Gleipnir-Technology/nidus-sync/db"
	"github.com/Gleipnir-Technology/nidus-sync/db/models"
	"github.com/Gleipnir-Technology/nidus-sync/platform"
	"github.com/Gleipnir-Technology/nidus-sync/platform/background"
	"github.com/Gleipnir-Technology/nidus-sync/platform/file"
	"github.com/aarondl/opt/omit"
	"github.com/aarondl/opt/omitnull"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/rs/zerolog/log"
)

func apiAudioPost(w http.ResponseWriter, r *http.Request, u platform.User) {
	vars := mux.Vars(r)
	id := vars["uuid"]
	noteUUID, err := uuid.Parse(id)
	if err != nil {
		http.Error(w, "Failed to decode the uuid", http.StatusBadRequest)
		return
	}

	var payload NoteAudioPayload
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read the payload", http.StatusBadRequest)
		return
	}
	if err := json.Unmarshal(body, &payload); err != nil {
		//debugSaveRequest(body, err, "Audio note POST JSON decode error")
		http.Error(w, "Failed to decode the payload", http.StatusBadRequest)
		return
	}
	ctx := r.Context()
	setter := models.NoteAudioSetter{
		Created:                 omit.From(payload.Created),
		CreatorID:               omit.From(int32(u.ID)),
		Deleted:                 omitnull.FromPtr(payload.Deleted),
		DeletorID:               omitnull.FromPtr(payload.DeletorID),
		Duration:                omit.From(payload.Duration),
		OrganizationID:          omit.From(u.Organization.ID),
		Transcription:           omitnull.FromPtr(payload.Transcription),
		TranscriptionUserEdited: omit.From(payload.TranscriptionUserEdited),
		Version:                 omit.From(payload.Version),
		UUID:                    omit.From(noteUUID),
	}
	if err := platform.NoteAudioCreate(ctx, u, setter); err != nil {
		renderShim(w, r, errRender(err))
		return
	}
	w.WriteHeader(http.StatusAccepted)
}

func apiAudioContentPost(w http.ResponseWriter, r *http.Request, user platform.User) {
	vars := mux.Vars(r)
	u_str := vars["uuid"]
	u, err := uuid.Parse(u_str)
	if err != nil {
		http.Error(w, "Failed to parse image UUID", http.StatusBadRequest)
		return
	}
	err = file.FileContentWrite(r.Body, file.CollectionAudioRaw, u)
	if err != nil {
		log.Printf("Failed to write content file: %v", err)
		http.Error(w, "failed to write content file", http.StatusInternalServerError)
	}
	ctx := r.Context()
	a, err := models.NoteAudios.Query(
		models.SelectWhere.NoteAudios.UUID.EQ(u),
		models.SelectWhere.NoteAudios.OrganizationID.EQ(user.Organization.ID),
	).One(ctx, db.PGInstance.BobDB)
	if err != nil {
		log.Printf("Failed to get note audio %s for org %d: %w", u_str, user.Organization.ID, err)
		http.Error(w, "failed to update database", http.StatusBadRequest)
	}

	background.NewAudioTranscode(ctx, db.PGInstance.BobDB, a.ID)
	w.WriteHeader(http.StatusOK)
}
