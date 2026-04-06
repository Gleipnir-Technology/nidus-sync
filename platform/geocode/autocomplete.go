package geocode

import (
	"context"
	"fmt"

	"github.com/Gleipnir-Technology/nidus-sync/db/models"
	"github.com/Gleipnir-Technology/nidus-sync/stadia"
	"github.com/rs/zerolog/log"
)

type AutocompleteResult struct {
	Detail   string
	GID      string
	Layer    string // 'poi', 'postalcode', 'street',
	Locality string
}

func Autocomplete(ctx context.Context, org *models.Organization, address string) ([]*AutocompleteResult, error) {
	req := stadia.RequestGeocodeAutocomplete{
		Text: address,
	}
	maybeAddServiceArea(&req, org)
	resp, err := client.GeocodeAutocomplete(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("client raw geocode failure on %s: %w", address, err)
	}
	result := make([]*AutocompleteResult, len(resp.Features))
	for i, r := range resp.Features {
		if r.Type != "Feature" {
			log.Error().Str("type", r.Type).Msg("should be handled from Stadia")
			continue
		}
		var locality string
		if r.Properties.CoarseLocation != nil {
			locality = *r.Properties.CoarseLocation
		} else {
			locality = "???"
		}
		result[i] = &AutocompleteResult{
			Detail:   r.Properties.Name,
			GID:      r.Properties.GID,
			Layer:    r.Properties.Layer,
			Locality: locality,
		}
	}
	return result, nil
}
