package resource

import (
	"context"
	"net/http"

	nhttp "github.com/Gleipnir-Technology/nidus-sync/http"
	ngeocode "github.com/Gleipnir-Technology/nidus-sync/platform/geocode"
	"github.com/Gleipnir-Technology/nidus-sync/platform/types"
	"github.com/gorilla/mux"
	//"github.com/rs/zerolog/log"
	"github.com/uber/h3-go/v4"
)

type geocodeR struct {
	router *router
}
type geocode struct {
	Address  types.Address  `json:"address"`
	Cell     h3.Cell        `json:"cell"`
	Location types.Location `json:"location"`
}

func newGeocode(g *ngeocode.GeocodeResult) *geocode {
	return &geocode{
		Address:  g.Address,
		Cell:     g.Cell,
		Location: g.Location,
	}
}

type geocodeSuggestion struct {
	Detail   string `json:"detail"`
	GID      string `json:"gid"`
	Locality string `json:"locality"`
	Type     string `json:"type"`
}

func Geocode(r *router) *geocodeR {
	return &geocodeR{
		router: r,
	}
}

func (res *geocodeR) ByGID(ctx context.Context, r *http.Request, query QueryParams) (*geocode, *nhttp.ErrorWithStatus) {
	vars := mux.Vars(r)
	gid := vars["id"]
	if gid == "" {
		return nil, nhttp.NewBadRequest("no id")
	}
	g, err := ngeocode.ByGID(ctx, gid)
	if err != nil {
		return nil, nhttp.NewError("bygid: %w", err)
	}
	return newGeocode(g), nil
}
func (res *geocodeR) Reverse(ctx context.Context, r *http.Request, location types.Location) (*geocode, *nhttp.ErrorWithStatus) {
	g, err := ngeocode.ReverseGeocode(ctx, location)
	if err != nil {
		return nil, nhttp.NewError("reverse: %w", err)
	}
	return newGeocode(g), nil
}
func (res *geocodeR) SuggestionList(ctx context.Context, r *http.Request, query QueryParams) ([]*geocodeSuggestion, *nhttp.ErrorWithStatus) {
	if query.Query == nil {
		return nil, nhttp.NewBadRequest("you must include a query")
	}
	completions, err := ngeocode.Autocomplete(ctx, nil, *query.Query)
	if err != nil {
		return nil, nhttp.NewError("geocode: %w", err)
	}
	result := make([]*geocodeSuggestion, len(completions))
	for i, c := range completions {
		result[i] = &geocodeSuggestion{
			Detail:   c.Detail,
			GID:      c.GID,
			Locality: c.Locality,
			Type:     c.Layer,
		}
	}
	return result, nil
}
