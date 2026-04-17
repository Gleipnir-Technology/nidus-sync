package text

import (
	"context"
	"fmt"

	"github.com/Gleipnir-Technology/nidus-sync/db"
	"github.com/Gleipnir-Technology/nidus-sync/db/models"
	"github.com/Gleipnir-Technology/nidus-sync/platform/types"
	//"github.com/rs/zerolog/log"
)

func JobRespond(ctx context.Context, log_id int32) error {
	return respondText(ctx, log_id)
}
func JobSend(ctx context.Context, job_id int32) error {
	bxn := db.PGInstance.BobDB
	job, err := models.FindCommsTextJob(ctx, bxn, job_id)
	if err != nil {
		return fmt.Errorf("find text: %w", err)
	}
	//log.Debug().Int32("job.id", job.ID).Msg("completing text job")
	return sendTextComplete(ctx, job)
}
func handleWaitingTextJobs(ctx context.Context, dst types.E164) error {
	bxn := db.PGInstance.BobDB
	jobs, err := models.CommsTextJobs.Query(
		models.SelectWhere.CommsTextJobs.Destination.EQ(dst.PhoneString()),
		models.SelectWhere.CommsTextJobs.Completed.IsNull(),
	).All(ctx, bxn)
	if err != nil {
		return fmt.Errorf("query jobs: %w", err)
	}
	for _, job := range jobs {
		err = sendTextComplete(ctx, job)
		if err != nil {
			return fmt.Errorf("send text complete: %w", err)
		}
	}
	return nil
}
