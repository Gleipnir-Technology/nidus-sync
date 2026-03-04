package background

import (
	"context"
	"fmt"
	"sync"

	"github.com/Gleipnir-Technology/bob"
	"github.com/Gleipnir-Technology/bob/dialect/psql"
	"github.com/Gleipnir-Technology/bob/dialect/psql/sm"
	commsemail "github.com/Gleipnir-Technology/nidus-sync/comms/email"
	"github.com/Gleipnir-Technology/nidus-sync/config"
	"github.com/Gleipnir-Technology/nidus-sync/db"
	"github.com/Gleipnir-Technology/nidus-sync/db/enums"
	"github.com/Gleipnir-Technology/nidus-sync/platform/email"
	"github.com/Gleipnir-Technology/nidus-sync/platform/text"
	"github.com/rs/zerolog/log"
	"github.com/stephenafamo/scan"
)

var waitGroup sync.WaitGroup

func Start(ctx context.Context) {
	newOAuthTokenChannel = make(chan struct{}, 10)

	channelJobAudio = make(chan jobAudio, 100)  // Buffered channel to prevent blocking
	channelJobCSV = make(chan jobCSV, 100)      // Buffered channel to prevent blocking
	channelJobEmail = make(chan email.Job, 100) // Buffered channel to prevent blocking
	channelJobText = make(chan text.Job, 100)   // Buffered channel to prevent blocking

	waitGroup.Add(1)
	go func() {
		defer waitGroup.Done()
		commsemail.StartWebsocket(ctx, config.ForwardEmailAPIToken)
	}()

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
		startWorkerCSV(ctx, channelJobCSV)
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
			psql.Raw("file.status").EQ(psql.Arg(enums.FileuploadFilestatustypeUploaded)),
		),
	), scan.StructMapper[Row_]())

	if err != nil {
		return fmt.Errorf("Failed to query file uploads: %w", err)
	}
	for _, row := range rows {
		report_id := row.ID
		enqueueJobCSV(jobCSV{
			fileID:  report_id,
			csvType: row.Type,
		})
	}
	return nil
}
