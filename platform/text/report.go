package text

import (
	"context"
	"fmt"

	"github.com/Gleipnir-Technology/bob"
	"github.com/Gleipnir-Technology/bob/dialect/psql"
	"github.com/Gleipnir-Technology/bob/dialect/psql/sm"
	"github.com/Gleipnir-Technology/nidus-sync/db"
	"github.com/Gleipnir-Technology/nidus-sync/db/enums"
	//"github.com/Gleipnir-Technology/nidus-sync/db/models"
	"github.com/Gleipnir-Technology/nidus-sync/platform/types"
	//"github.com/rs/zerolog/log"
	"github.com/stephenafamo/scan"
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

type reportIDs struct {
	ID             int32  `db:"id"`
	PublicID       string `db:"public_id"`
	OrganizationID int32  `db:"organization_id"`
}

// Get the list of reports that are still open for a particular text message recipient
// 'still open' is not well-defined throughout the system, but for now we'll go with
// 'not reviewed in any way'.
func reportsForTextRecipient(ctx context.Context, txn bob.Executor, destination types.E164) ([]reportIDs, error) {
	rows, err := bob.All(ctx, db.PGInstance.BobDB, psql.Select(
		sm.Columns(
			"r.id",
			"r.public_id",
			"r.organization_id",
		),
		sm.From("comms.text_job").As("t"),
		sm.InnerJoin("publicreport.report").As("r").OnEQ(
			psql.Quote("t", "report_id"),
			psql.Quote("r", "id"),
		),
		sm.Where(psql.Quote("t", "report_id").IsNotNull()),
		sm.Where(psql.Quote("t", "destination").EQ(psql.Arg(destination.PhoneString()))),
		sm.Where(psql.Quote("r", "status").EQ(psql.Arg(enums.PublicreportReportstatustypeReported))),
	), scan.StructMapper[reportIDs]())
	if err != nil {
		return []reportIDs{}, fmt.Errorf("query reports: %w", err)
	}

	return rows, nil
}
