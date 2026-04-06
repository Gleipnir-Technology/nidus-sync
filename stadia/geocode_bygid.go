package stadia

import (
	"context"
	"fmt"

	"github.com/google/go-querystring/query"
)

type RequestGeocodeByGID struct {
	GIDs []string `url:"ids,comma"`

	// Other parameters
	Lang *string `url:"lang,omitempty" json:"lang,omitempty"`
}

func (s *StadiaMaps) GeocodeByGID(ctx context.Context, req RequestGeocodeByGID) (*GeocodeResponse, error) {
	// https://docs.stadiamaps.com/geocoding-search-autocomplete/place-details/
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
		Get("https://{urlBase}/geocoding/v2/place_details")
	if err != nil {
		return nil, fmt.Errorf("autocomplete get: %w", err)
	}

	if !resp.IsSuccess() {
		return nil, parseError(resp)
	}
	return &result, nil
}
