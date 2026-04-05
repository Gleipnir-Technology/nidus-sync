package resource

import (
	"context"
	nhttp "github.com/Gleipnir-Technology/nidus-sync/http"
	"github.com/Gleipnir-Technology/nidus-sync/platform/geocode"
	"net/http"
	//"github.com/rs/zerolog/log"
)

type geocodeR struct {
	router *router
}

type geocodeSuggestion struct {
	Detail   string `json:"detail"`
	Locality string `json:"locality"`
}

func Geocode(r *router) *geocodeR {
	return &geocodeR{
		router: r,
	}
}

func (res *geocodeR) SuggestionList(ctx context.Context, r *http.Request, query QueryParams) ([]*geocodeSuggestion, *nhttp.ErrorWithStatus) {
	if query.Query == nil {
		return nil, nhttp.NewBadRequest("you must include a query")
	}
	completions, err := geocode.Autocomplete(ctx, nil, *query.Query)
	if err != nil {
		return nil, nhttp.NewError("geocode: %w", err)
	}
	result := make([]*geocodeSuggestion, len(completions))
	for i, c := range completions {
		result[i] = &geocodeSuggestion{
			Detail:   c.Detail,
			Locality: c.Locality,
		}
	}
	return result, nil
}
