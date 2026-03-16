package text

import (
	"context"
	"fmt"

	"github.com/Gleipnir-Technology/nidus-sync/db"
	"github.com/Gleipnir-Technology/nidus-sync/db/enums"
	"github.com/Gleipnir-Technology/nidus-sync/platform/types"
	//"github.com/rs/zerolog/log"
)

func sendReportSubscription(ctx context.Context, source, destination types.E164, content string) error {
	err := EnsureInDB(ctx, db.PGInstance.BobDB, destination)
	if err != nil {
		return fmt.Errorf("Failed to ensure text message destination is in the DB: %w", err)
	}
	status, err := phoneStatus(ctx, destination)
	if err != nil {
		return fmt.Errorf("Failed to check if subscribed: %w", err)
	}
	switch status {
	case enums.CommsPhonestatustypeUnconfirmed:
		err = delayMessage(ctx, enums.CommsTextjobsourceRmo, destination, content, enums.CommsTextjobtypeReportConfirmation)
		if err != nil {
			return fmt.Errorf("Failed to delay report subscription message: %w", err)
		}
		err := ensureInitialText(ctx, source, destination)
		if err != nil {
			return fmt.Errorf("Failed to ensure initial text has been sent: %w", err)
		}
		return nil
	case enums.CommsPhonestatustypeOkToSend:
		err = sendTextBegin(ctx, source, destination, content, enums.CommsTextoriginWebsiteAction, false, true)
		if err != nil {
			return fmt.Errorf("Failed to send report subscription confirmation: %w", err)
		}
	case enums.CommsPhonestatustypeStopped:
		resendInitialText(ctx, source, destination)
	}
	return nil
}
