package platform

import (
	"context"
	"fmt"

	"github.com/Gleipnir-Technology/bob"
	"github.com/Gleipnir-Technology/nidus-sync/db"
	"github.com/Gleipnir-Technology/nidus-sync/db/models"
	//"github.com/Gleipnir-Technology/nidus-sync/platform/background"
	"github.com/Gleipnir-Technology/nidus-sync/platform/subprocess"
	//"github.com/google/uuid"
	//"github.com/rs/zerolog/log"
)

func processAudioFile(ctx context.Context, txn bob.Executor, audio_id int32) error {
	a, err := models.NoteAudios.Query(
		models.SelectWhere.NoteAudios.ID.EQ(audio_id),
	).One(ctx, db.PGInstance.BobDB)

	if err != nil {
		return fmt.Errorf("note audio query: %w", err)
	}
	// Normalize audio
	err = subprocess.NormalizeAudio(a.UUID)
	if err != nil {
		return fmt.Errorf("failed to normalize audio %s: %v", a.UUID, err)
	}

	// Transcode to OGG
	err = subprocess.TranscodeToOgg(a.UUID)
	if err != nil {
		return fmt.Errorf("failed to transcode audio %s to OGG: %v", a.UUID, err)
	}

	//background.NewLabelStudioAudioCreate(ctx, db.PGInstance.BobDB, audio_id)
	return nil
}
