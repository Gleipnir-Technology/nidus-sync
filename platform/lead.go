package platform

import (
	"context"
	"fmt"
	"time"

	"github.com/Gleipnir-Technology/bob/dialect/psql"
	"github.com/Gleipnir-Technology/bob/dialect/psql/um"
	"github.com/Gleipnir-Technology/nidus-sync/db"
	"github.com/Gleipnir-Technology/nidus-sync/db/enums"
	"github.com/Gleipnir-Technology/nidus-sync/db/models"
	"github.com/Gleipnir-Technology/nidus-sync/platform/geom"
	"github.com/Gleipnir-Technology/nidus-sync/platform/types"
	"github.com/aarondl/opt/omit"
	"github.com/aarondl/opt/omitnull"
	"github.com/rs/zerolog/log"
)

// Create a lead from the given signal and site
func LeadCreate(ctx context.Context, user User, signal_id int32, site_id int32, pool_location *types.Location) (*int32, error) {
	txn, err := db.PGInstance.BobDB.BeginTx(ctx, nil)
	defer txn.Rollback(ctx)
	if err != nil {
		return nil, fmt.Errorf("start transaction: %w", err)
	}

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
	_, err = psql.Update(
		um.Table("signal"),
		um.SetCol("addressed").ToArg(time.Now()),
		um.SetCol("addressor").ToArg(user.ID),
		um.Where(psql.Quote("id").EQ(psql.Arg(signal_id))),
	).Exec(ctx, txn)
	if err != nil {
		return nil, fmt.Errorf("failed to update signal %d: %w", signal_id, err)
	}
	if pool_location != nil {
		log.Info().Float64("lat", pool_location.Latitude).Float64("lng", pool_location.Longitude).Msg("got pool location")
		geom_query := geom.PostgisPointQuery(*pool_location)
		_, err = psql.Update(
			um.Table("pool"),
			um.SetCol("geometry").To(geom_query),
			um.From("signal_pool"),
			um.Where(psql.Quote("signal_pool", "pool_id").EQ(psql.Quote("pool", "id"))),
			um.Where(psql.Quote("signal_pool", "signal_id").EQ(psql.Arg(signal_id))),
		).Exec(ctx, txn)
		if err != nil {
			return nil, fmt.Errorf("failed to update pool through signal %d: %w", signal_id, err)
		}
	}
	txn.Commit(ctx)
	return &lead.ID, nil
}
