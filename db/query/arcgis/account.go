package arcgis

import (
	"context"

	"github.com/Gleipnir-Technology/bob"
	"github.com/Gleipnir-Technology/nidus-sync/db"
	"github.com/Gleipnir-Technology/nidus-sync/db/gen/nidus-sync/arcgis/model"
	"github.com/Gleipnir-Technology/nidus-sync/db/gen/nidus-sync/arcgis/table"
	"github.com/go-jet/jet/v2/postgres"
)

func AccountFromID(ctx context.Context, org_id string) (*model.Account, error) {
	statement := table.Account.SELECT(
		table.Account.AllColumns,
	).FROM(table.Account).
		WHERE(table.Account.ID.EQ(postgres.String(org_id)))
	return db.ExecuteOne[model.Account](ctx, statement)
}
func AccountInsert(ctx context.Context, txn bob.Tx, m *model.Account) (*model.Account, error) {
	statement := table.Account.INSERT(table.Account.AllColumns).
		MODEL(m)
	return db.ExecuteOneTx[model.Account](ctx, txn, statement)
}
