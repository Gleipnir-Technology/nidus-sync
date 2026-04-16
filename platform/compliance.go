package platform

import (
	"context"
	"fmt"
	"time"

	"github.com/Gleipnir-Technology/bob/dialect/psql"
	"github.com/Gleipnir-Technology/bob/dialect/psql/sm"
	"github.com/Gleipnir-Technology/nidus-sync/db"
	"github.com/Gleipnir-Technology/nidus-sync/db/models"
	"github.com/Gleipnir-Technology/nidus-sync/platform/types"
	"github.com/aarondl/opt/omit"
	"github.com/aarondl/opt/omitnull"
)

func ComplianceRequestMailerCreate(ctx context.Context, user User, site_id int32) (int32, error) {
	txn, err := db.PGInstance.BobDB.BeginTx(ctx, nil)
	if err != nil {
		return 0, fmt.Errorf("start txn: %w", err)
	}
	defer txn.Rollback(ctx)
	site, err := models.FindSite(ctx, txn, site_id)
	if err != nil {
		return 0, fmt.Errorf("find site: %w", err)
	}
	if site.OrganizationID != user.Organization.ID {
		return 0, fmt.Errorf("permission denied")
	}
	features, err := models.Features.Query(
		models.SelectWhere.Features.SiteID.EQ(site.ID),
	).All(ctx, txn)
	if err != nil {
		return 0, fmt.Errorf("find features: %w", err)
	}
	feature_ids := make([]int32, len(features))
	for i, f := range features {
		feature_ids[i] = f.ID
	}
	feature_pools, err := models.FeaturePools.Query(
		sm.Where(
			models.FeaturePools.Columns.FeatureID.EQ(psql.Any(feature_ids)),
		),
	).All(ctx, txn)
	if err != nil {
		return 0, fmt.Errorf("find feature pools: %w", err)
	}
	if len(feature_pools) != 1 {
		return 0, fmt.Errorf("wrong number of pools: %d", len(feature_pools))
	}
	feature_pool := feature_pools[0]
	var feature *models.Feature
	for _, f := range features {
		if f.ID == feature_pool.FeatureID {
			feature = f
		}
	}
	if feature == nil {
		return 0, fmt.Errorf("match feature %d", feature_pool.FeatureID)
	}
	location := types.Location{
		Latitude:  feature.LocationLatitude.GetOr(0),
		Longitude: feature.LocationLongitude.GetOr(0),
	}
	signal, err := SignalCreateFromPool(ctx, txn, user, site.ID, feature_pool.FeatureID, location)
	if err != nil {
		return 0, fmt.Errorf("create signal from ppol: %w", err)
	}
	lead_id, err := leadCreate(ctx, txn, user, *signal, site.ID, &location)
	if err != nil {
		return 0, fmt.Errorf("create lead from ppol: %w", err)
	}
	public_id, err := GenerateReportID()
	if err != nil {
		return 0, fmt.Errorf("create public id: %w", err)
	}
	setter := models.ComplianceReportRequestSetter{
		Created: omit.From(time.Now()),
		Creator: omit.From(int32(user.ID)),
		// ID
		PublicID: omit.From(public_id),
		LeadID:   omitnull.From(*lead_id),
	}
	req, err := models.ComplianceReportRequests.Insert(&setter).One(ctx, txn)
	if err != nil {
		return 0, fmt.Errorf("create compliance report request: %w", err)
	}
	return req.ID, nil
}
