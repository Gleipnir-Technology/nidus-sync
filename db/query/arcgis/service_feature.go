package arcgis

import (
	"context"

	"github.com/Gleipnir-Technology/bob"
	"github.com/Gleipnir-Technology/nidus-sync/db"
	"github.com/Gleipnir-Technology/nidus-sync/db/gen/nidus-sync/arcgis/model"
	"github.com/Gleipnir-Technology/nidus-sync/db/gen/nidus-sync/arcgis/table"
	"github.com/go-jet/jet/v2/postgres"
)

func ServiceFeatureFromID(ctx context.Context, id string) (model.ServiceFeature, error) {
	statement := table.ServiceFeature.SELECT(
		table.ServiceFeature.AllColumns,
	).FROM(table.ServiceFeature).
		WHERE(table.ServiceFeature.ItemID.EQ(postgres.String(id)))
	return db.ExecuteOne[model.ServiceFeature](ctx, statement)
}
func ServiceFeatureFromURL(ctx context.Context, url string) (model.ServiceFeature, error) {
	statement := table.ServiceFeature.SELECT(
		table.ServiceFeature.AllColumns,
	).FROM(table.ServiceFeature).
		WHERE(table.ServiceFeature.URL.EQ(postgres.String(url)))
	return db.ExecuteOne[model.ServiceFeature](ctx, statement)
}
func ServiceFeatureInsert(ctx context.Context, txn bob.Tx, m model.ServiceFeature) error {
	statement := table.ServiceMap.INSERT(table.ServiceMap.MutableColumns).
		MODEL(m).
		RETURNING(table.ServiceFeature.AllColumns)
	return db.ExecuteNoneTxBob(ctx, txn, statement)
}
