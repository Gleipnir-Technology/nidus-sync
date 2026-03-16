package text

import (
	"context"
	"fmt"

	"github.com/Gleipnir-Technology/bob"
	"github.com/Gleipnir-Technology/nidus-sync/config"
	"github.com/Gleipnir-Technology/nidus-sync/db/enums"
	"github.com/Gleipnir-Technology/nidus-sync/platform/types"
)

func Job(ctx context.Context, txn bob.Executor, text_id int32) error {
	return sendTextComplete(ctx, txn, text_id)
}

func ReportSubscriptionConfirmationText(ctx context.Context, destination types.E164, report_id string) error {
	content := fmt.Sprintf("Thanks for submitting mosquito report %s. Text for any questions. We'll send you updates as we get them.", report_id)
	origin := enums.CommsTextoriginWebsiteAction
	err := sendTextBegin(ctx, *types.NewE164(&config.PhoneNumberReport), destination, content, origin, true, true)
	if err != nil {
		return fmt.Errorf("Failed to send initial confirmation: %w", err)
	}
	return err
}
