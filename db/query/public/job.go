package public

import (
	"context"

	"github.com/Gleipnir-Technology/nidus-sync/db"
	"github.com/Gleipnir-Technology/nidus-sync/db/gen/nidus-sync/public/model"
	"github.com/Gleipnir-Technology/nidus-sync/db/gen/nidus-sync/public/table"
	//"github.com/go-jet/jet/v2/postgres"
)

func JobInsert(ctx context.Context, txn db.Ex, m model.Job) (model.Job, error) {
	statement := table.Job.INSERT(table.Job.MutableColumns).
		MODEL(m).
		RETURNING(table.Job.AllColumns)
	return db.ExecuteOne[model.Job](ctx, statement)
}
