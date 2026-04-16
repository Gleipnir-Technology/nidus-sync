package platform

import (
	"context"
	"fmt"
	"time"

	"github.com/Gleipnir-Technology/bob"
	"github.com/Gleipnir-Technology/nidus-sync/db"
	"github.com/Gleipnir-Technology/nidus-sync/db/enums"
	"github.com/Gleipnir-Technology/nidus-sync/db/models"
	"github.com/Gleipnir-Technology/nidus-sync/platform/types"
	"github.com/aarondl/opt/omit"
	"github.com/aarondl/opt/omitnull"
	//"github.com/rs/zerolog/log"
)

// Create a lead from the given signal and site
func LeadCreate(ctx context.Context, user User, signal_id int32, site_id int32, pool_location *types.Location) (*int32, error) {
	txn, err := db.PGInstance.BobDB.BeginTx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("start transaction: %w", err)
	}
	defer txn.Rollback(ctx)

	lead_id, err := leadCreate(ctx, txn, user, signal_id, site_id, pool_location)
	if err != nil {
		return nil, fmt.Errorf("inner leadcreate: %w", err)
	}
	txn.Commit(ctx)
	return lead_id, nil
}

func leadCreate(ctx context.Context, txn bob.Executor, user User, signal_id int32, site_id int32, pool_location *types.Location) (*int32, error) {
	lead, err := models.Leads.Insert(&models.LeadSetter{
		Created: omit.From(time.Now()),
		Creator: omit.From(int32(user.ID)),
		// ID
		OrganizationID: omit.From(int32(user.Organization.ID)),
		SiteID:         omitnull.From(site_id),
		Type:           omit.From(enums.LeadtypeGreenPool),
	}).One(ctx, txn)
	if err != nil {
		return nil, fmt.Errorf("failed to create lead: %w", err)
	}
	return &lead.ID, nil
}
