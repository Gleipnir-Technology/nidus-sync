package tomtom

import (
	"fmt"
)

type PointShort struct {
	Latitude  float64 `json:"lat"`
	Longitude float64 `json:"lon"`
}

func (ps PointShort) AsPoint() Point {
	return Point{
		Latitude:  ps.Latitude,
		Longitude: ps.Longitude,
	}
}

type GeocodeResult struct {
	Type            string          `json:"type"`
	ID              string          `json:"id"`
	Score           float64         `json:"score"`
	Dist            float64         `json:"dist"`
	MatchConfidence MatchConfidence `json:"matchConfidence"`
	Address         Address         `json:"address"`
	Position        PointShort      `json:"position"`
	Viewport        Viewport        `json:"viewport"`
	EntryPoints     []EntryPoint    `json:"entryPoints"`
}

// MatchConfidence represents the confidence score for a match
type MatchConfidence struct {
	Score float64 `json:"score"`
}

// Address contains detailed address information
type Address struct {
	StreetNumber                string `json:"streetNumber"`
	StreetName                  string `json:"streetName"`
	Municipality                string `json:"municipality"`
	CountrySecondarySubdivision string `json:"countrySecondarySubdivision"`
	CountrySubdivision          string `json:"countrySubdivision"`
	CountrySubdivisionName      string `json:"countrySubdivisionName"`
	CountrySubdivisionCode      string `json:"countrySubdivisionCode"`
	PostalCode                  string `json:"postalCode"`
	ExtendedPostalCode          string `json:"extendedPostalCode"`
	CountryCode                 string `json:"countryCode"`
	Country                     string `json:"country"`
	CountryCodeISO3             string `json:"countryCodeISO3"`
	FreeformAddress             string `json:"freeformAddress"`
	LocalName                   string `json:"localName"`
}

// Viewport defines a geographic bounding box
type Viewport struct {
	TopLeftPoint  PointShort `json:"topLeftPoint"`
	BtmRightPoint PointShort `json:"btmRightPoint"`
}

// EntryPoint contains information about a point of entry to a location
type EntryPoint struct {
	Type     string     `json:"type"`
	Position PointShort `json:"position"`
}
type GeocodeSummary struct {
	Query        string     `json:"query"`
	QueryType    string     `json:"queryType"`
	QueryTime    uint       `json:"queryTime"`
	NumResults   uint       `json:"numResults"`
	Offset       uint       `json:"offset"`
	TotalResults uint       `json:"totalResults"`
	FuzzyLevel   uint       `json:"fuzzyLevel"`
	GeoBias      PointShort `json:"geoBias"`
}
type GeocodeResponse struct {
	Summary GeocodeSummary  `json:"summary"`
	Results []GeocodeResult `json:"results"`
}

// CalculateRoute sends a route calculation request to TomTom API
func (c *TomTom) Geocode(address string) (*GeocodeResponse, error) {
	var result GeocodeResponse

	resp, err := c.client.R().
		SetResult(&result).
		SetPathParam("address", address).
		SetPathParam("urlBase", c.urlBase).
		SetQueryParam("key", c.APIKey).
		SetQueryParam("storeResult", "false").
		Get("https://{urlBase}/search/2/geocode/{address}.json")
	if err != nil {
		return nil, fmt.Errorf("calculate route get: %w", err)
	}
	if !resp.IsSuccess() {
		return nil, fmt.Errorf("calculate route status: %d", resp.Status)
	}

	return &result, nil
}
