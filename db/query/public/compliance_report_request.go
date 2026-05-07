package public

import (
	"context"

	"github.com/Gleipnir-Technology/nidus-sync/db"
	"github.com/Gleipnir-Technology/nidus-sync/db/gen/nidus-sync/public/model"
	"github.com/Gleipnir-Technology/nidus-sync/db/gen/nidus-sync/public/table"
	//"github.com/go-jet/jet/v2/postgres"
)

func ComplianceReportRequestInsert(ctx context.Context, txn db.Ex, m model.ComplianceReportRequest) (model.ComplianceReportRequest, error) {
	statement := table.ComplianceReportRequest.INSERT(table.ComplianceReportRequest.MutableColumns).
		MODEL(m).
		RETURNING(table.ComplianceReportRequest.AllColumns)
	return db.ExecuteOne[model.ComplianceReportRequest](ctx, statement)
}
