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
	"github.com/Gleipnir-Technology/nidus-sync/db/enums"
	"github.com/Gleipnir-Technology/nidus-sync/db/models"
	"github.com/Gleipnir-Technology/nidus-sync/h3utils"
	"github.com/Gleipnir-Technology/nidus-sync/stadia"
	"github.com/stephenafamo/scan"
	//"github.com/rs/zerolog/log"
	"github.com/uber/h3-go/v4"
)

type Address struct {
	Country    enums.Countrytype
	Locality   string
	Number     string
	PostalCode string
	Region     string
	Street     string
	Unit       string
}
type GeocodeResult struct {
	Address   Address
	Cell      h3.Cell
	Longitude float64
	Latitude  float64
}

func (a Address) String() string {
	return fmt.Sprintf("%s %s, %s, %s, %s, %s", a.Number, a.Street, a.Locality, a.Region, a.PostalCode, a.Country)
}

var client *stadia.StadiaMaps

func InitializeStadia(key string) {
	client = stadia.NewStadiaMaps(key)
}

// Either get an address that matches, or create a new address. Either way, return an address
// This will make a call to a structured geocode service, so it's slow.
func EnsureAddress(ctx context.Context, txn bob.Tx, org *models.Organization, a Address) (*models.Address, error) {
	address, err := models.Addresses.Query(
		models.SelectWhere.Addresses.Country.EQ(a.Country),
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
	geo, err := Geocode(ctx, org, a)
	if err != nil {
		return nil, fmt.Errorf("geocode: %w", err)
	}

	type _row struct {
		ID int32 `db:"id"`
	}
	created := time.Now()
	row, err := bob.One(ctx, txn, psql.Insert(
		im.Into("address", "country", "created", "geom", "h3cell", "id", "locality", "number_", "postal_code", "region", "street", "unit"),
		im.Values(
			psql.Arg(geo.Address.Country),
			psql.Arg(created),
			psql.F("ST_Point", geo.Longitude, geo.Latitude, 4326),
			psql.Arg(geo.Cell),
			psql.Raw("DEFAULT"),
			psql.Arg(geo.Address.Locality),
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
		Country:    geo.Address.Country,
		Created:    created,
		Geom:       "",
		H3cell:     "",
		ID:         row.ID,
		Locality:   geo.Address.Locality,
		PostalCode: geo.Address.PostalCode,
		Street:     geo.Address.Street,
		Unit:       geo.Address.Unit,
		Region:     geo.Address.Region,
		Number:     geo.Address.Number,
	}, nil
}

func Geocode(ctx context.Context, org *models.Organization, a Address) (GeocodeResult, error) {
	street := fmt.Sprintf("%s %s", a.Number, a.Street)
	country_s := a.Country.String()
	/*
		sublog := log.With().
			Str("street", street).
			Str("country", country).
			Str("locality", a.Locality).
			Str("postal", a.PostalCode).
			Str("region", a.Region).
			Logger()
	*/
	req := stadia.StructuredGeocodeRequest{
		Address:    &street,
		Country:    &country_s,
		Locality:   &a.Locality,
		PostalCode: &a.PostalCode,
		Region:     &a.Region,
	}
	maybeAddServiceArea(&req, org)
	resp, err := client.StructuredGeocode(ctx, req)
	if err != nil {
		return GeocodeResult{}, fmt.Errorf("client structured geocode failure on %s: %w", a.String(), err)
	}
	if len(resp.Features) > 1 {
		return GeocodeResult{}, fmt.Errorf("%s matched more than one location", a.String())
	}
	feature := resp.Features[0]
	if feature.Geometry.Type != "Point" {
		return GeocodeResult{}, fmt.Errorf("wrong type %s from %s", feature.Geometry.Type, a.String())
	}
	longitude := feature.Geometry.Coordinates[0]
	latitude := feature.Geometry.Coordinates[1]
	cell, err := h3utils.GetCell(longitude, latitude, 15)
	if err != nil {
		return GeocodeResult{}, fmt.Errorf("failed to convert lat %f lng %f to h3 cell", longitude, latitude)
	}
	var country enums.Countrytype
	country_s = strings.ToLower(feature.Properties.CountryA)
	err = country.Scan(country_s)
	if err != nil {
		return GeocodeResult{}, fmt.Errorf("failed to scan country '%s': %w", country_s, err)
	}
	return GeocodeResult{
		Address: Address{
			Country:    country,
			Locality:   feature.Properties.Locality,
			Number:     feature.Properties.HouseNumber,
			PostalCode: feature.Properties.PostalCode,
			Region:     feature.Properties.Region,
			Street:     feature.Properties.Street,
			Unit:       "",
		},
		Cell:      cell,
		Longitude: feature.Geometry.Coordinates[0],
		Latitude:  feature.Geometry.Coordinates[1],
	}, nil
}

// Get the parcel for a given address, if one can be found
func GetParcel(ctx context.Context, txn bob.Tx, a *models.Address) (*models.Parcel, error) {
	result, err := models.Parcels.Query(
		sm.InnerJoin("address").On(psql.F("ST_Contains", psql.Raw("parcel.geometry"), psql.Raw("address.geom"))),
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
func maybeAddServiceArea(req *stadia.StructuredGeocodeRequest, org *models.Organization) {
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
	req.BoundaryRectMaxLon = &xmax
	req.BoundaryRectMaxLat = &ymax
	req.BoundaryRectMinLon = &xmin
	req.BoundaryRectMinLat = &ymin

	if org.ServiceAreaCentroidX.IsNull() || org.ServiceAreaCentroidY.IsNull() {
		return
	}
	centroid_x := org.ServiceAreaCentroidX.MustGet()
	centroid_y := org.ServiceAreaCentroidY.MustGet()

	req.FocusPointLat = &centroid_y
	req.FocusPointLng = &centroid_x
}
