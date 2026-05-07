package public

import (
	"context"

	//"github.com/Gleipnir-Technology/bob"
	"github.com/Gleipnir-Technology/nidus-sync/db"
	//"github.com/Gleipnir-Technology/nidus-sync/db/gen/nidus-sync/public/enum"
	"github.com/Gleipnir-Technology/nidus-sync/db/gen/nidus-sync/public/model"
	"github.com/Gleipnir-Technology/nidus-sync/db/gen/nidus-sync/public/table"
	//"github.com/go-jet/jet/v2/postgres"
)

func SignalInsert(ctx context.Context, txn db.Tx, m model.Signal) (model.Signal, error) {
	statement := table.Signal.INSERT(table.Signal.MutableColumns).
		MODEL(m).
		RETURNING(table.Signal.AllColumns)
	return db.ExecuteOne[model.Signal](ctx, statement)
}
