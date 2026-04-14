package geocode

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/Gleipnir-Technology/bob"
	"github.com/Gleipnir-Technology/bob/dialect/psql"
	"github.com/Gleipnir-Technology/bob/dialect/psql/sm"
	bobtypes "github.com/Gleipnir-Technology/bob/types"
	"github.com/Gleipnir-Technology/nidus-sync/db"
	"github.com/Gleipnir-Technology/nidus-sync/db/models"
	"github.com/Gleipnir-Technology/nidus-sync/h3utils"
	"github.com/Gleipnir-Technology/nidus-sync/platform/types"
	"github.com/Gleipnir-Technology/nidus-sync/stadia"
	"github.com/aarondl/opt/omit"
	//"github.com/rs/zerolog/log"
	"github.com/uber/h3-go/v4"
	"resty.dev/v3"
)

type GeocodeResult struct {
	Address types.Address
	Cell    h3.Cell
}

var client *stadia.StadiaMaps

func InitializeStadia(key string) {
	client = stadia.NewStadiaMaps(key)
	client.AddResponseMiddleware(restyMiddleware)
}
func restyMiddleware(rclient *resty.Client, response *resty.Response) error {
	//log.Info().Msg("middleware")
	ctx := context.Background()
	var body bobtypes.JSON[json.RawMessage]
	err := body.UnmarshalJSON(response.Bytes())
	if err != nil {
		return fmt.Errorf("unmarshal json in middleware: %w", err)
	}
	models.StadiaAPIRequests.Insert(&models.StadiaAPIRequestSetter{
		CreatedAt: omit.From(time.Now()),
		Request:   omit.From(response.Request.URL),
		Response:  omit.From(body),
	}).One(ctx, db.PGInstance.BobDB)
	return nil
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
	insertAddresses(ctx, db.PGInstance.BobDB, resp.Features)
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
	insertAddresses(ctx, db.PGInstance.BobDB, resp.Features)
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
	insertAddresses(ctx, db.PGInstance.BobDB, resp.Features)
	return toGeocodeResult(*resp, location.String())

}
func toGeocodeResult(resp stadia.GeocodeResponse, address_msg string) (*GeocodeResult, error) {
	if len(resp.Features) < 1 {
		return nil, fmt.Errorf("%s matched no locations", address_msg)
	}
	feature := resp.Features[0]
	if len(resp.Features) > 1 {
		if !allFeaturesIdenticalEnough(resp.Features) {
			return nil, fmt.Errorf("%s matched more than one location, and they differ a lot", address_msg)
		}
	}
	if feature.Geometry.Type != "Point" {
		return nil, fmt.Errorf("wrong type %s from %s", feature.Geometry.Type, address_msg)
	}
	longitude := feature.Geometry.Coordinates[0]
	latitude := feature.Geometry.Coordinates[1]
	cell, err := h3utils.GetCell(longitude, latitude, 15)
	if err != nil {
		return nil, fmt.Errorf("failed to convert lat %f lng %f to h3 cell", longitude, latitude)
	}
	country_s := strings.ToLower(feature.Properties.CountryA)
	// Depending on what kind of request we made we'll get wildly different result structures
	// This first structure generally works for forword geocoding
	address := types.Address{
		Country:  country_s,
		GID:      feature.Properties.GID,
		Locality: feature.Properties.Locality,
		Location: &types.Location{
			Longitude: feature.Geometry.Coordinates[0],
			Latitude:  feature.Geometry.Coordinates[1],
		},
		Number:     feature.Properties.HouseNumber,
		PostalCode: feature.Properties.PostalCode,
		Region:     feature.Properties.Region,
		Raw:        feature.Properties.FormattedAddressLine,
		Street:     feature.Properties.Street,
		Unit:       "",
	}
	// If we don't have a locality, try populating for reverse geocoding
	if address.Country == "" {
		address.Country = strings.ToLower(feature.Properties.Context.ISO3166A3)
		address.Locality = feature.Properties.Context.WhosOnFirst.Locality.Name
		address.Number = feature.Properties.AddressComponents.Number
		address.PostalCode = feature.Properties.AddressComponents.PostalCode
		address.Street = feature.Properties.AddressComponents.Street
	}
	return &GeocodeResult{
		Address: address,
		Cell:    cell,
	}, nil
}

// Get the parcel for a given address, if one can be found
func GetParcel(ctx context.Context, txn bob.Executor, a types.Address) (*models.Parcel, error) {
	if a.ID == nil {
		return nil, fmt.Errorf("nil address ID")
	}
	result, err := models.Parcels.Query(
		sm.InnerJoin("address").On(psql.F("ST_Contains", psql.Raw("parcel.geometry"), psql.Raw("address.location"))),
		models.SelectWhere.Addresses.ID.EQ(*a.ID),
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
	if org == nil {
		return
	}
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
