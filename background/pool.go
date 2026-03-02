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

type jobImportCSV struct {
	fileID int32
	type_  enums.FileuploadCsvtype
}

var channelJobImportCSV chan jobImportCSV

func ProcessUpload(file_id int32, t enums.FileuploadCsvtype) {
	enqueueUploadJob(jobImportCSV{
		fileID: file_id,
		type_:  t,
	})
}

func enqueueUploadJob(job jobImportCSV) {
	select {
	case channelJobImportCSV <- job:
		log.Info().Int32("file_id", job.fileID).Msg("Enqueued csv job")
	default:
		log.Warn().Int32("file_id", job.fileID).Msg("csv channel is full, dropping job")
	}
}
func startWorkerCSV(ctx context.Context, channelJobImport chan jobImportCSV) {
	go func() {
		for {
			select {
			case <-ctx.Done():
				log.Info().Msg("CSV worker shutting down.")
				return
			case job := <-channelJobImport:
				log.Info().Int32("id", job.fileID).Msg("Processing CSV job")
				err := csv.ProcessJob(ctx, job.fileID, job.type_)
				if err != nil {
					log.Error().Err(err).Int32("id", job.fileID).Msg("Error processing CSV file")
				}
			}
		}
	}()
}
