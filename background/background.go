package background

import (
	"context"
	"fmt"
	"sync"

	"github.com/Gleipnir-Technology/nidus-sync/db"
	"github.com/Gleipnir-Technology/nidus-sync/db/enums"
	"github.com/Gleipnir-Technology/nidus-sync/db/models"
	"github.com/Gleipnir-Technology/nidus-sync/platform/email"
	"github.com/Gleipnir-Technology/nidus-sync/platform/text"
	"github.com/rs/zerolog/log"
)

var waitGroup sync.WaitGroup

func Start(ctx context.Context) {
	newOAuthTokenChannel = make(chan struct{}, 10)

	channelJobAudio = make(chan jobAudio, 100)                 // Buffered channel to prevent blocking
	channelJobImportCSVPool = make(chan jobImportCSVPool, 100) // Buffered channel to prevent blocking
	channelJobEmail = make(chan email.Job, 100)                // Buffered channel to prevent blocking
	channelJobText = make(chan text.Job, 100)                  // Buffered channel to prevent blocking

	waitGroup.Add(1)
	go func() {
		defer waitGroup.Done()
		refreshFieldseekerData(ctx, newOAuthTokenChannel)
	}()

	waitGroup.Add(1)
	go func() {
		defer waitGroup.Done()
		startWorkerAudio(ctx, channelJobAudio)
	}()

	waitGroup.Add(1)
	go func() {
		defer waitGroup.Done()
		startWorkerCSV(ctx, channelJobImportCSVPool)
	}()

	waitGroup.Add(1)
	go func() {
		defer waitGroup.Done()
		startWorkerEmail(ctx, channelJobEmail)
	}()

	waitGroup.Add(1)
	go func() {
		defer waitGroup.Done()
		startWorkerText(ctx, channelJobText)
	}()

	err := addWaitingJobs(ctx)
	if err != nil {
		log.Error().Err(err).Msg("Failed to add waiting background jobs")
	}
}

func WaitForExit() {

	waitGroup.Wait()
}

func addWaitingJobs(ctx context.Context) error {
	rows, err := models.FileuploadFiles.Query(
		models.SelectWhere.FileuploadFiles.Status.EQ(
			enums.FileuploadFilestatustypeUploaded,
		),
	).All(ctx, db.PGInstance.BobDB)
	if err != nil {
		return fmt.Errorf("Failed to query file uploads: %w", err)
	}
	for _, row := range rows {
		report_id := row.ID
		job := jobImportCSVPool{
			fileID: report_id,
		}
		select {
		case channelJobImportCSVPool <- job:
			log.Info().Int32("report_id", report_id).Msg("CSV upload job queued")
		default:
			log.Warn().Int32("report_id", report_id).Msg("CSV upload job failed to queue, channel full")
		}
	}
	return nil
}
