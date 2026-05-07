package publicreport

import (
	"context"
	//"time"

	"github.com/Gleipnir-Technology/nidus-sync/db"
	"github.com/Gleipnir-Technology/nidus-sync/db/gen/nidus-sync/publicreport/model"
	"github.com/Gleipnir-Technology/nidus-sync/db/gen/nidus-sync/publicreport/table"
	//"github.com/go-jet/jet/v2/postgres"
)

func WaterInsert(ctx context.Context, txn db.Ex, m model.Water) (model.Water, error) {
	statement := table.Water.INSERT(table.Water.AllColumns).
		MODEL(m)
	return db.ExecuteOneTx[model.Water](ctx, txn, statement)
}
