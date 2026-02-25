package stadia

import (
	"context"
	"fmt"
	"strings"

	"github.com/google/go-querystring/query"
)

// StructuredGeocodeRequest represents the query parameters for structured geocoding
type StructuredGeocodeRequest struct {
	// Address components
	Address       *string `url:"address,omitempty" json:"address,omitempty"`
	Neighbourhood *string `url:"neighbourhood,omitempty" json:"neighbourhood,omitempty"`
	Borough       *string `url:"borough,omitempty" json:"borough,omitempty"`
	Locality      *string `url:"locality,omitempty" json:"locality,omitempty"`
	County        *string `url:"county,omitempty" json:"county,omitempty"`
	Region        *string `url:"region,omitempty" json:"region,omitempty"`
	PostalCode    *string `url:"postalcode,omitempty" json:"postalcode,omitempty"`
	Country       *string `url:"country,omitempty" json:"country,omitempty"`

	// Focus point
	FocusPoint *FocusPoint `url:",omitempty" json:",omitempty"`

	// Boundary parameters
	BoundaryRect    *BoundaryRect   `url:",omitempty" json:",omitempty"`
	BoundaryCircle  *BoundaryCircle `url:",omitempty" json:",omitempty"`
	BoundaryCountry []string        `url:"boundary.country,omitempty,comma" json:"boundary.country,omitempty,comma"`
	BoundaryGid     *string         `url:"boundary.gid,omitempty" json:"boundary.gid,omitempty"`

	// Other parameters
	Layers  []string `url:"layers,omitempty,comma" json:"layers,omitempty,comma"`
	Sources []string `url:"sources,omitempty,comma" json:"sources,omitempty,comma"`
	Size    *int     `url:"size,omitempty" json:"size,omitempty"`
	Lang    *string  `url:"lang,omitempty" json:"lang,omitempty"`
}

func (s *StadiaMaps) StructuredGeocode(ctx context.Context, req StructuredGeocodeRequest) (*GeocodeResponse, error) {
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
		SetError(&result).
		SetResult(&result).
		SetPathParam("urlBase", s.urlBase).
		SetQueryParam("api_key", s.APIKey).
		Get("https://{urlBase}/geocoding/v1/search/structured")
	if err != nil {
		return nil, fmt.Errorf("structured geocoding get: %w", err)
	}

	if !resp.IsSuccess() {
		/*
			if api_error.Error() != "" {
				return nil, &api_error
			}
			content, err := io.ReadAll(resp.Body)
			if err != nil {
				return nil, fmt.Errorf("read all failure: %w", err)
			}
		*/
		fmt.Printf("geocoding error: %s\n", result.Geocode.Error)
		if len(result.Geocode.Errors) > 0 {
			joined := strings.Join(result.Geocode.Errors, ", ")
			return nil, fmt.Errorf("structured geocoding failure: %d '%s'", resp.StatusCode(), joined)
		} else if result.Geocode.Error != "" {
			return nil, fmt.Errorf("structured geocoding failure: %d '%s'", resp.StatusCode(), result.Geocode.Error)
		} else {
			return nil, fmt.Errorf("structured geocoding failure: %d", resp.StatusCode())
		}
	}
	return &result, nil
}

func (sgr StructuredGeocodeRequest) endpoint() string {
	return "/v1/search/structured"
}
