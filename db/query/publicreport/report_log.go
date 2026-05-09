package publicreport

import (
	"context"

	"github.com/Gleipnir-Technology/nidus-sync/db"
	"github.com/Gleipnir-Technology/nidus-sync/db/gen/nidus-sync/publicreport/model"
	"github.com/Gleipnir-Technology/nidus-sync/db/gen/nidus-sync/publicreport/table"
	//"github.com/Gleipnir-Technology/jet/postgres"
)

func ReportLogInsert(ctx context.Context, txn db.Ex, m model.ReportLog) (model.ReportLog, error) {
	statement := table.ReportLog.INSERT(table.ReportLog.MutableColumns).
		MODEL(m).
		RETURNING(table.ReportLog.AllColumns)
	return db.ExecuteOneTx[model.ReportLog](ctx, txn, statement)
}
