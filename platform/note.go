package platform

import (
	"context"

	"github.com/Gleipnir-Technology/nidus-sync/db"
	"github.com/Gleipnir-Technology/nidus-sync/db/models"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
)

func NoteAudioCreate(ctx context.Context, user User, setter models.NoteAudioSetter) error {
	err := user.Organization.model.InsertNoteAudios(ctx, db.PGInstance.BobDB, &setter)
	if err == nil {
		return nil
	}
	// Just ignore this failure, it means we already have this content
	if err.Error() == "insertOrganizationNoteAudios0: ERROR: duplicate key value violates unique constraint \"note_audio_pkey\" (SQLSTATE 23505)" {
		return nil
	}
	log.Warn().Err(err).Msg("Unrecognized error creating note audio")
	return err
}

func NoteAudioNormalized(uuid string) error {
	return nil
}
func NoteAudioTranscodedToOgg(uuid string) error {
	return nil
}
func NoteImageCreate(ctx context.Context, user User, setter models.NoteImageSetter) error {
	err := user.Organization.model.InsertNoteImages(ctx, db.PGInstance.BobDB, &setter)
	if err == nil {
		return nil
	}
	// Just ignore this failure, it means we already have this content
	if err.Error() == "insertOrganizationNoteImages0: ERROR: duplicate key value violates unique constraint \"note_image_pkey\" (SQLSTATE 23505)" {
		return nil
	}
	log.Warn().Err(err).Msg("Unrecognized error creating note audio")
	return err
}

func NoteUpdate(ctx context.Context, noteUUID uuid.UUID, setter models.NoteAudioSetter) error {
	return nil
}
