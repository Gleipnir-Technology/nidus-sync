package background

import (
	"context"
	//"fmt"

	"github.com/Gleipnir-Technology/nidus-sync/platform/csv"
	//"github.com/Gleipnir-Technology/nidus-sync/userfile"
	//"github.com/google/uuid"
	"github.com/rs/zerolog/log"
)

// represents a job to import a pool CSV file
type jobImportCSVPool struct {
	fileID int32
}

var channelJobImportCSVPool chan jobImportCSVPool

func startWorkerCSV(ctx context.Context, channelJobImport chan jobImportCSVPool) {
	go func() {
		for {
			select {
			case <-ctx.Done():
				log.Info().Msg("CSV worker shutting down.")
				return
			case job := <-channelJobImport:
				log.Info().Int32("id", job.fileID).Msg("Processing CSV job")
				err := csv.ProcessJob(ctx, job.fileID)
				if err != nil {
					log.Error().Err(err).Int32("id", job.fileID).Msg("Error processing CSV file")
				}
			}
		}
	}()
}
