package platform

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/Gleipnir-Technology/bob"
	"github.com/Gleipnir-Technology/bob/dialect/psql"
	"github.com/Gleipnir-Technology/bob/dialect/psql/sm"
	"github.com/Gleipnir-Technology/nidus-sync/db"
	modelpublic "github.com/Gleipnir-Technology/nidus-sync/db/gen/nidus-sync/public/model"
	"github.com/Gleipnir-Technology/nidus-sync/db/models"
	querypublic "github.com/Gleipnir-Technology/nidus-sync/db/query/public"
	"github.com/Gleipnir-Technology/nidus-sync/platform/background"
	"github.com/Gleipnir-Technology/nidus-sync/platform/event"
	"github.com/Gleipnir-Technology/nidus-sync/platform/types"
	//"github.com/rs/zerolog/log"
	"github.com/stephenafamo/scan"
)

func ComplianceRequestMailerCreate(ctx context.Context, user User, site_id int64) (int32, error) {
	txn, err := db.BeginTxn(ctx)
	if err != nil {
		return 0, fmt.Errorf("start txn: %w", err)
	}
	defer txn.Rollback(ctx)
	site, err := querypublic.SiteFromIDForOrg(ctx, txn, site_id, int64(user.Organization.ID))
	if err != nil {
		return 0, fmt.Errorf("find site: %w", err)
	}
	if site.OrganizationID != user.Organization.ID {
		return 0, fmt.Errorf("permission denied")
	}
	address, err := querypublic.AddressFromID(ctx, txn, int64(site.AddressID))
	if err != nil {
		return 0, fmt.Errorf("find address %d: %w", site.AddressID, err)
	}
	if address.PostalCode == "" {
		return 0, fmt.Errorf("address %d does not have a postal code", address.ID)
	}
	features, err := querypublic.FeaturesFromSiteID(ctx, txn, int64(site.ID))
	if err != nil {
		return 0, fmt.Errorf("find features: %w", err)
	}
	feature_ids := make([]int64, len(features))
	for i, f := range features {
		feature_ids[i] = int64(f.ID)
	}
	feature_pools, err := querypublic.FeaturePoolsFromFeatures(ctx, txn, feature_ids)
	if err != nil {
		return 0, fmt.Errorf("find feature pools: %w", err)
	}
	if len(feature_pools) != 1 {
		return 0, fmt.Errorf("wrong number of pools: %d", len(feature_pools))
	}
	feature_pool := feature_pools[0]
	var feature *modelpublic.Feature
	for _, f := range features {
		if f.ID == feature_pool.FeatureID {
			feature = &f
		}
	}
	if feature == nil {
		return 0, fmt.Errorf("match feature %d", feature_pool.FeatureID)
	}
	if feature.Location == nil {
		return 0, fmt.Errorf("nil location %d", feature.ID)
	}
	location, err := types.LocationFromGeom(*feature.Location)
	if err != nil {
		return 0, fmt.Errorf("location from geom: %w", err)
	}
	signal, err := SignalCreateFromPool(ctx, txn, user, site.ID, feature_pool.FeatureID, location)
	if err != nil {
		return 0, fmt.Errorf("create signal from ppol: %w", err)
	}
	lead, err := leadCreate(ctx, txn, user, signal.ID, site.ID, &location)
	if err != nil {
		return 0, fmt.Errorf("create lead from ppol: %w", err)
	}
	public_id, err := GenerateReportID()
	if err != nil {
		return 0, fmt.Errorf("create public id: %w", err)
	}
	setter := modelpublic.ComplianceReportRequest{
		Created: time.Now(),
		Creator: int32(user.ID),
		// ID
		PublicID: public_id,
		LeadID:   &lead.ID,
	}
	req, err := querypublic.ComplianceReportRequestInsert(ctx, txn, setter)
	if err != nil {
		return 0, fmt.Errorf("create compliance report request: %w", err)
	}
	err = background.NewComplianceMailer(ctx, txn, req.ID)
	if err != nil {
		return 0, fmt.Errorf("create background compliance mailer job: %w", err)
	}
	event.Updated(event.TypeSite, user.Organization.ID, strconv.Itoa(int(site.ID)))
	txn.Commit(ctx)

	return req.ID, nil
}

func ComplianceReportRequestByLeadID(ctx context.Context, lead_ids []int32) (map[int32][]*types.ComplianceReportRequest, error) {
	rows, err := models.ComplianceReportRequests.Query(
		sm.Where(models.ComplianceReportRequests.Columns.LeadID.EQ(psql.Any(lead_ids))),
	).All(ctx, db.PGInstance.BobDB)
	if err != nil {
		return nil, fmt.Errorf("query reports: %w", err)
	}
	results := make(map[int32][]*types.ComplianceReportRequest, len(lead_ids))
	for _, lead_id := range lead_ids {
		results[lead_id] = make([]*types.ComplianceReportRequest, 0)
	}
	for _, row := range rows {
		lead_id := row.LeadID.MustGet()
		crrs, ok := results[lead_id]
		if !ok {
			return nil, fmt.Errorf("impossible")
		}
		crrs = append(crrs, types.ComplianceReportRequestFromModel(row))
		results[lead_id] = crrs
	}
	return results, nil
}
func ComplianceReportRequestFromPublicID(ctx context.Context, public_id string) (*types.ComplianceReportRequest, error) {
	row, err := models.ComplianceReportRequests.Query(
		sm.Where(models.ComplianceReportRequests.Columns.PublicID.EQ(psql.Arg(public_id))),
	).One(ctx, db.PGInstance.BobDB)
	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			return nil, nil
		}
		return nil, fmt.Errorf("query CRR: %w", err)
	}
	return types.ComplianceReportRequestFromModel(row), nil
}
func OrganizationIDForComplianceReportRequest(ctx context.Context, public_id string) (int32, error) {
	type _Row struct {
		ID int32 `db:"organization_id"`
	}
	row, err := bob.One(ctx, db.PGInstance.BobDB, psql.Select(
		sm.Columns(
			models.Sites.Columns.OrganizationID,
		),
		sm.From(models.ComplianceReportRequests.NameAs()),
		sm.InnerJoin(models.Leads.NameAs()).On(
			models.ComplianceReportRequests.Columns.LeadID.EQ(models.Leads.Columns.ID)),
		sm.InnerJoin(models.Sites.NameAs()).On(
			models.Leads.Columns.SiteID.EQ(models.Sites.Columns.ID)),
		sm.Where(models.ComplianceReportRequests.Columns.PublicID.EQ(psql.Arg(public_id))),
	), scan.StructMapper[_Row]())
	if err != nil {
		return 0, fmt.Errorf("query compliance report request")
	}
	return row.ID, nil
}
