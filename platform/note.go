package platform

import (
	"context"
	"fmt"
	"strconv"

	"github.com/Gleipnir-Technology/nidus-sync/db"
	"github.com/Gleipnir-Technology/nidus-sync/db/models"
	"github.com/Gleipnir-Technology/nidus-sync/platform/event"
	//"github.com/google/uuid"
	//"github.com/rs/zerolog/log"
)

func NoteAudioCreate(ctx context.Context, user User, setter models.NoteAudioSetter) error {
	txn, err := db.PGInstance.BobDB.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("create txn: %w", err)
	}
	defer txn.Rollback(ctx)

	note_audio, err := models.NoteAudios.Insert(&setter).One(ctx, txn)
	if err != nil {
		// Just ignore this failure, it means we already have this content
		if err.Error() != "insertOrganizationNoteAudios0: ERROR: duplicate key value violates unique constraint \"note_audio_pkey\" (SQLSTATE 23505)" {
			return fmt.Errorf("create note_audio: %w", err)
		}
	}
	event.Created(event.TypeNoteAudio, user.Organization.ID, strconv.Itoa(int(note_audio.ID)))
	txn.Commit(ctx)

	return nil
}

func NoteImageCreate(ctx context.Context, user User, setter models.NoteImageSetter) error {
	txn, err := db.PGInstance.BobDB.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("create txn: %w", err)
	}
	defer txn.Rollback(ctx)
	note_image, err := models.NoteImages.Insert(&setter).One(ctx, db.PGInstance.BobDB)
	if err != nil {
		// Just ignore this failure, it means we already have this content
		if err.Error() != "insertOrganizationNoteImages0: ERROR: duplicate key value violates unique constraint \"note_image_pkey\" (SQLSTATE 23505)" {
			return fmt.Errorf("create note_image: %w", err)
		}
	}
	event.Created(event.TypeNoteImage, user.Organization.ID, strconv.Itoa(int(note_image.ID)))
	txn.Commit(ctx)

	return err
}
