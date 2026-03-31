package platform

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/Gleipnir-Technology/bob"
	"github.com/Gleipnir-Technology/bob/dialect/psql"
	"github.com/Gleipnir-Technology/bob/dialect/psql/sm"
	"github.com/Gleipnir-Technology/bob/types/pgtypes"
	"github.com/Gleipnir-Technology/nidus-sync/db"
	"github.com/Gleipnir-Technology/nidus-sync/db/models"
	nhttp "github.com/Gleipnir-Technology/nidus-sync/http"
	"github.com/Gleipnir-Technology/nidus-sync/platform/geocode"
	//"github.com/Gleipnir-Technology/nidus-sync/platform/types"
	"github.com/aarondl/opt/omit"
	"github.com/aarondl/opt/omitnull"
	"github.com/stephenafamo/scan"
)

func SiteFromSignal(ctx context.Context, user User, signal_id int32) (*int32, error) {
	type _Row struct {
		ID int32 `db:"site_id"`
	}
	site, err := bob.One(ctx, db.PGInstance.BobDB, psql.Select(
		sm.Columns(
			"pool.site_id AS site_id",
		),
		sm.From("signal_pool"),
		sm.InnerJoin("pool").OnEQ(
			psql.Quote("signal_pool", "pool_id"),
			psql.Quote("pool", "id"),
		),
		sm.InnerJoin("site").On(
			psql.Quote("pool", "site_id").EQ(psql.Quote("site", "id")),
		),
		sm.Where(psql.Quote("signal_pool", "signal_id").EQ(psql.Arg(signal_id))),
		sm.Where(psql.Quote("site", "organization_id").EQ(psql.Arg(user.Organization.ID))),
	), scan.StructMapper[_Row]())
	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			return nil, nhttp.NewErrorStatus(http.StatusBadRequest, "Can't make a lead from signal %d: %w", signal_id, err)
		}
		return nil, fmt.Errorf("failed getting site: %w", err)
	}
	return &site.ID, nil
}
func SiteCreate(ctx context.Context, txn bob.Tx, user User, address_id int32) (*models.Site, error) {
	return models.Sites.Insert(&models.SiteSetter{
		AddressID: omit.From(address_id),
		Created:   omit.From(time.Now()),
		CreatorID: omit.From(int32(user.ID)),
		FileID:    omitnull.FromPtr[int32](nil),
		//ID:
		Notes:          omit.From(""),
		OrganizationID: omit.From(user.Organization.ID),
		OwnerName:      omit.From(""),
		OwnerPhoneE164: omitnull.FromPtr[string](nil),
		ParcelID:       omitnull.FromPtr[int32](nil),
		ResidentOwned:  omitnull.FromPtr[bool](nil),
		Tags:           omit.From(pgtypes.HStore{}),
		Version:        omit.From(int32(1)),
	}).One(ctx, txn)
}
func siteFromAddress(ctx context.Context, txn bob.Tx, user User, address_id int32) (*models.Site, error) {
	site, err := models.Sites.Query(
		models.SelectWhere.Sites.AddressID.EQ(address_id),
		models.SelectWhere.Sites.OrganizationID.EQ(user.Organization.ID),
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
