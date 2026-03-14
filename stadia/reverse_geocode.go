package stadia

import (
	"context"
	"fmt"

	"github.com/google/go-querystring/query"
)

type RequestReverseGeocode struct {
	Latitude  float64 `url:"point.lat" json:"point.lat"`
	Longitude float64 `url:"point.lon" json:"point.lon"`

	// Boundary circle parameters
	BoundaryCircleRadius *float64 `url:"boundary.circle.radius,omitempty"`
	BoundaryCountry      []string `url:"boundary.country,omitempty"`
	BoundaryGID          string   `url:"boundary.gid,omitempty"`

	// Other parameters
	Layers  []string `url:"layers,omitempty,comma" json:"layers,omitempty"`
	Size    *int     `url:"size,omitempty" json:"size,omitempty"`
	Sources []string `url:"sources,omitempty,comma" json:"sources,omitempty"`
}

func (s *StadiaMaps) ReverseGeocode(ctx context.Context, req RequestReverseGeocode) (*GeocodeResponse, error) {
	// https://docs.stadiamaps.com/geocoding-search-autocomplete/reverse-search/
	var result GeocodeResponse

	query, err := query.Values(req)
	if err != nil {
		return nil, fmt.Errorf("reverse geocode query: %w", err)
	}
	//var api_error Error
	resp, err := s.client.R().
		SetQueryParamsFromValues(query).
		SetContext(ctx).
		SetResult(&result).
		SetPathParam("urlBase", s.urlBase).
		SetQueryParam("api_key", s.APIKey).
		Get("https://{urlBase}/geocoding/v2/reverse")
	if err != nil {
		return nil, fmt.Errorf("reverse geocoding get: %w", err)
	}

	if !resp.IsSuccess() {
		return nil, parseError(resp)
	}
	return &result, nil
}
