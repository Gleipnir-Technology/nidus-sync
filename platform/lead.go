package platform

import (
	"context"
	"fmt"
	"time"

	"github.com/Gleipnir-Technology/bob"
	"github.com/Gleipnir-Technology/bob/dialect/psql"
	"github.com/Gleipnir-Technology/bob/dialect/psql/um"
	"github.com/Gleipnir-Technology/nidus-sync/db"
	"github.com/Gleipnir-Technology/nidus-sync/db/enums"
	"github.com/Gleipnir-Technology/nidus-sync/db/models"
	"github.com/Gleipnir-Technology/nidus-sync/platform/geocode"
	"github.com/Gleipnir-Technology/nidus-sync/platform/geom"
	"github.com/aarondl/opt/omit"
	"github.com/aarondl/opt/omitnull"
	"github.com/rs/zerolog/log"
)

// Create a lead from the given signal and site
func LeadCreate(ctx context.Context, user User, signal_id int32, site_id int32, pool_location *Location) (*int32, error) {
	txn, err := db.PGInstance.BobDB.BeginTx(ctx, nil)
	defer txn.Rollback(ctx)
	if err != nil {
		return nil, fmt.Errorf("start transaction: %w", err)
	}

	lead, err := models.Leads.Insert(&models.LeadSetter{
		Created: omit.From(time.Now()),
		Creator: omit.From(int32(user.ID)),
		// ID
		OrganizationID: omit.From(int32(user.Organization.ID())),
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

func siteFromAddress(ctx context.Context, txn bob.Tx, user User, address_id int32) (*models.Site, error) {
	site, err := models.Sites.Query(
		models.SelectWhere.Sites.AddressID.EQ(address_id),
		models.SelectWhere.Sites.OrganizationID.EQ(user.Organization.ID()),
	).One(ctx, txn)
	if err == nil {
		return site, nil
	}
	if err.Error() != "sql: no rows in result set" {
		return nil, fmt.Errorf("query site: %w", err)
	}
	return SiteCreate(ctx, txn, user, address_id)
}
func siteFromAddressRaw(ctx context.Context, txn bob.Tx, user User, address string) (*models.Site, error) {
	// Geocode
	geo, err := geocode.GeocodeRaw(ctx, user.Organization.model, address)
	if err != nil {
		return nil, fmt.Errorf("geocode: %w", err)
	}
	a, err := geocode.EnsureAddress(ctx, txn, geo.Address, geo.Location)
	if err != nil {
		return nil, fmt.Errorf("ensure address: %w", err)
	}
	return siteFromAddress(ctx, txn, user, a.ID)
}
func siteFromLocation(ctx context.Context, txn bob.Tx, user User, location Location) (*models.Site, error) {
	// Reverse geocode at the location
	resp, err := geocode.ReverseGeocode(ctx, location)
	if err != nil {
		return nil, fmt.Errorf("reverse geocode: %w", err)
	}
	// Ensure we have an address at that newly created location
	a, err := geocode.EnsureAddress(ctx, txn, resp.Address, resp.Location)
	if err != nil {
		return nil, fmt.Errorf("ensure address: %w", err)
	}
	return siteFromAddress(ctx, txn, user, a.ID)
}
