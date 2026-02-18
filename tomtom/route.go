package tomtom

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/google/go-querystring/query"
)

// CalculateRouteRequest combines both path parameters and POST body
type CalculateRouteRequest struct {
	// Path parameters
	Locations Locations

	// Query parameters
	MaxAlternatives                              *int
	AlternativeType                              string
	MinDeviationDistance                         *int
	MinDeviationTime                             *int
	InstructionsType                             string
	Language                                     string
	ComputeBestOrder                             *bool
	RouteRepresentation                          string
	ComputeTravelTimeFor                         string
	VehicleHeading                               *int
	SectionType                                  []string
	IncludeTollPaymentTypes                      string
	Callback                                     string
	Report                                       string
	DepartAt                                     string
	ArriveAt                                     string
	RouteType                                    string
	Traffic                                      *bool
	Avoid                                        []string
	TravelMode                                   string
	Hilliness                                    string
	Windingness                                  string
	VehicleMaxSpeed                              *int
	VehicleWeight                                *int
	VehicleAxleWeight                            *int
	VehicleNumberOfAxles                         *int
	VehicleLength                                *float64
	VehicleWidth                                 *float64
	VehicleHeight                                *float64
	VehicleCommercial                            *bool
	VehicleLoadType                              []string
	VehicleAdrTunnelRestrictionCode              string
	VehicleHasElectricTollCollectionTransponder  string
	VehicleEngineType                            string
	ConstantSpeedConsumptionInLitersPerHundredkm string
	CurrentFuelInLiters                          *float64
	AuxiliaryPowerInLitersPerHour                *float64
	FuelEnergyDensityInMJoulesPerLiter           *float64
	AccelerationEfficiency                       *float64
	DecelerationEfficiency                       *float64
	UphillEfficiency                             *float64
	DownhillEfficiency                           *float64
	ConsumptionInkWhPerkmAltitudeGain            *float64
	RecuperationInkWhPerkmAltitudeLoss           *float64
	ConstantSpeedConsumptionInkWhPerHundredkm    string
	CurrentChargeInkWh                           *float64
	MaxChargeInkWh                               *float64
	AuxiliaryPowerInkW                           *float64
}

func (sgr CalculateRouteRequest) toQueryParams() (url.Values, error) {
	return query.Values(sgr)
}

type CalculateRouteResponse struct {
	Routes []Route `json:"routes"`
}

// CalculateRoute sends a route calculation request to TomTom API
func (c *TomTom) CalculateRoute(req *CalculateRouteRequest) (*CalculateRouteResponse, error) {
	/*url, err := req.BuildURL(c.APIKey)
	if err != nil {
		return nil, err
	}*/

	var result CalculateRouteResponse
	/*
		query, err := req.toQueryParams()
		if err != nil {
			return nil, fmt.Errorf("structured geocode query: %w", err)
		}
	*/

	resp, err := c.client.R().
		//SetQueryParamsFromValues(query).
		SetResult(&result).
		SetPathParam("locations", locationString(req.Locations)).
		SetPathParam("urlBase", c.urlBase).
		SetQueryParam("key", c.APIKey).
		Get("https://{urlBase}/routing/1/calculateRoute/{locations}/json")
	if err != nil {
		return nil, fmt.Errorf("calculate route get: %w", err)
	}
	if !resp.IsSuccess() {
		return nil, fmt.Errorf("calculate route status: %d", resp.Status)
	}

	return &result, nil
}

func P(lat, lng float64) Point {
	return Point{
		Latitude:  lat,
		Longitude: lng,
	}
}

type Locations = []Point

func locationString(locations Locations) string {
	var sb strings.Builder
	for i, p := range locations {
		if i == 0 {
			sb.WriteString(fmt.Sprintf("%f,%f", p.Latitude, p.Longitude))
		} else {
			sb.WriteString(fmt.Sprintf(":%f,%f", p.Latitude, p.Longitude))
		}
	}
	return sb.String()
}
