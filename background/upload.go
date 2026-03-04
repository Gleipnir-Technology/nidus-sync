package background

import (
	"context"
	"fmt"

	"github.com/Gleipnir-Technology/bob"
	"github.com/Gleipnir-Technology/bob/dialect/psql"
	"github.com/Gleipnir-Technology/bob/dialect/psql/sm"
	"github.com/Gleipnir-Technology/nidus-sync/db"
	"github.com/Gleipnir-Technology/nidus-sync/db/enums"
	"github.com/Gleipnir-Technology/nidus-sync/platform/csv"
	//"github.com/Gleipnir-Technology/nidus-sync/userfile"
	//"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"github.com/stephenafamo/scan"
)

type jobCSVAction = int

const (
	jobCSVActionCommit jobCSVAction = iota
	jobCSVActionImport
)

type jobCSV struct {
	action  jobCSVAction
	csvType enums.FileuploadCsvtype
	fileID  int32
}

var channelJobCSV chan jobCSV

func CommitUpload(file_id int32) {
	enqueueJobCSV(jobCSV{
		action: jobCSVActionCommit,
		fileID: file_id,
	})
}
func ProcessUpload(file_id int32, t enums.FileuploadCsvtype) {
	enqueueJobCSV(jobCSV{
		csvType: t,
		fileID:  file_id,
	})
}

func addWaitingJobsCommit(ctx context.Context) error {
	return addWaitingJobsForType(ctx, enums.FileuploadFilestatustypeCommitting, jobCSVActionCommit)
}
func addWaitingJobsImport(ctx context.Context) error {
	return addWaitingJobsForType(ctx, enums.FileuploadFilestatustypeUploaded, jobCSVActionImport)
}
func addWaitingJobsForType(ctx context.Context, status enums.FileuploadFilestatustype, action jobCSVAction) error {
	type Row_ struct {
		ID   int32                   `db:"id"`
		Type enums.FileuploadCsvtype `db:"type"`
	}
	rows, err := bob.All(ctx, db.PGInstance.BobDB, psql.Select(
		sm.Columns(
			"file.id AS id",
			"csv.type_ AS type",
		),
		sm.From("fileupload.file").As("file"),
		sm.InnerJoin("fileupload.csv").As("csv").OnEQ(psql.Raw("file.id"), psql.Raw("csv.file_id")),
		sm.Where(
			psql.Raw("file.status").EQ(psql.Arg(status)),
		),
	), scan.StructMapper[Row_]())

	if err != nil {
		return fmt.Errorf("Failed to query file uploads: %w", err)
	}
	for _, row := range rows {
		report_id := row.ID
		enqueueJobCSV(jobCSV{
			action:  action,
			fileID:  report_id,
			csvType: row.Type,
		})
	}
	return nil
}
func enqueueJobCSV(job jobCSV) {
	select {
	case channelJobCSV <- job:
		log.Info().Int32("file_id", job.fileID).Msg("Enqueued csv job")
	default:
		log.Warn().Int32("file_id", job.fileID).Msg("csv channel is full, dropping job")
	}
}
func startWorkerCSV(ctx context.Context, channelJobImport chan jobCSV) {
	go func() {
		for {
			select {
			case <-ctx.Done():
				log.Info().Msg("CSV worker shutting down.")
				return
			case job := <-channelJobImport:
				switch job.action {
				case jobCSVActionCommit:
					log.Info().Int32("id", job.fileID).Msg("Processing CSV commit job")
					err := csv.JobCommit(ctx, job.fileID)
					if err != nil {
						log.Error().Err(err).Int32("id", job.fileID).Msg("Error processing CSV file")
						continue
					}
				case jobCSVActionImport:
					log.Info().Int32("id", job.fileID).Msg("Processing CSV import job")
					err := csv.JobImport(ctx, job.fileID, job.csvType)
					if err != nil {
						log.Error().Err(err).Int32("id", job.fileID).Msg("Error processing CSV file")
						continue
					}
				default:
					log.Error().Msg("Unrecognized job action")
					return
				}
				log.Info().Int32("id", job.fileID).Msg("Done processing CSV job")
			}
		}
	}()
}
