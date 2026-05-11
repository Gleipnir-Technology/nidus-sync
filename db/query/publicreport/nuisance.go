package publicreport

import (
	"context"
	//"time"

	"github.com/Gleipnir-Technology/nidus-sync/db"
	"github.com/Gleipnir-Technology/nidus-sync/db/gen/nidus-sync/publicreport/model"
	"github.com/Gleipnir-Technology/nidus-sync/db/gen/nidus-sync/publicreport/table"
	//"github.com/Gleipnir-Technology/jet/postgres"
)

func NuisanceInsert(ctx context.Context, txn db.Ex, m model.Nuisance) (model.Nuisance, error) {
	statement := table.Nuisance.INSERT(table.Nuisance.AllColumns).
		MODEL(m).
		RETURNING(table.Nuisance.AllColumns)
	return db.ExecuteOneTx[model.Nuisance](ctx, txn, statement)
}
