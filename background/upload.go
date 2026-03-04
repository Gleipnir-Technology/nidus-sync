package background

import (
	"context"
	//"fmt"

	"github.com/Gleipnir-Technology/nidus-sync/db/enums"
	"github.com/Gleipnir-Technology/nidus-sync/platform/csv"
	//"github.com/Gleipnir-Technology/nidus-sync/userfile"
	//"github.com/google/uuid"
	"github.com/rs/zerolog/log"
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
				log.Info().Int32("id", job.fileID).Msg("Processing CSV job")
				switch job.action {
				case jobCSVActionCommit:
					err := csv.JobCommit(ctx, job.fileID)
					if err != nil {
						log.Error().Err(err).Int32("id", job.fileID).Msg("Error processing CSV file")
						continue
					}
				case jobCSVActionImport:
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
