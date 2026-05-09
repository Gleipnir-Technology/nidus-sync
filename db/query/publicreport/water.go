package publicreport

import (
	"context"
	//"time"

	"github.com/Gleipnir-Technology/nidus-sync/db"
	"github.com/Gleipnir-Technology/nidus-sync/db/gen/nidus-sync/publicreport/model"
	"github.com/Gleipnir-Technology/nidus-sync/db/gen/nidus-sync/publicreport/table"
	//"github.com/Gleipnir-Technology/jet/postgres"
)

func WaterInsert(ctx context.Context, txn db.Ex, m model.Water) (model.Water, error) {
	statement := table.Water.INSERT(table.Water.MutableColumns).
		MODEL(m).
		RETURNING(table.Water.AllColumns)
	return db.ExecuteOneTx[model.Water](ctx, txn, statement)
}
