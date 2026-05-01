package public

import (
	"context"
	//"time"

	//"github.com/Gleipnir-Technology/bob"
	"github.com/Gleipnir-Technology/nidus-sync/db"
	"github.com/Gleipnir-Technology/nidus-sync/db/gen/nidus-sync/publicreport/model"
	"github.com/Gleipnir-Technology/nidus-sync/db/gen/nidus-sync/publicreport/table"
	"github.com/go-jet/jet/v2/postgres"
)

/*
	func CommunicationInsert(ctx context.Context, txn bob.Tx, m *model.Communication) (*model.Communication, error) {
		m.Created = time.Now()
		statement := table.Communication.INSERT(table.Communication.MutableColumns).
			MODEL(m)
		return db.ExecuteOne[model.Communication](ctx, statement)
	}
*/
func PublicReportsFromIDs(ctx context.Context, report_ids []int64) ([]*model.Report, error) {
	sql_ids := make([]postgres.Expression, len(report_ids))
	for i, report_id := range report_ids {
		sql_ids[i] = postgres.Int(report_id)
	}
	statement := table.Report.SELECT(
		table.Report.AllColumns,
	).FROM(table.Report).
		WHERE(table.Report.ID.IN(sql_ids...))
	return db.ExecuteMany[model.Report](ctx, statement)
}
