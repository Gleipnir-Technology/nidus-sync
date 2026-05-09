package public

import (
	"context"
	"fmt"

	"github.com/Gleipnir-Technology/jet/postgres"
	"github.com/Gleipnir-Technology/nidus-sync/db"
	"github.com/Gleipnir-Technology/nidus-sync/db/gen/nidus-sync/public/model"
	"github.com/Gleipnir-Technology/nidus-sync/db/gen/nidus-sync/public/table"
)

func FeaturePoolsFromFeatures(ctx context.Context, txn db.Ex, feature_ids []int64) ([]model.FeaturePool, error) {
	sql_ids := make([]postgres.Expression, len(feature_ids))
	for i, site_id := range feature_ids {
		sql_ids[i] = postgres.Int(site_id)
	}
	statement := table.FeaturePool.SELECT(
		table.FeaturePool.AllColumns,
	).FROM(table.FeaturePool).
		WHERE(table.FeaturePool.FeatureID.IN(sql_ids...))
	result, err := db.ExecuteManyTx[model.FeaturePool](ctx, txn, statement)
	if err != nil {
		return []model.FeaturePool{}, fmt.Errorf("query: %w", err)
	}
	return result, nil
}
