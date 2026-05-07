package platform

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/Gleipnir-Technology/bob"
	"github.com/Gleipnir-Technology/bob/dialect/psql"
	"github.com/Gleipnir-Technology/bob/dialect/psql/dialect"
	"github.com/Gleipnir-Technology/bob/dialect/psql/sm"
	"github.com/Gleipnir-Technology/bob/types/pgtypes"
	"github.com/Gleipnir-Technology/nidus-sync/db"
	"github.com/Gleipnir-Technology/nidus-sync/db/gen/nidus-sync/public/model"
	"github.com/Gleipnir-Technology/nidus-sync/db/models"
	querypublic "github.com/Gleipnir-Technology/nidus-sync/db/query/public"
	nhttp "github.com/Gleipnir-Technology/nidus-sync/http"
	"github.com/Gleipnir-Technology/nidus-sync/platform/geocode"
	"github.com/Gleipnir-Technology/nidus-sync/platform/types"
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
func SiteByID(ctx context.Context, user User, id int32) (*types.Site, error) {
	query := siteQuery()
	query.Apply(
		sm.Where(models.Sites.Columns.ID.EQ(psql.Arg(id))),
		sm.Where(models.Sites.Columns.OrganizationID.EQ(psql.Arg(user.Organization.ID))),
	)
	sites, err := siteQueryToRows(ctx, query)
	if err != nil {
		return nil, err
	}
	return sites[id], nil
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
func SiteList(ctx context.Context, user User, limit int) ([]*types.Site, error) {
	query := siteQuery()
	query.Apply(
		sm.Where(psql.Quote("site", "organization_id").EQ(psql.Arg(user.Organization.ID))),
		sm.OrderBy(models.Sites.Columns.Created),
		sm.Limit(limit),
	)
	return siteQueryToRows(ctx, query)
}
func SitesByID(ctx context.Context, ids []int32) (map[int32]*models.Site, error) {
	rows, err := models.Sites.Query(
		sm.Where(
			models.Sites.Columns.ID.EQ(psql.Any(ids)),
		),
	).All(ctx, db.PGInstance.BobDB)
	if err != nil {
		return nil, fmt.Errorf("query sites: %w", err)
	}
	results := make(map[int32]*models.Site, len(rows))
	for _, row := range rows {
		results[row.ID] = row
	}
	return results, err
}
func siteFromAddressRaw(ctx context.Context, txn db.Ex, user User, address string) (*model.Site, error) {
	// Geocode
	geo, err := geocode.GeocodeRaw(ctx, user.Organization.model, address)
	if err != nil {
		return nil, fmt.Errorf("geocode: %w", err)
	}
	a, err := geocode.EnsureAddress(ctx, txn, geo.Address)
	if err != nil {
		return nil, fmt.Errorf("ensure address: %w", err)
	}
	return querypublic.SiteFromAddressIDForOrg(ctx, txn, int64(*a.ID), int64(user.Organization.ID))
}
func siteFromLocation(ctx context.Context, txn db.Ex, user User, location types.Location) (*model.Site, error) {
	// Reverse geocode at the location
	resp, err := geocode.ReverseGeocode(ctx, location)
	if err != nil {
		return nil, fmt.Errorf("reverse geocode: %w", err)
	}
	// Ensure we have an address at that newly created location
	a, err := geocode.EnsureAddress(ctx, txn, resp.Address)
	if err != nil {
		return nil, fmt.Errorf("ensure address: %w", err)
	}
	return querypublic.SiteFromAddressIDForOrg(ctx, txn, int64(*a.ID), int64(user.Organization.ID))
}
func siteQuery() bob.BaseQuery[*dialect.SelectQuery] {
	return psql.Select(
		sm.Columns(
			"address.country AS \"address.country\"",
			"address.locality AS \"address.locality\"",
			"COALESCE(address.location_latitude, 0) AS \"address.location.latitude\"",
			"COALESCE(address.location_longitude, 0) AS \"address.location.longitude\"",
			"address.number_ AS \"address.number_\"",
			"address.postal_code AS \"address.postal_code\"",
			"address.region AS \"address.region\"",
			"address.street AS \"address.street\"",
			"address.unit AS \"address.unit\"",
			"site.created AS \"created\"",
			"site.id AS \"id\"",
			"site.notes AS \"notes\"",
			"site.owner_name AS \"owner.name\"",
			"site.owner_phone_e164 AS \"owner.phone\"",
			"COALESCE(site.parcel_id, 0) AS \"parcel.id\"",
			"COALESCE(parcel.apn, '') AS \"parcel.apn\"",
			"COALESCE(parcel.description, '') AS \"parcel.description\"",
		),
		sm.From("site"),
		sm.InnerJoin("address").OnEQ(
			psql.Quote("site", "address_id"),
			psql.Quote("address", "id"),
		),
		sm.LeftJoin("parcel").OnEQ(
			psql.Quote("site", "parcel_id"),
			psql.Quote("parcel", "id"),
		),
	)
}
func siteQueryToRows(ctx context.Context, query bob.BaseQuery[*dialect.SelectQuery]) ([]*types.Site, error) {
	rows, err := bob.All(ctx, db.PGInstance.BobDB, query, scan.StructMapper[types.Site]())
	if err != nil {
		return nil, fmt.Errorf("query sites: %w", err)
	}
	site_ids := make([]int32, len(rows))
	results := make([]*types.Site, len(rows))
	for i, row := range rows {
		results[i] = &row
		site_ids[i] = row.ID
	}
	features_by_site_id, err := featuresBySiteID(ctx, site_ids)
	if err != nil {
		return nil, fmt.Errorf("query features for sites: %w", err)
	}
	for _, result := range results {
		features, ok := features_by_site_id[result.ID]
		if !ok {
			return nil, fmt.Errorf("impossible")
		}
		result.Features = features
	}
	leads_by_site_id, err := leadsBySiteID(ctx, site_ids)
	if err != nil {
		return nil, fmt.Errorf("query leads for sites: %w", err)
	}
	for _, result := range results {
		leads, ok := leads_by_site_id[result.ID]
		if !ok {
			return nil, fmt.Errorf("impossible")
		}
		result.Leads = leads
	}

	return results, nil
}
