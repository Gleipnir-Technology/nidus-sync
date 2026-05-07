package publicreport

import (
	"context"

	"github.com/Gleipnir-Technology/nidus-sync/db"
	"github.com/Gleipnir-Technology/nidus-sync/db/gen/nidus-sync/publicreport/model"
	"github.com/Gleipnir-Technology/nidus-sync/db/gen/nidus-sync/publicreport/table"
	//"github.com/go-jet/jet/v2/postgres"
)

func ImageInsert(ctx context.Context, txn db.Ex, m model.Image) (model.Image, error) {
	statement := table.Image.INSERT(table.Image.AllColumns).
		MODEL(m).
		RETURNING(table.Image.AllColumns)
	return db.ExecuteOneTx[model.Image](ctx, txn, statement)
}
