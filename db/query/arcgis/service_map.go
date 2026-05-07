package arcgis

import (
	"context"

	"github.com/Gleipnir-Technology/bob"
	"github.com/Gleipnir-Technology/nidus-sync/db"
	"github.com/Gleipnir-Technology/nidus-sync/db/gen/nidus-sync/arcgis/model"
	"github.com/Gleipnir-Technology/nidus-sync/db/gen/nidus-sync/arcgis/table"
	"github.com/go-jet/jet/v2/postgres"
)

func ServiceMapFromID(ctx context.Context, id string) (model.ServiceMap, error) {
	statement := table.ServiceMap.SELECT(
		table.ServiceMap.AllColumns,
	).FROM(table.ServiceMap).
		WHERE(table.ServiceMap.ArcgisID.EQ(postgres.String(id)))
	return db.ExecuteOne[model.ServiceMap](ctx, statement)
}
func ServiceMapsFromAccountID(ctx context.Context, account_id string) ([]model.ServiceMap, error) {
	statement := table.ServiceMap.SELECT(
		table.ServiceMap.AllColumns,
	).FROM(table.ServiceMap).
		WHERE(table.ServiceMap.AccountID.EQ(postgres.String(account_id)))
	return db.ExecuteMany[model.ServiceMap](ctx, statement)
}
func ServiceMapInsert(ctx context.Context, txn bob.Tx, m *model.ServiceMap) error {
	statement := table.ServiceMap.INSERT(table.ServiceMap.MutableColumns).
		MODEL(m)
	return db.ExecuteNoneTxBob(ctx, txn, statement)
}
