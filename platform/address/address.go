package address

import (
	"context"
	"fmt"
	"time"

	"github.com/Gleipnir-Technology/nidus-sync/db"
	"github.com/Gleipnir-Technology/nidus-sync/db/gen/nidus-sync/public/model"
	querypublic "github.com/Gleipnir-Technology/nidus-sync/db/query/public"
	"github.com/Gleipnir-Technology/nidus-sync/h3utils"
	"github.com/Gleipnir-Technology/nidus-sync/platform/types"
	"github.com/Gleipnir-Technology/nidus-sync/stadia"
	"github.com/rs/zerolog/log"
	"github.com/twpayne/go-geom"
)

func InsertAddress(ctx context.Context, txn db.Ex, address types.Address) (types.Address, error) {
	lng := address.Location.Longitude
	lat := address.Location.Latitude
	cell, err := h3utils.GetCell(lng, lat, 15)
	if err != nil {
		return types.Address{}, fmt.Errorf("failed to convert lat %f lng %f to h3 cell", lat, lng)
	}
	addr := model.Address{
		Country: address.Country,
		Created: time.Now(),
		Gid:     address.GID,
		H3cell:  cell.String(),
		//ID:
		Locality:   address.Locality,
		Location:   address.Location.ToGeom(),
		Number:     address.Number,
		PostalCode: address.PostalCode,
		Region:     address.Region,
		Street:     address.Street,
		Unit:       "",
	}
	m, err := querypublic.AddressInsert(ctx, txn, addr)
	if err != nil {
		return types.Address{}, fmt.Errorf("address insert: %w", err)
	}
	log.Info().Int32("id", m.ID).Msg("inserted address")
	return types.AddressFromModel(m), nil
}
func InsertAddressFeature(ctx context.Context, txn db.Ex, feature stadia.GeocodeFeature) (types.Address, error) {
	m, err := addressModelFromFeature(feature)
	if err != nil {
		return types.Address{}, fmt.Errorf("address from feature: %w", err)
	}
	row, err := querypublic.AddressInsert(ctx, txn, m)
	if err != nil {
		return types.Address{}, fmt.Errorf("address insert: %w", err)
	}
	return types.AddressFromModel(row), nil
}
func InsertAddresses(ctx context.Context, txn db.Ex, features []stadia.GeocodeFeature) ([]types.Address, error) {
	models := make([]model.Address, len(features))
	for i, feature := range features {
		m, err := addressModelFromFeature(feature)
		if err != nil {
			return nil, fmt.Errorf("address from feature: %w", err)
		}
		models[i] = m
	}
	addresses, err := querypublic.AddressInserts(ctx, txn, models)
	if err != nil {
		return nil, fmt.Errorf("inserts: %w", err)
	}
	results := make([]types.Address, len(addresses))
	for i, address := range addresses {
		results[i] = types.AddressFromModel(address)
	}
	return results, nil
}
func geomFromLngLat(lng, lat float64) geom.T {
	return geom.NewPointFlat(geom.XY, []float64{lng, lat})
}
func addressModelFromFeature(feature stadia.GeocodeFeature) (model.Address, error) {
	lng := feature.Geometry.Coordinates[0]
	lat := feature.Geometry.Coordinates[1]
	cell, err := h3utils.GetCell(lng, lat, 15)
	if err != nil {
		return model.Address{}, fmt.Errorf("failed to convert lat %f lng %f to h3 cell", lat, lng)
	}
	return model.Address{
		Country: feature.CountryCode(),
		Created: time.Now(),
		Gid:     feature.Properties.GID,
		H3cell:  cell.String(),
		//ID:
		Locality:   feature.Locality(),
		Location:   geomFromLngLat(lng, lat),
		Number:     feature.Number(),
		PostalCode: feature.PostalCode(),
		Region:     feature.Region(),
		Street:     feature.Street(),
		Unit:       "",
	}, nil
}
