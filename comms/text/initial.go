package text

import (
	"context"
	"fmt"

	"github.com/Gleipnir-Technology/nidus-sync/db"
	"github.com/Gleipnir-Technology/nidus-sync/db/enums"
	"github.com/Gleipnir-Technology/nidus-sync/db/models"
)

func ensureInitialText(ctx context.Context, src string, dst string) error {
	//
	origin := enums.CommsTextoriginWebsiteAction
	rows, err := models.CommsTextLogs.Query(
		models.SelectWhere.CommsTextLogs.Destination.EQ(dst),
		models.SelectWhere.CommsTextLogs.IsWelcome.EQ(true),
	).All(ctx, db.PGInstance.BobDB)
	if err != nil {
		return fmt.Errorf("Failed to query text logs: %w", err)
	}
	if len(rows) > 0 {
		return nil
	}
	content := "Welcome to Report Mosquitoes Online. We received your request and want to confirm text updates. Reply YES to continue. Reply STOP at any time to unsubscribe"
	err = sendText(ctx, src, dst, content, origin)
	if err != nil {
		return fmt.Errorf("Failed to send initial confirmation: %w", err)
	}
	return nil
}
