package publicreport

import (
	"context"
	//"time"

	//"github.com/Gleipnir-Technology/bob"
	"github.com/Gleipnir-Technology/nidus-sync/db"
	"github.com/Gleipnir-Technology/nidus-sync/db/gen/nidus-sync/publicreport/model"
	"github.com/Gleipnir-Technology/nidus-sync/db/gen/nidus-sync/publicreport/table"
	"github.com/go-jet/jet/v2/postgres"
)

type ComplianceUpdater = db.Updater[table.ComplianceTable, model.Compliance]

func NewUpdaterCompliance() db.Updater[table.ComplianceTable, model.Compliance] {
	return db.NewUpdater[table.ComplianceTable, model.Compliance](
		*table.Compliance,
		table.Compliance.ReportID,
	)

}
func ComplianceFromID(ctx context.Context, txn db.Tx, report_id int64) (model.Compliance, error) {
	statement := table.Report.SELECT(
		table.Compliance.AllColumns,
	).FROM(table.Compliance).
		WHERE(table.Compliance.ReportID.EQ(postgres.Int(report_id)))
	return db.ExecuteOneTx[model.Compliance](ctx, txn, statement)
}
func ComplianceInsert(ctx context.Context, txn db.Ex, m model.Compliance) (model.Compliance, error) {
	statement := table.Compliance.INSERT(table.Compliance.MutableColumns).
		MODEL(m).
		RETURNING(table.Compliance.AllColumns)
	return db.ExecuteOneTx[model.Compliance](ctx, txn, statement)
}
