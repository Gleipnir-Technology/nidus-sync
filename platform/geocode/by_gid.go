package geocode

import (
	"context"
	"fmt"

	"github.com/Gleipnir-Technology/nidus-sync/db"
	"github.com/Gleipnir-Technology/nidus-sync/h3utils"
	"github.com/Gleipnir-Technology/nidus-sync/platform/types"
	"github.com/Gleipnir-Technology/nidus-sync/stadia"
)

func ByGID(ctx context.Context, gid string) (*GeocodeResult, error) {
	req := stadia.RequestGeocodeByGID{
		GIDs: []string{gid},
	}
	resp, err := client.GeocodeByGID(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("geocodebygid: %w", err)
	}
	if len(resp.Features) < 1 {
		return nil, fmt.Errorf("no features in result")
	}
	feature := resp.Features[0]
	location := types.Location{
		Latitude:  feature.Geometry.Coordinates[1],
		Longitude: feature.Geometry.Coordinates[0],
	}
	cell, err := h3utils.GetCell(location.Longitude, location.Latitude, 15)
	if err != nil {
		return nil, fmt.Errorf("latlngtocell: %w", err)
	}
	addr, err := ensureAddressFromFeature(ctx, db.PGInstance.PGXPool, feature)
	if err != nil {
		return nil, fmt.Errorf("insert address: %w", err)
	}
	return &GeocodeResult{
		Address: types.Address{
			Country:    feature.Properties.Context.ISO3166A3,
			GID:        feature.Properties.GID,
			ID:         addr.ID,
			Locality:   feature.Properties.Context.WhosOnFirst.Locality.Name,
			Location:   &location,
			Number:     feature.Properties.AddressComponents.Number,
			PostalCode: feature.Properties.AddressComponents.PostalCode,
			Raw:        feature.Properties.FormattedAddressLine,
			Region:     feature.Properties.Context.WhosOnFirst.Region.Name,
			Street:     feature.Properties.AddressComponents.Street,
			Unit:       "",
		},
		Cell: cell,
	}, nil
}
