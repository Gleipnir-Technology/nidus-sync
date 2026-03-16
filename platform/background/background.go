package background

import (
	"context"
	"fmt"

	"github.com/Gleipnir-Technology/bob"
	"github.com/Gleipnir-Technology/nidus-sync/db/enums"
	"github.com/Gleipnir-Technology/nidus-sync/db/models"
	"github.com/aarondl/opt/omit"
	//"github.com/rs/zerolog/log"
)

func NewAudioTranscode(ctx context.Context, txn bob.Executor, audio_id int32) error {
	return newJob(ctx, txn, enums.JobtypeCSVCommit, audio_id)
}
func NewCSVCommit(ctx context.Context, txn bob.Executor, csv_id int32) error {
	return newJob(ctx, txn, enums.JobtypeCSVCommit, csv_id)
}
func NewCSVImport(ctx context.Context, txn bob.Executor, csv_id int32) error {
	return newJob(ctx, txn, enums.JobtypeCSVImport, csv_id)
}
func NewEmailSend(ctx context.Context, txn bob.Executor, email_id int32) error {
	return newJob(ctx, txn, enums.JobtypeEmailSend, email_id)
}
func NewLabelStudioAudioCreate(ctx context.Context, txn bob.Executor, note_audio_id int32) error {
	return newJob(ctx, txn, enums.JobtypeLabelStudioAudioCreate, note_audio_id)
}
func NewTextSend(ctx context.Context, txn bob.Executor, text_id int32) error {
	return newJob(ctx, txn, enums.JobtypeTextSend, text_id)
}
func newJob(ctx context.Context, txn bob.Executor, t enums.Jobtype, id int32) error {
	_, err := models.Jobs.Insert(&models.JobSetter{
		// ID
		Type:  omit.From(t),
		RowID: omit.From(id),
	}).One(ctx, txn)
	if err != nil {
		return fmt.Errorf("insert job: %w", err)
	}
	return nil
}
