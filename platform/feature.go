package platform

import (
	"context"
	"fmt"

	//"github.com/aarondl/opt/omit"
	//"github.com/aarondl/opt/omitnull"
	"github.com/Gleipnir-Technology/bob"
	"github.com/Gleipnir-Technology/bob/dialect/psql"
	"github.com/Gleipnir-Technology/bob/dialect/psql/sm"
	"github.com/Gleipnir-Technology/nidus-sync/db"
	"github.com/Gleipnir-Technology/nidus-sync/db/models"
	"github.com/Gleipnir-Technology/nidus-sync/platform/types"
	"github.com/stephenafamo/scan"
)

func featuresBySiteID(ctx context.Context, site_ids []int32) (map[int32][]types.Feature, error) {
	rows, err := bob.All(ctx, db.PGInstance.BobDB, psql.Select(
		sm.Columns(
			"feature.id AS id",
			"feature.site_id AS site_id",
			"COALESCE(ST_X(feature.location), 0) AS \"location.longitude\"",
			"COALESCE(ST_Y(feature.location), 0) AS \"location.latitude\"",
			"'pool' AS type",
		),
		sm.From("feature"),
		sm.InnerJoin("feature_pool").OnEQ(
			psql.Quote("feature", "id"),
			psql.Quote("feature_pool", "feature_id"),
		),
		sm.Where(
			models.Features.Columns.ID.EQ(psql.Any(site_ids)),
		),
	), scan.StructMapper[types.Feature]())
	if err != nil {
		return nil, fmt.Errorf("query features: %w", err)
	}
	results := make(map[int32][]types.Feature, len(site_ids))
	for _, site_id := range site_ids {
		results[site_id] = make([]types.Feature, 0)
	}
	for _, row := range rows {
		features, ok := results[row.SiteID]
		if !ok {
			return nil, fmt.Errorf("impossible")
		}
		features = append(features, types.Feature{
			ID:   row.ID,
			Type: "pool",
		})
		results[row.SiteID] = features
	}
	return results, nil
}
