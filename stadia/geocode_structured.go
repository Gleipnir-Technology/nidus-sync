package stadia

import (
	"context"
	"fmt"

	"github.com/google/go-querystring/query"
)

// RequestGeocodeStructured represents the query parameters for structured geocoding
type RequestGeocodeStructured struct {
	// Address components
	Address       *string `url:"address,omitempty" json:"address,omitempty"`
	Neighbourhood *string `url:"neighbourhood,omitempty" json:"neighbourhood,omitempty"`
	Borough       *string `url:"borough,omitempty" json:"borough,omitempty"`
	Locality      *string `url:"locality,omitempty" json:"locality,omitempty"`
	County        *string `url:"county,omitempty" json:"county,omitempty"`
	Region        *string `url:"region,omitempty" json:"region,omitempty"`
	PostalCode    *string `url:"postalcode,omitempty" json:"postalcode,omitempty"`
	Country       *string `url:"country,omitempty" json:"country,omitempty"`

	// Boundary circle parameters
	BoundaryCircleLat    *float64 `url:"boundary.circle.lat,omitempty"`
	BoundaryCircleLon    *float64 `url:"boundary.circle.lon,omitempty"`
	BoundaryCircleRadius *float64 `url:"boundary.circle.radius,omitempty"`

	BoundaryCountry []string `url:"boundary.country,omitempty,comma" json:"boundary.country,omitempty"`

	BoundaryGid *string `url:"boundary.gid,omitempty" json:"boundary.gid,omitempty"`
	// Boundary parameters
	BoundaryRectMaxLat *float64 `url:"boundary.rect.max_lat,omitempty"`
	BoundaryRectMinLat *float64 `url:"boundary.rect.min_lat,omitempty"`
	BoundaryRectMaxLon *float64 `url:"boundary.rect.max_lon,omitempty"`
	BoundaryRectMinLon *float64 `url:"boundary.rect.min_lon,omitempty"`

	// Focus point
	FocusPointLat *float64 `url:"focus.point.lat,omitempty" json:",omitempty"`
	FocusPointLng *float64 `url:"focus.point.lon,omitempty" json:",omitempty"`

	// Other parameters
	Layers  []string `url:"layers,omitempty,comma" json:"layers,omitempty"`
	Sources []string `url:"sources,omitempty,comma" json:"sources,omitempty"`
	Size    *int     `url:"size,omitempty" json:"size,omitempty"`
	Lang    *string  `url:"lang,omitempty" json:"lang,omitempty"`
}

func (r *RequestGeocodeStructured) SetBoundaryRect(xmin, ymin, xmax, ymax float64) {
	r.BoundaryRectMaxLat = &ymax
	r.BoundaryRectMinLat = &ymin
	r.BoundaryRectMaxLon = &xmax
	r.BoundaryRectMinLon = &xmin
}
func (r *RequestGeocodeStructured) SetFocusPoint(x, y float64) {
	r.FocusPointLat = &y
	r.FocusPointLng = &x
}

func (s *StadiaMaps) GeocodeStructured(ctx context.Context, req RequestGeocodeStructured) (*GeocodeResponse, error) {
	// https://docs.stadiamaps.com/geocoding-search-autocomplete/structured-search/
	// curl "https://api.stadiamaps.com/geocoding/v1/search/structured?address=P%C3%B5hja%20pst%2027a&region=Harju&country=EE&api_key=YOUR-API-KEY"
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
		Get("https://{urlBase}/geocoding/v1/search/structured")
	if err != nil {
		return nil, fmt.Errorf("structured geocoding get: %w", err)
	}

	if !resp.IsSuccess() {
		return nil, parseError(resp)
	}
	return &result, nil
}
