package tomtom

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

// Base URLs and API constants
const (
	BaseURL          = "https://api.tomtom.com"
	RouteTypeFastest = "fastest"
	TravelModeCar    = "car"
)

// Coordinates represents latitude and longitude values
type Coordinates struct {
	Latitude  string `json:"latitude" xml:"latitude,attr"`
	Longitude string `json:"longitude" xml:"longitude,attr"`
}

// Rectangle represents a geographic rectangle
type Rectangle struct {
	SouthWestCorner Coordinates `json:"southWestCorner" xml:"southWestCorner"`
	NorthEastCorner Coordinates `json:"northEastCorner" xml:"northEastCorner"`
}

// AvoidAreas represents areas to avoid in routing
type AvoidAreas struct {
	Rectangles []Rectangle `json:"rectangles" xml:"rectangles>rectangle"`
}

// SupportingPoint represents a supporting point in the route calculation
type SupportingPoint struct {
	Latitude  string `json:"latitude" xml:"latitude,attr"`
	Longitude string `json:"longitude" xml:"longitude,attr"`
}

// Client represents a TomTom API client
type Client struct {
	APIKey     string
	HTTPClient *http.Client
}

// CalculateRoutePostData represents the POST body for Calculate Route API
type CalculateRoutePostData struct {
	SupportingPoints []SupportingPoint `json:"supportingPoints,omitempty" xml:"supportingPoints>supportingPoint,omitempty"`
	AvoidVignette    []string          `json:"avoidVignette,omitempty" xml:"avoidVignette,omitempty"`
	AllowVignette    []string          `json:"allowVignette,omitempty" xml:"allowVignette,omitempty"`
	AvoidAreas       *AvoidAreas       `json:"avoidAreas,omitempty" xml:"avoidAreas,omitempty"`
}

// Route response structures - These would need to be completed based on actual API response
type Summary struct {
	LengthInMeters                 int     `json:"lengthInMeters"`
	TravelTimeInSeconds            int     `json:"travelTimeInSeconds"`
	TrafficDelayInSeconds          int     `json:"trafficDelayInSeconds"`
	DepartureTime                  string  `json:"departureTime"`
	ArrivalTime                    string  `json:"arrivalTime"`
	FuelConsumptionInLiters        float64 `json:"fuelConsumptionInLiters,omitempty"`
	ElectricEnergyConsumptionInkWh float64 `json:"electricEnergyConsumptionInkWh,omitempty"`
}

type Point struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

type Polyline struct {
	Points []Point `json:"points"`
}

type Leg struct {
	Summary  Summary  `json:"summary"`
	Points   []Point  `json:"points,omitempty"`
	Polyline Polyline `json:"polyline,omitempty"`
}

type Route struct {
	Summary  Summary  `json:"summary"`
	Legs     []Leg    `json:"legs,omitempty"`
	Polyline Polyline `json:"polyline,omitempty"`
}

type CalculateRouteResponse struct {
	Routes []Route `json:"routes"`
}

// CalculateReachableRange API structures

// CalculateReachableRangeParams holds the parameters for the Calculate Reachable Range API
type CalculateReachableRangeParams struct {
	// Path parameters
	Origin      string
	ContentType string // "json" or "jsonp"

	// Query parameters
	FuelBudgetInLiters                           *float64
	EnergyBudgetInkWh                            *float64
	TimeBudgetInSec                              *float64
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

// CalculateReachableRangePostData represents the POST body for Calculate Reachable Range API
type CalculateReachableRangePostData struct {
	AvoidVignette []string    `json:"avoidVignette,omitempty" xml:"avoidVignette,omitempty"`
	AllowVignette []string    `json:"allowVignette,omitempty" xml:"allowVignette,omitempty"`
	AvoidAreas    *AvoidAreas `json:"avoidAreas,omitempty" xml:"avoidAreas,omitempty"`
}

// CalculateReachableRangeRequest combines both path parameters and POST body
type CalculateReachableRangeRequest struct {
	Params   CalculateReachableRangeParams
	PostData *CalculateReachableRangePostData
}

// Reachable Range response structures
type Polygon struct {
	Exterior []Point   `json:"exterior"`
	Interior [][]Point `json:"interior,omitempty"`
}

type CalculateReachableRangeResponse struct {
	Polygon Polygon `json:"polygon"`
	Summary struct {
		DistanceLimit          float64 `json:"distanceLimit,omitempty"`
		TimeLimit              int     `json:"timeLimit,omitempty"`
		FuelConsumptionLimit   float64 `json:"fuelConsumptionLimit,omitempty"`
		EnergyConsumptionLimit float64 `json:"energyConsumptionLimit,omitempty"`
	} `json:"summary"`
}

// BuildURL builds the URL for the Calculate Reachable Range request
func (req *CalculateReachableRangeRequest) BuildURL(apiKey string) (string, error) {
	baseURL := fmt.Sprintf("%s/routing/%d/calculateReachableRange/%s/%s",
		BaseURL,
		req.Params.Origin,
		req.Params.ContentType)

	// Add query parameters
	query := url.Values{}
	query.Add("key", apiKey)

	if req.Params.FuelBudgetInLiters != nil {
		query.Add("fuelBudgetInLiters", fmt.Sprintf("%f", *req.Params.FuelBudgetInLiters))
	}

	if req.Params.EnergyBudgetInkWh != nil {
		query.Add("energyBudgetInkWh", fmt.Sprintf("%f", *req.Params.EnergyBudgetInkWh))
	}

	if req.Params.TimeBudgetInSec != nil {
		query.Add("timeBudgetInSec", fmt.Sprintf("%f", *req.Params.TimeBudgetInSec))
	}

	// Add other parameters similarly...

	return baseURL + "?" + query.Encode(), nil
}

// Client methods for executing requests

// CalculateReachableRange sends a reachable range calculation request to TomTom API
func (c *Client) CalculateReachableRange(req *CalculateReachableRangeRequest) (*CalculateReachableRangeResponse, error) {
	url, err := req.BuildURL(c.APIKey)
	if err != nil {
		return nil, err
	}

	var response CalculateReachableRangeResponse
	var httpReq *http.Request

	if req.PostData != nil {
		// POST request
		jsonData, err := json.Marshal(req.PostData)
		if err != nil {
			return nil, err
		}
		httpReq, err = http.NewRequest("POST", url, strings.NewReader(string(jsonData)))
		if err != nil {
			return nil, err
		}
		httpReq.Header.Set("Content-Type", "application/json")
	} else {
		// GET request
		httpReq, err = http.NewRequest("GET", url, nil)
		if err != nil {
			return nil, err
		}
	}

	resp, err := c.HTTPClient.Do(httpReq)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API error: %s", resp.Status)
	}

	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}
