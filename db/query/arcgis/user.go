package arcgis

import (
	"context"

	"github.com/Gleipnir-Technology/bob"
	"github.com/Gleipnir-Technology/nidus-sync/db"
	"github.com/Gleipnir-Technology/nidus-sync/db/gen/nidus-sync/arcgis/model"
	"github.com/Gleipnir-Technology/nidus-sync/db/gen/nidus-sync/arcgis/table"
	"github.com/go-jet/jet/v2/postgres"
)

func UserFromID(ctx context.Context, id string) (model.User, error) {
	statement := table.User.SELECT(table.User.AllColumns).
		FROM(table.User).
		WHERE(table.User.ID.EQ(postgres.String(id)))
	return db.ExecuteOne[model.User](ctx, statement)
}
func UserInsert(ctx context.Context, txn bob.Tx, m *model.User) (model.User, error) {
	statement := table.User.INSERT(table.User.MutableColumns).
		MODEL(m).
		RETURNING(table.User.AllColumns)
	return db.ExecuteOneTxBob[model.User](ctx, txn, statement)
}
