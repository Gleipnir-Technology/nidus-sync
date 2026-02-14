package stadia

import (
	"fmt"
	"net/url"

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

func (s *StadiaMaps) StructuredGeocode(req StructuredGeocodeRequest) (*GeocodeResponse, error) {
	// https://docs.stadiamaps.com/geocoding-search-autocomplete/structured-search/
	// curl "https://api.stadiamaps.com/geocoding/v1/search/structured?address=P%C3%B5hja%20pst%2027a&region=Harju&country=EE&api_key=YOUR-API-KEY"
	var result GeocodeResponse

	query, err := req.toQueryParams()
	if err != nil {
		return nil, fmt.Errorf("structured geocode query: %w", err)
	}
	resp, err := s.client.R().
		SetQueryParamsFromValues(query).
		SetResult(&result).
		SetPathParam("urlBase", s.urlBase).
		SetQueryParam("api_key", s.APIKey).
		Get("https://{urlBase}/geocoding/v1/search/structured")
	if err != nil {
		return nil, fmt.Errorf("structured geocoding get: %w", err)
	}

	if !resp.IsSuccess() {
		return nil, fmt.Errorf("structude geocoding status: %w", err)
	}
	return &result, nil
}

func (sgr StructuredGeocodeRequest) endpoint() string {
	return "/v1/search/structured"
}
func (sgr StructuredGeocodeRequest) toQueryParams() (url.Values, error) {
	return query.Values(sgr)
}
