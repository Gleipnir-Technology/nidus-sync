package stadia

import (
	"context"
	"fmt"

	"github.com/google/go-querystring/query"
)

type RequestGeocodeAutocomplete struct {
	Text string `url:"text" json:"text"`

	// Boundary circle parameters
	BoundaryCircleLat    *float64 `url:"boundary.circle.lat,omitempty"`
	BoundaryCircleLon    *float64 `url:"boundary.circle.lon,omitempty"`
	BoundaryCircleRadius *float64 `url:"boundary.circle.radius,omitempty"`

	BoundaryCountry *string `url:"boundary.country,omitempty"` //comma-delimited ISO 2 or 3 character code
	BoundaryGID     *string `url:"boundary.gid,omitempty"`     // The GID of a region to limit the search to

	// Boundary parameters
	BoundaryRectMaxLat *float64 `url:"boundary.rect.max_lat,omitempty"`
	BoundaryRectMinLat *float64 `url:"boundary.rect.min_lat,omitempty"`
	BoundaryRectMaxLon *float64 `url:"boundary.rect.max_lon,omitempty"`
	BoundaryRectMinLon *float64 `url:"boundary.rect.min_lon,omitempty"`

	// Focus point
	FocusPointLat *float64 `url:"focus.point.lat,omitempty" json:",omitempty"`
	FocusPointLng *float64 `url:"focus.point.lon,omitempty" json:",omitempty"`

	// Other parameters
	Lang    *string  `url:"lang,omitempty" json:"lang,omitempty"`
	Layers  []string `url:"layers,omitempty,comma" json:"layers,omitempty"`
	Size    *int     `url:"size,omitempty" json:"size,omitempty"`
	Sources []string `url:"sources,omitempty,comma" json:"sources,omitempty"`
}

func (r *RequestGeocodeAutocomplete) SetBoundaryRect(xmin, ymin, xmax, ymax float64) {
	r.BoundaryRectMaxLat = &ymax
	r.BoundaryRectMinLat = &ymin
	r.BoundaryRectMaxLon = &xmax
	r.BoundaryRectMinLon = &xmin
}
func (r *RequestGeocodeAutocomplete) SetFocusPoint(x, y float64) {
	r.FocusPointLat = &y
	r.FocusPointLng = &x
}
func (s *StadiaMaps) GeocodeAutocomplete(ctx context.Context, req RequestGeocodeAutocomplete) (*GeocodeResponse, error) {
	// https://docs.stadiamaps.com/geocoding-search-autocomplete/search/
	var result GeocodeResponse

	query, err := query.Values(req)
	if err != nil {
		return nil, fmt.Errorf("structured geocode query: %w", err)
	}
	//var api_error Error
	resp, err := s.client.R().
		SetQueryParamsFromValues(query).
		SetContext(ctx).
		SetResult(&result).
		SetPathParam("urlBase", s.urlBase).
		SetQueryParam("api_key", s.APIKey).
		Get("https://{urlBase}/geocoding/v2/autocomplete")
	if err != nil {
		return nil, fmt.Errorf("autocomplete get: %w", err)
	}

	if !resp.IsSuccess() {
		return nil, parseError(resp)
	}
	return &result, nil
}
