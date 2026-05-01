package geocode

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"sort"
	"time"

	"github.com/Gleipnir-Technology/bob"
	"github.com/Gleipnir-Technology/bob/dialect/psql"
	"github.com/Gleipnir-Technology/bob/dialect/psql/sm"
	bobtypes "github.com/Gleipnir-Technology/bob/types"
	"github.com/Gleipnir-Technology/nidus-sync/db"
	"github.com/Gleipnir-Technology/nidus-sync/db/gen/nidus-sync/stadia/model"
	"github.com/Gleipnir-Technology/nidus-sync/db/gen/nidus-sync/stadia/table"
	"github.com/Gleipnir-Technology/nidus-sync/db/models"
	"github.com/Gleipnir-Technology/nidus-sync/h3utils"
	"github.com/Gleipnir-Technology/nidus-sync/platform/types"
	"github.com/Gleipnir-Technology/nidus-sync/stadia"
	//"github.com/aarondl/opt/omit"
	"github.com/rs/zerolog/log"
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
func redactQueryParam(u string, param string) (string, error) {
	parsedURL, err := url.Parse(u)
	if err != nil {
		return "", fmt.Errorf("failed to parse URL: %w", err)
	}

	queryParams := parsedURL.Query()
	queryParams.Del(param)
	parsedURL.RawQuery = queryParams.Encode()

	return parsedURL.String(), nil
}
func restyMiddleware(rclient *resty.Client, response *resty.Response) error {
	//log.Info().Msg("middleware")
	ctx := context.Background()
	var body bobtypes.JSON[json.RawMessage]
	resp_bytes := response.Bytes()
	err := body.UnmarshalJSON(resp_bytes)
	if err != nil {
		return fmt.Errorf("unmarshal json in middleware: %w", err)
	}
	u, err := redactQueryParam(response.Request.URL, "api_key")
	if err != nil {
		log.Error().Err(err).Str("url", response.Request.URL).Msg("failed to redact url")
		return nil
	}
	statement := table.APIRequest.INSERT(table.APIRequest.MutableColumns).
		MODEL(model.APIRequest{
			CreatedAt: time.Now(),
			Request:   u,
			Response:  string(resp_bytes),
		}).RETURNING(table.APIRequest.AllColumns)
	data, err := db.ExecuteOne[model.APIRequest](ctx, statement)
	if err != nil {
		log.Error().Err(err).Msg("failed to insert stadia request")
	} else {
		log.Debug().Int64("id", data.ID).Msg("Created stadia request cache entry")
	}
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
	addresses, err := insertAddresses(ctx, db.PGInstance.BobDB, resp.Features)
	if err != nil {
		return nil, fmt.Errorf("insert addresses: %w", err)
	}
	return toGeocodeResult(resp.Features, address, addresses)
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
	addresses, err := insertAddresses(ctx, db.PGInstance.BobDB, resp.Features)
	if err != nil {
		return nil, fmt.Errorf("insert addresses: %w", err)
	}
	return toGeocodeResult(resp.Features, a.String(), addresses)
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
func ReverseGeocode(ctx context.Context, location types.Location) (*GeocodeResult, error) {
	req := stadia.RequestReverseGeocode{
		Latitude:  location.Latitude,
		Longitude: location.Longitude,
	}
	resp, err := client.ReverseGeocode(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("client reverse geocode failure on %s: %w", location.String(), err)
	}
	addresses, err := insertAddresses(ctx, db.PGInstance.BobDB, resp.Features)
	if err != nil {
		return nil, fmt.Errorf("insert addresses: %w", err)
	}
	return toGeocodeResult(resp.Features, location.String(), addresses)

}
func ReverseGeocodeClosest(ctx context.Context, location types.Location) (*GeocodeResult, error) {
	req := stadia.RequestReverseGeocode{
		Latitude:  location.Latitude,
		Longitude: location.Longitude,
	}
	resp, err := client.ReverseGeocode(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("client reverse geocode failure on %s: %w", location.String(), err)
	}
	addresses, err := insertAddresses(ctx, db.PGInstance.BobDB, resp.Features)
	if err != nil {
		return nil, fmt.Errorf("insert addresses: %w", err)
	}
	/*
		sorter := SortAddressByDistance{
			Addresses: addresses,
			Location: location,
		}
	*/
	sort.Sort(SortFeaturesByDistance(resp.Features))
	return toGeocodeResult(resp.Features[:1], location.String(), addresses[:1])

}
func toGeocodeResult(features []stadia.GeocodeFeature, address_msg string, addresses []types.Address) (*GeocodeResult, error) {
	if len(features) < 1 {
		return nil, fmt.Errorf("%s matched no locations", address_msg)
	}
	if len(addresses) < 1 {
		return nil, fmt.Errorf("no addresses")
	}
	if len(features) > 1 {
		if !allFeaturesIdenticalEnough(features) {
			return nil, fmt.Errorf("%s matched more than one location, and they differ a lot", address_msg)
		}
	}
	feature := features[0]
	address := addresses[0]
	if feature.Geometry.Type != "Point" {
		return nil, fmt.Errorf("wrong type %s from %s", feature.Geometry.Type, address_msg)
	}
	cell, err := h3utils.GetCell(address.Location.Longitude, address.Location.Latitude, 15)
	if err != nil {
		return nil, fmt.Errorf("failed to convert lat %f lng %f to h3 cell", address.Location.Longitude, address.Location.Latitude)
	}
	return &GeocodeResult{
		Address: address,
		Cell:    cell,
	}, nil
}
func allFeaturesIdenticalEnough(features []stadia.GeocodeFeature) bool {
	if len(features) < 2 {
		return true
	}
	// We may have a 'fallback feature' which is a match for a street, but not the rest
	// We therefore search for a feature that is 'full' or has all the properties we care about
	// So long as we find only one, and no conflicts, we're fine.
	var full_feature *stadia.GeocodeFeature
	for _, feature := range features {
		if feature.Number() != "" &&
			feature.Street() != "" &&
			feature.Locality() != "" &&
			feature.Region() != "" {
			if full_feature == nil {
				full_feature = &feature
			} else if feature.Number() == full_feature.Number() &&
				feature.Street() == full_feature.Street() &&
				feature.Locality() == full_feature.Locality() &&
				feature.Region() == full_feature.Region() {
				continue
			} else {
				log.Info().
					Str("ff_number", full_feature.Number()).
					Str("ff_street", full_feature.Street()).
					Str("ff_locality", full_feature.Locality()).
					Str("ff_region", full_feature.Region()).
					Str("number", feature.Number()).
					Str("street", feature.Street()).
					Str("locality", feature.Locality()).
					Str("region", feature.Region()).
					Msg("This is where the features first disagreed")
				return false
			}
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
