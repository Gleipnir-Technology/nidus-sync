package background

import (
	"context"
	"fmt"
	"sync"

	//commsemail "github.com/Gleipnir-Technology/nidus-sync/comms/email"
	//"github.com/Gleipnir-Technology/nidus-sync/config"
	"github.com/Gleipnir-Technology/nidus-sync/platform/email"
	"github.com/Gleipnir-Technology/nidus-sync/platform/text"
	"github.com/rs/zerolog/log"
)

var waitGroup sync.WaitGroup

func Start(ctx context.Context) {
	newOAuthTokenChannel = make(chan struct{}, 10)

	channelJobAudio = make(chan jobAudio, 100)  // Buffered channel to prevent blocking
	channelJobCSV = make(chan jobCSV, 100)      // Buffered channel to prevent blocking
	channelJobEmail = make(chan email.Job, 100) // Buffered channel to prevent blocking
	channelJobText = make(chan text.Job, 100)   // Buffered channel to prevent blocking

	/*
		waitGroup.Add(1)
		go func() {
			defer waitGroup.Done()
			commsemail.StartWebsocket(ctx, config.ForwardEmailAPIToken)
		}()
	*/

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
	err := addWaitingJobsCommit(ctx)
	if err != nil {
		return fmt.Errorf("commit: %w", err)
	}
	err = addWaitingJobsImport(ctx)
	if err != nil {
		return fmt.Errorf("commit: %w", err)
	}
	return nil
}
