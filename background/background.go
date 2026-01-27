package background

import (
	"context"
	"sync"

	"github.com/Gleipnir-Technology/nidus-sync/comms/email"
	"github.com/Gleipnir-Technology/nidus-sync/platform/text"
)

var waitGroup sync.WaitGroup

func Start(ctx context.Context) {
	newOAuthTokenChannel = make(chan struct{}, 10)

	channelJobAudio = make(chan jobAudio, 100)  // Buffered channel to prevent blocking
	channelJobEmail = make(chan email.Job, 100) // Buffered channel to prevent blocking
	channelJobText = make(chan text.Job, 100)   // Buffered channel to prevent blocking

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
		startWorkerEmail(ctx, channelJobEmail)
	}()

	waitGroup.Add(1)
	go func() {
		defer waitGroup.Done()
		startWorkerText(ctx, channelJobText)
	}()
}
func WaitForExit() {

	waitGroup.Wait()
}
