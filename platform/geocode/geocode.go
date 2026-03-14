package geocode

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/Gleipnir-Technology/bob"
	"github.com/Gleipnir-Technology/bob/dialect/psql"
	"github.com/Gleipnir-Technology/bob/dialect/psql/im"
	"github.com/Gleipnir-Technology/bob/dialect/psql/sm"
	"github.com/Gleipnir-Technology/nidus-sync/db/models"
	"github.com/Gleipnir-Technology/nidus-sync/h3utils"
	"github.com/Gleipnir-Technology/nidus-sync/platform/types"
	"github.com/Gleipnir-Technology/nidus-sync/stadia"
	"github.com/stephenafamo/scan"
	//"github.com/rs/zerolog/log"
	"github.com/uber/h3-go/v4"
)

type GeocodeResult struct {
	Address  types.Address
	Cell     h3.Cell
	Location types.Location
}

var client *stadia.StadiaMaps

func InitializeStadia(key string) {
	client = stadia.NewStadiaMaps(key)
}

// Ensure the provided address exists. If it doesn't add it to the database.
func EnsureAddress(ctx context.Context, txn bob.Tx, a types.Address, l types.Location) (*models.Address, error) {
	address, err := models.Addresses.Query(
		models.SelectWhere.Addresses.Country.EQ(a.CountryEnum()),
		models.SelectWhere.Addresses.Locality.EQ(a.Locality),
		models.SelectWhere.Addresses.Number.EQ(a.Number),
		models.SelectWhere.Addresses.PostalCode.EQ(a.PostalCode),
		models.SelectWhere.Addresses.Region.EQ(a.Region),
		models.SelectWhere.Addresses.Street.EQ(a.Street),
		models.SelectWhere.Addresses.Unit.EQ(a.Unit),
	).One(ctx, txn)
	if err == nil {
		return address, nil
	}
	cell, err := h3utils.GetCell(l.Longitude, l.Latitude, 15)
	if err != nil {
		return nil, fmt.Errorf("failed to convert lat %f lng %f to h3 cell", l.Longitude, l.Latitude)
	}
	type _row struct {
		ID int32 `db:"id"`
	}
	created := time.Now()
	row, err := bob.One(ctx, txn, psql.Insert(
		im.Into("address", "country", "created", "h3cell", "id", "locality", "location", "number_", "postal_code", "region", "street", "unit"),
		im.Values(
			psql.Arg(a.CountryEnum()),
			psql.Arg(created),
			psql.Arg(cell),
			psql.Raw("DEFAULT"),
			psql.Arg(a.Locality),
			psql.F("ST_Point", l.Longitude, l.Latitude, 4326),
			psql.Arg(a.Number),
			psql.Arg(a.PostalCode),
			psql.Arg(a.Region),
			psql.Arg(a.Street),
			psql.Raw("''"),
		),
		im.Returning("id"),
	), scan.StructMapper[_row]())
	if err != nil {
		return nil, fmt.Errorf("insert: %w", err)
	}
	return &models.Address{
		Country:    a.CountryEnum(),
		Created:    created,
		H3cell:     "",
		ID:         row.ID,
		Locality:   a.Locality,
		Location:   "",
		PostalCode: a.PostalCode,
		Street:     a.Street,
		Unit:       a.Unit,
		Region:     a.Region,
		Number:     a.Number,
	}, nil
}

// Either get an address that matches, or create a new address. Either way, return an address
// This will make a call to a structured geocode service, so it's slow.
func EnsureAddressWithGeocode(ctx context.Context, txn bob.Tx, org *models.Organization, a types.Address) (*models.Address, error) {
	address, err := models.Addresses.Query(
		models.SelectWhere.Addresses.Country.EQ(a.CountryEnum()),
		models.SelectWhere.Addresses.Locality.EQ(a.Locality),
		models.SelectWhere.Addresses.Number.EQ(a.Number),
		models.SelectWhere.Addresses.PostalCode.EQ(a.PostalCode),
		models.SelectWhere.Addresses.Region.EQ(a.Region),
		models.SelectWhere.Addresses.Street.EQ(a.Street),
		models.SelectWhere.Addresses.Unit.EQ(a.Unit),
	).One(ctx, txn)
	if err == nil {
		return address, nil
	}
	// Geocode
	geo, err := GeocodeStructured(ctx, org, a)
	if err != nil {
		return nil, fmt.Errorf("geocode: %w", err)
	}

	type _row struct {
		ID int32 `db:"id"`
	}
	created := time.Now()
	row, err := bob.One(ctx, txn, psql.Insert(
		im.Into("address", "country", "created", "h3cell", "id", "locality", "location", "number_", "postal_code", "region", "street", "unit"),
		im.Values(
			psql.Arg(geo.Address.Country),
			psql.Arg(created),
			psql.Arg(geo.Cell),
			psql.Raw("DEFAULT"),
			psql.Arg(geo.Address.Locality),
			psql.F("ST_Point", geo.Location.Longitude, geo.Location.Latitude, 4326),
			psql.Arg(geo.Address.Number),
			psql.Arg(geo.Address.PostalCode),
			psql.Arg(geo.Address.Region),
			psql.Arg(geo.Address.Street),
			psql.Raw("''"),
		),
		im.Returning("id"),
	), scan.StructMapper[_row]())
	if err != nil {
		return nil, fmt.Errorf("insert: %w", err)
	}

	return &models.Address{
		Country:    geo.Address.CountryEnum(),
		Created:    created,
		H3cell:     "",
		ID:         row.ID,
		Locality:   geo.Address.Locality,
		Location:   "",
		PostalCode: geo.Address.PostalCode,
		Street:     geo.Address.Street,
		Unit:       geo.Address.Unit,
		Region:     geo.Address.Region,
		Number:     geo.Address.Number,
	}, nil
}
func GeocodeRaw(ctx context.Context, org *models.Organization, address string) (*GeocodeResult, error) {
	req := stadia.RequestGeocodeRaw{
		Text: address,
	}
	maybeAddServiceArea(&req, org)
	resp, err := client.GeocodeRaw(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("client raw geocode failure on %s: %w", address, err)
	}
	return toGeocodeResult(*resp, address)
}
func GeocodeStructured(ctx context.Context, org *models.Organization, a types.Address) (*GeocodeResult, error) {
	street := fmt.Sprintf("%s %s", a.Number, a.Street)
	req := stadia.RequestGeocodeStructured{
		Address: &street,
		//Country:    &a.Country,
		Locality:   &a.Locality,
		PostalCode: &a.PostalCode,
		Region:     &a.Region,
	}
	maybeAddServiceArea(&req, org)
	resp, err := client.GeocodeStructured(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("client structured geocode failure on %s: %w", a.String(), err)
	}
	return toGeocodeResult(*resp, a.String())
}
func ReverseGeocode(ctx context.Context, location types.Location) (*GeocodeResult, error) {
	req := stadia.RequestReverseGeocode{
		Latitude:  location.Latitude,
		Longitude: location.Longitude,
	}
	resp, err := client.ReverseGeocode(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("client reverse geocode failure on %s: %w", location.String(), err)
	}
	return toGeocodeResult(*resp, location.String())

}
func toGeocodeResult(resp stadia.GeocodeResponse, address string) (*GeocodeResult, error) {
	if len(resp.Features) < 1 {
		return nil, fmt.Errorf("%s matched no locations", address)
	}
	feature := resp.Features[0]
	if len(resp.Features) > 1 {
		if !allFeaturesIdenticalEnough(resp.Features) {
			return nil, fmt.Errorf("%s matched more than one location, and they differ a lot", address)
		}
	}
	if feature.Geometry.Type != "Point" {
		return nil, fmt.Errorf("wrong type %s from %s", feature.Geometry.Type, address)
	}
	longitude := feature.Geometry.Coordinates[0]
	latitude := feature.Geometry.Coordinates[1]
	cell, err := h3utils.GetCell(longitude, latitude, 15)
	if err != nil {
		return nil, fmt.Errorf("failed to convert lat %f lng %f to h3 cell", longitude, latitude)
	}
	country_s := strings.ToLower(feature.Properties.CountryA)
	return &GeocodeResult{
		Address: types.Address{
			Country:    country_s,
			Locality:   feature.Properties.Locality,
			Number:     feature.Properties.HouseNumber,
			PostalCode: feature.Properties.PostalCode,
			Region:     feature.Properties.Region,
			Street:     feature.Properties.Street,
			Unit:       "",
		},
		Cell: cell,
		Location: types.Location{
			Longitude: feature.Geometry.Coordinates[0],
			Latitude:  feature.Geometry.Coordinates[1],
		},
	}, nil
}

// Get the parcel for a given address, if one can be found
func GetParcel(ctx context.Context, txn bob.Tx, a *models.Address) (*models.Parcel, error) {
	result, err := models.Parcels.Query(
		sm.InnerJoin("address").On(psql.F("ST_Contains", psql.Raw("parcel.geometry"), psql.Raw("address.location"))),
		models.SelectWhere.Addresses.ID.EQ(a.ID),
	).One(ctx, txn)
	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			return nil, nil
		}
		return nil, fmt.Errorf("Get parcel from address %d: %w", a.ID, err)
	}
	return result, nil
}
func allFeaturesIdenticalEnough(features []stadia.GeocodeFeature) bool {
	if len(features) < 2 {
		return true
	}
	f := features[0].Properties
	for _, feature := range features {
		if feature.Properties.CountryCode != f.CountryCode ||
			feature.Properties.County != f.County ||
			feature.Properties.HouseNumber != f.HouseNumber ||
			feature.Properties.Locality != f.Locality ||
			feature.Properties.RegionA != f.RegionA {
			return false
		}
	}
	return true
}
func maybeAddServiceArea(req stadia.RequestGeocode, org *models.Organization) {
	if org.ServiceAreaXmax.IsNull() ||
		org.ServiceAreaYmax.IsNull() ||
		org.ServiceAreaXmin.IsNull() ||
		org.ServiceAreaYmin.IsNull() {
		return
	}
	xmax := org.ServiceAreaXmax.MustGet()
	ymax := org.ServiceAreaYmax.MustGet()
	xmin := org.ServiceAreaXmin.MustGet()
	ymin := org.ServiceAreaYmin.MustGet()
	req.SetBoundaryRect(xmin, ymin, xmax, ymax)

	if org.ServiceAreaCentroidX.IsNull() || org.ServiceAreaCentroidY.IsNull() {
		return
	}
	centroid_x := org.ServiceAreaCentroidX.MustGet()
	centroid_y := org.ServiceAreaCentroidY.MustGet()

	req.SetFocusPoint(centroid_x, centroid_y)
}
