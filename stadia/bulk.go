package stadia

import (
	"fmt"
	"io"
)

type BulkGeocodeQuery interface {
	endpoint() string
}

// BulkGeocodeRequestItem represents a single request in a bulk geocoding operation
type BulkGeocodeRequestItem struct {
	Endpoint string           `json:"endpoint"`
	Query    BulkGeocodeQuery `json:"query"`
}

// BulkGeocodeResponseItem represents a single response in a bulk geocoding operation
type BulkGeocodeResponseItem struct {
	Response *GeocodeResponse `json:"response,omitempty"`
	Status   int              `json:"status"`
	Message  string           `json:"msg,omitempty"`
}

func (s *StadiaMaps) BulkGeocode(requests []BulkGeocodeQuery) ([]BulkGeocodeResponseItem, error) {
	// https://docs.stadiamaps.com/geocoding-search-autocomplete/bulk-geocoding-search/
	// POST 'https://api.stadiamaps.com/geocoding/v1/search/bulk?api_key=YOUR-API-KEY'
	body := make([]BulkGeocodeRequestItem, 0)
	for _, r := range requests {
		body = append(body, BulkGeocodeRequestItem{
			Endpoint: r.endpoint(),
			Query:    r,
		})
	}
	var results []BulkGeocodeResponseItem
	var api_error Error
	resp, err := s.client.R().
		SetBody(body).
		SetContentType("application/json").
		SetPathParam("urlBase", s.urlBase).
		SetQueryParam("api_key", s.APIKey).
		SetError(&api_error).
		SetResult(&results).
		Post("https://{urlBase}/geocoding/v1/search/bulk")

	if err != nil {
		return nil, fmt.Errorf("bulk geocode request: %w", err)
	}

	if !resp.IsSuccess() {
		if api_error.Error() != "" {
			return nil, &api_error
		}
		content, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("read all failure: %w", err)
		}
		return nil, fmt.Errorf("bulk geocoding request failed with status code: %d: %s", resp.StatusCode(), content)
	}

	return results, nil
}
