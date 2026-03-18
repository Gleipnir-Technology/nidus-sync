package text

import (
	"context"
	"fmt"

	"github.com/Gleipnir-Technology/bob"
	"github.com/Gleipnir-Technology/nidus-sync/db/enums"
	"github.com/Gleipnir-Technology/nidus-sync/db/models"
	"github.com/Gleipnir-Technology/nidus-sync/platform/types"
	//"github.com/rs/zerolog/log"
)

// Send a message from a district to a public reporter within the context of the public report
func ReportMessage(ctx context.Context, txn bob.Executor, user_id int32, report_id int32, destination types.E164, content string) (*int32, error) {
	job_id, err := sendTextBegin(ctx, txn, &user_id, &report_id, destination, content, enums.CommsTextjobtypeReportMessage)
	if err != nil {
		return nil, fmt.Errorf("Failed to send initial confirmation: %w", err)
	}
	return job_id, nil
}

// Send a message from the system to a public reporter indicating they are subscribed to updates on the report
func ReportSubscriptionConfirmationText(ctx context.Context, txn bob.Executor, destination types.E164, report_id string) error {
	content := fmt.Sprintf("Thanks for submitting mosquito report %s. Text for any questions. We'll send you updates as we get them.", report_id)
	_, err := sendTextBegin(ctx, txn, nil, nil, destination, content, enums.CommsTextjobtypeReportConfirmation)
	if err != nil {
		return fmt.Errorf("Failed to send initial confirmation: %w", err)
	}
	return err
}
func reportForTextRecipient(ctx context.Context, txn bob.Executor, destination types.E164) (*models.PublicreportReport, error) {
	/*return models.ReportText
	psql.Query(
	return Addresses.Query(
		sm.Where(Addresses.Columns.ID.EQ(psql.Arg(IDPK))),
	).Exists(ctx, exec)
	*/
	return nil, nil
}
