package text

import (
	"context"
	"fmt"

	"github.com/Gleipnir-Technology/bob"
	"github.com/Gleipnir-Technology/nidus-sync/db/models"
	"github.com/Gleipnir-Technology/nidus-sync/platform/types"
	"github.com/rs/zerolog/log"
)

func JobRespond(ctx context.Context, txn bob.Executor, log_id int32) error {
	return respondText(ctx, txn, log_id)
}
func JobSend(ctx context.Context, txn bob.Executor, job_id int32) error {
	job, err := models.FindCommsTextJob(ctx, txn, job_id)
	if err != nil {
		return fmt.Errorf("find text: %w", err)
	}
	log.Debug().Int32("job.id", job.ID).Msg("completing text job")
	return sendTextComplete(ctx, txn, job)
}
func handleWaitingTextJobs(ctx context.Context, txn bob.Executor, dst types.E164) error {
	jobs, err := models.CommsTextJobs.Query(
		models.SelectWhere.CommsTextJobs.Destination.EQ(dst.PhoneString()),
		models.SelectWhere.CommsTextJobs.Completed.IsNull(),
	).All(ctx, txn)
	if err != nil {
		return fmt.Errorf("query jobs: %w", err)
	}
	for _, job := range jobs {
		err = sendTextComplete(ctx, txn, job)
		if err != nil {
			return fmt.Errorf("send text complete: %w", err)
		}
	}
	return nil
}
