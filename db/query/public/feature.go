package public

import (
	"context"
	"fmt"

	"github.com/Gleipnir-Technology/jet/postgres"
	"github.com/Gleipnir-Technology/nidus-sync/db"
	"github.com/Gleipnir-Technology/nidus-sync/db/gen/nidus-sync/public/model"
	"github.com/Gleipnir-Technology/nidus-sync/db/gen/nidus-sync/public/table"
)

func FeaturesFromSiteID(ctx context.Context, txn db.Ex, site_id int64) ([]model.Feature, error) {
	statement := table.Feature.SELECT(
		table.Feature.AllColumns,
	).FROM(table.Feature).
		WHERE(table.Feature.SiteID.EQ(postgres.Int(site_id)))
	result, err := db.ExecuteManyTx[model.Feature](ctx, txn, statement)
	if err != nil {
		return []model.Feature{}, fmt.Errorf("query: %w", err)
	}
	return result, nil
}
func FeaturesFromSiteIDs(ctx context.Context, txn db.Ex, site_ids []int64) ([]model.Feature, error) {
	sql_ids := make([]postgres.Expression, len(site_ids))
	for i, site_id := range site_ids {
		sql_ids[i] = postgres.Int(site_id)
	}
	statement := table.Feature.SELECT(
		table.Feature.AllColumns,
	).FROM(table.Feature).
		WHERE(table.Feature.SiteID.IN(sql_ids...))
	result, err := db.ExecuteManyTx[model.Feature](ctx, txn, statement)
	if err != nil {
		return []model.Feature{}, fmt.Errorf("query: %w", err)
	}
	return result, nil
}
