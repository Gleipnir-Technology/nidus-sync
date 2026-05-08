package publicreport

import (
	"context"
	//"time"

	"github.com/Gleipnir-Technology/nidus-sync/db"
	"github.com/Gleipnir-Technology/nidus-sync/db/gen/nidus-sync/publicreport/model"
	"github.com/Gleipnir-Technology/nidus-sync/db/gen/nidus-sync/publicreport/table"
	//"github.com/go-jet/jet/v2/postgres"
)

func NuisanceInsert(ctx context.Context, txn db.Ex, m model.Nuisance) (model.Nuisance, error) {
	statement := table.Nuisance.INSERT(table.Nuisance.MutableColumns).
		MODEL(m).
		RETURNING(table.Nuisance.AllColumns)
	return db.ExecuteOneTx[model.Nuisance](ctx, txn, statement)
}
