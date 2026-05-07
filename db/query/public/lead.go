package public

import (
	"context"

	"github.com/Gleipnir-Technology/nidus-sync/db"
	"github.com/Gleipnir-Technology/nidus-sync/db/gen/nidus-sync/public/model"
	"github.com/Gleipnir-Technology/nidus-sync/db/gen/nidus-sync/public/table"
	//"github.com/go-jet/jet/v2/postgres"
)

func LeadInsert(ctx context.Context, txn db.Ex, m model.Lead) (model.Lead, error) {
	statement := table.Lead.INSERT(table.Lead.MutableColumns).
		MODEL(m).
		RETURNING(table.Lead.AllColumns)
	return db.ExecuteOne[model.Lead](ctx, statement)
}
