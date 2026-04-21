package platform

import (
	"context"
	"fmt"
	"strconv"
	"sync"
	"time"

	//"github.com/Gleipnir-Technology/bob/dialect/psql"
	//"github.com/Gleipnir-Technology/bob/dialect/psql/sm"
	"github.com/Gleipnir-Technology/nidus-sync/config"
	"github.com/Gleipnir-Technology/nidus-sync/db"
	"github.com/Gleipnir-Technology/nidus-sync/db/enums"
	"github.com/Gleipnir-Technology/nidus-sync/db/models"
	//"github.com/Gleipnir-Technology/nidus-sync/platform/background"
	"github.com/Gleipnir-Technology/nidus-sync/platform/csv"
	"github.com/Gleipnir-Technology/nidus-sync/platform/email"
	"github.com/Gleipnir-Technology/nidus-sync/platform/file"
	"github.com/Gleipnir-Technology/nidus-sync/platform/geocode"
	"github.com/Gleipnir-Technology/nidus-sync/platform/mailer"
	"github.com/Gleipnir-Technology/nidus-sync/platform/text"
	//"github.com/Gleipnir-Technology/nidus-sync/userfile"
	//"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	bobpgx "github.com/stephenafamo/bob/drivers/pgx"
)

var waitGroup sync.WaitGroup
var newOAuthTokenChannel chan struct{}

func StartAll(ctx context.Context) error {
	err := email.LoadTemplates()
	if err != nil {
		return fmt.Errorf("Failed to load email templates: %w", err)
	}

	err = text.StoreSources()
	if err != nil {
		return fmt.Errorf("Failed to store text source phone numbers: %w", err)
	}

	err = file.CreateDirectories()
	if err != nil {
		return fmt.Errorf("Failed to create file directories: %w", err)
	}

	err = initializeLabelStudio()
	if err != nil {
		return fmt.Errorf("init label studio: %w", err)
	}

	geocode.InitializeStadia(config.StadiaMapsAPIKey)

	newOAuthTokenChannel = make(chan struct{}, 10)

	waitGroup.Add(1)
	go func() {
		defer waitGroup.Done()
		refreshFieldseekerData(ctx, newOAuthTokenChannel)
	}()
	waitGroup.Add(1)
	go func() {
		defer waitGroup.Done()
		listenForJobs(ctx)
	}()

	err = addWaitingJobs(ctx)
	if err != nil {
		log.Error().Err(err).Msg("Failed to add waiting background jobs")
	}
	return nil
}

func WaitForExit() {
	waitGroup.Wait()
}

func addWaitingJobs(ctx context.Context) error {
	jobs, err := models.Jobs.Query().All(ctx, db.PGInstance.BobDB)
	if err != nil {
		return fmt.Errorf("Failed to query waiting jobs: %w", err)
	}
	go func() {
		for _, job := range jobs {
			sublog := log.With().Int32("job", job.ID).Int32("row_id", job.RowID).Str("type", string(job.Type)).Logger()
			sublog.Info().Msg("begin restarted background job")
			err = handleJob(ctx, job)
			if err != nil {
				sublog.Error().Err(err).Msg("failed handle job")
				continue
			}
			err = job.Delete(ctx, db.PGInstance.BobDB)
			if err != nil {
				sublog.Error().Err(err).Msg("failed delete job")
				continue
			}
			sublog.Info().Msg("job complete")
		}
	}()
	return nil
}
func handleJob(ctx context.Context, job *models.Job) error {
	switch job.Type {
	case enums.JobtypeAudioTranscode:
		return processAudioFile(ctx, job.RowID)
	case enums.JobtypeComplianceMailerSend:
		return mailer.ComplianceSend(ctx, job.RowID)
	case enums.JobtypeCSVCommit:
		return csv.JobCommit(ctx, job.RowID)
	case enums.JobtypeCSVImport:
		return csv.JobImport(ctx, job.RowID)
	case enums.JobtypeLabelStudioAudioCreate:
		return jobLabelStudioAudioCreate(ctx, job.RowID)
	case enums.JobtypeEmailSend:
		return email.Job(ctx, job.RowID)
	case enums.JobtypeTextRespond:
		return text.JobRespond(ctx, job.RowID)
	case enums.JobtypeTextSend:
		return text.JobSend(ctx, job.RowID)
	default:
		return fmt.Errorf("No handler for job type %s", string(job.Type))
	}
}
func listenForJobs(ctx context.Context) {
	for {
		//es.SendQueuedEmails(ctx) // send any emails queued prior to listening for notificiations
		err := listenAndDoOneJob(ctx)
		if err != nil {
			if err.Error() == "context canceled" {
				log.Debug().Msg("Exiting listenForJobs")
				return
			}
			log.Error().Err(err).Msg("Crashed listenAndDoOneJob")
		}

		select {
		case <-ctx.Done():
			return
		default:
			// If listenAndSendOneConn returned and ctx has not been cancelled that means there was a fatal database error.
			// Wait a while to avoid busy-looping while the database is unreachable.
			time.Sleep(time.Minute)
		}
	}
}
func listenAndDoOneJob(ctx context.Context) error {
	conn, err := db.PGInstance.PGXPool.Acquire(ctx)
	if err != nil {
		//if !pgconn.Timeout(err) {
		return fmt.Errorf("failed to acquire database connection to listen for queued emails: %w", err)
	}
	defer conn.Release()

	_, err = conn.Exec(ctx, "LISTEN new_job")
	if err != nil {
		//if !pgconn.Timeout(err) {
		return fmt.Errorf("failed to execute 'LISTEN new_job': %w", err)
	}

	for {
		//log.Debug().Msg("wait for notification")
		notification, err := conn.Conn().WaitForNotification(ctx)
		if err != nil {
			//if !pgconn.Timeout(err) {
			return fmt.Errorf("failed while waiting for notification of new job: %w", err)
		}

		job_id, err := strconv.Atoi(notification.Payload)
		if err != nil {
			return fmt.Errorf("failed to parse int from payload '%s': %w", notification.Payload, err)
		}
		//log.Debug().Int("job_id", job_id).Msg("got notification for job")

		c := bobpgx.NewConn(conn.Conn())
		job, err := models.FindJob(ctx, c, int32(job_id))
		if err != nil {
			return fmt.Errorf("Failed to find job %d: %w", job_id, err)
		}
		sublog := log.With().Int32("job", job.ID).Int32("row_id", job.RowID).Str("type", string(job.Type)).Logger()

		err = handleJob(ctx, job)
		if err != nil {
			sublog.Error().Err(err).Msg("failed to handle job")
			return nil
		}
		err = job.Delete(ctx, db.PGInstance.BobDB)
		if err != nil {
			sublog.Error().Err(err).Msg("failed to delete job")
			return fmt.Errorf("delete job: %w", err)
		}
		sublog.Debug().Msg("job complete")
	}
}
