package platform

import (
	"context"
	"fmt"
	"time"

	"github.com/Gleipnir-Technology/bob"
	"github.com/Gleipnir-Technology/bob/dialect/psql"
	"github.com/Gleipnir-Technology/bob/dialect/psql/sm"
	"github.com/Gleipnir-Technology/nidus-sync/db"
	"github.com/Gleipnir-Technology/nidus-sync/db/gen/nidus-sync/public/model"
	"github.com/Gleipnir-Technology/nidus-sync/db/models"
	query "github.com/Gleipnir-Technology/nidus-sync/db/query/public"
	"github.com/Gleipnir-Technology/nidus-sync/platform/types"
	//"github.com/rs/zerolog/log"
	"github.com/stephenafamo/scan"
)

// Create a lead from the given signal and site
func LeadCreate(ctx context.Context, user User, signal_id int32, site_id int32, pool_location *types.Location) (model.Lead, error) {
	txn, err := db.BeginTxn(ctx)
	if err != nil {
		return model.Lead{}, fmt.Errorf("start transaction: %w", err)
	}
	defer txn.Rollback(ctx)

	lead, err := leadCreate(ctx, txn, user, signal_id, site_id, pool_location)
	if err != nil {
		return model.Lead{}, fmt.Errorf("inner leadcreate: %w", err)
	}
	txn.Commit(ctx)
	return lead, nil
}

func leadCreate(ctx context.Context, txn db.Ex, user User, signal_id int32, site_id int32, pool_location *types.Location) (model.Lead, error) {
	lead := model.Lead{
		Created: time.Now(),
		Creator: int32(user.ID),
		// ID
		OrganizationID: int32(user.Organization.ID),
		SiteID:         &site_id,
		Type:           model.Leadtype_GreenPool,
	}
	lead, err := query.LeadInsert(ctx, txn, lead)
	if err != nil {
		return model.Lead{}, fmt.Errorf("failed to create lead: %w", err)
	}
	return lead, nil
}
func leadsBySiteID(ctx context.Context, site_ids []int64) (map[int32][]*types.Lead, error) {
	rows, err := bob.All(ctx, db.PGInstance.BobDB, psql.Select(
		sm.Columns(
			models.Leads.Columns.ID.As("id"),
			models.Leads.Columns.SiteID.As("site_id"),
			models.Leads.Columns.Type.As("type"),
		),
		sm.From(models.Leads.Name()),
		sm.Where(
			models.Leads.Columns.SiteID.EQ(psql.Any(site_ids)),
		),
	), scan.StructMapper[*types.Lead]())
	if err != nil {
		return nil, fmt.Errorf("query leads: %w", err)
	}
	lead_ids := make([]int32, len(rows))
	for i, row := range rows {
		lead_ids[i] = row.ID
	}
	compliance_report_requests, err := ComplianceReportRequestByLeadID(ctx, lead_ids)
	for _, row := range rows {
		crrs, ok := compliance_report_requests[row.ID]
		if !ok {
			return nil, fmt.Errorf("impossible")
		}
		row.ComplianceReportRequests = crrs
	}
	results := make(map[int32][]*types.Lead, len(site_ids))
	for _, site_id := range site_ids {
		results[int32(site_id)] = make([]*types.Lead, 0)
	}
	for _, row := range rows {
		leads, ok := results[row.SiteID]
		if !ok {
			return nil, fmt.Errorf("impossible")
		}
		leads = append(leads, row)
		results[row.SiteID] = leads
	}
	return results, nil
}
