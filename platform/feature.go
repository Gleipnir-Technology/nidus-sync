package platform

import (
	"context"
	"fmt"

	//"github.com/aarondl/opt/omit"
	//"github.com/aarondl/opt/omitnull"
	"github.com/Gleipnir-Technology/nidus-sync/db"
	modelpublic "github.com/Gleipnir-Technology/nidus-sync/db/gen/nidus-sync/public/model"
	querypublic "github.com/Gleipnir-Technology/nidus-sync/db/query/public"
	"github.com/Gleipnir-Technology/nidus-sync/platform/types"
	"github.com/rs/zerolog/log"
)

func FeaturesForSite(ctx context.Context, site_id int64) ([]types.Feature, error) {
	features, err := featuresBySiteID(ctx, []int64{site_id})
	if err != nil {
		return nil, fmt.Errorf("features by site ID: %w", err)
	}
	return features[int32(site_id)], nil
}

func featuresBySiteID(ctx context.Context, site_ids []int64) (map[int32][]types.Feature, error) {
	features, err := querypublic.FeaturesFromSiteIDs(ctx, db.PGInstance.PGXPool, site_ids)
	if err != nil {
		return nil, fmt.Errorf("query features: %w", err)
	}
	feature_ids := make([]int64, len(features))
	for i, feature := range features {
		feature_ids[i] = int64(feature.ID)
	}

	feature_pools, err := querypublic.FeaturePoolsFromFeatures(ctx, db.PGInstance.PGXPool, feature_ids)
	if err != nil {
		return nil, fmt.Errorf("query feature pools: %w", err)
	}
	feature_pools_by_feature_id := make(map[int32]modelpublic.FeaturePool, len(feature_pools))
	for _, feature_pool := range feature_pools {
		feature_pools_by_feature_id[feature_pool.FeatureID] = feature_pool
	}

	results := make(map[int32][]types.Feature, len(site_ids))
	for _, site_id := range site_ids {
		results[int32(site_id)] = make([]types.Feature, 0)
	}
	for _, row := range features {
		features, ok := results[row.SiteID]
		if !ok {
			return nil, fmt.Errorf("impossible")
		}
		/*
			feature_pools, ok := feature_pools_by_feature_id[row.ID]
			if !ok {
				return nil, fmt.Errorf("impossible 2")
			}
		*/
		if row.Location == nil {
			log.Warn().Int32("id", row.ID).Msg("nil location")
			continue
		}
		location, err := types.LocationFromGeom(*row.Location)
		if err != nil {
			return nil, fmt.Errorf("location from geom on %d: %w", row.SiteID, err)
		}
		features = append(features, types.Feature{
			ID:       row.ID,
			Location: location,
			SiteID:   row.SiteID,
			Type:     "pool",
		})
		results[row.SiteID] = features
	}
	return results, nil
}
