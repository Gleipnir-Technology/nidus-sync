package background

import (
	"context"
	"fmt"
	"time"

	"github.com/Gleipnir-Technology/bob"
	"github.com/Gleipnir-Technology/nidus-sync/db"
	"github.com/Gleipnir-Technology/nidus-sync/db/enums"
	"github.com/Gleipnir-Technology/nidus-sync/db/gen/nidus-sync/public/model"
	"github.com/Gleipnir-Technology/nidus-sync/db/models"
	query "github.com/Gleipnir-Technology/nidus-sync/db/query/public"
	"github.com/aarondl/opt/omit"
	//"github.com/rs/zerolog/log"
)

func NewAudioTranscode(ctx context.Context, txn bob.Executor, audio_id int32) error {
	return newJob(ctx, txn, enums.JobtypeAudioTranscode, audio_id)
}
func NewComplianceMailer(ctx context.Context, txn db.Ex, compliance_report_request_id int32) error {
	return newJob2(ctx, txn, model.Jobtype_ComplianceMailerSend, compliance_report_request_id)
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
func NewTextRespond(ctx context.Context, txn bob.Executor, text_id int32) error {
	return newJob(ctx, txn, enums.JobtypeTextRespond, text_id)
}
func NewTextSend(ctx context.Context, txn bob.Executor, job_id int32) error {
	return newJob(ctx, txn, enums.JobtypeTextSend, job_id)
}
func newJob(ctx context.Context, txn bob.Executor, t enums.Jobtype, id int32) error {
	_, err := models.Jobs.Insert(&models.JobSetter{
		Created: omit.From(time.Now()),
		// ID
		Type:  omit.From(t),
		RowID: omit.From(id),
	}).One(ctx, txn)
	if err != nil {
		return fmt.Errorf("insert job: %w", err)
	}
	return nil
}
func newJob2(ctx context.Context, txn db.Ex, t model.Jobtype, id int32) error {
	job := model.Job{
		Created: time.Now(),
		Type:    t,
		RowID:   id,
	}
	_, err := query.JobInsert(ctx, txn, job)
	if err != nil {
		return fmt.Errorf("insert job: %w", err)
	}
	return nil
}
