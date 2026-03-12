package platform

import (
	"context"

	"github.com/Gleipnir-Technology/nidus-sync/platform/background"
)

func BackgroundStart(ctx context.Context) {
	background.Start(ctx)
}
func BackgroundWaitForExit() {
	background.WaitForExit()
}
