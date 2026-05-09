package geocode

import (
	"context"
	"fmt"

	"github.com/Gleipnir-Technology/nidus-sync/db"
	//"github.com/Gleipnir-Technology/nidus-sync/db/gen/nidus-sync/public/model"
	querypublic "github.com/Gleipnir-Technology/nidus-sync/db/query/public"
	platformaddress "github.com/Gleipnir-Technology/nidus-sync/platform/address"
	"github.com/Gleipnir-Technology/nidus-sync/platform/types"
	"github.com/Gleipnir-Technology/nidus-sync/stadia"
	//"github.com/rs/zerolog/log"
)

// Ensure the provided address exists. If it doesn't add it to the database.
func EnsureAddress(ctx context.Context, txn db.Ex, a types.Address) (types.Address, error) {
	existing, err := querypublic.AddressFromGID(ctx, txn, a.GID)
	if err != nil {
		return types.Address{}, fmt.Errorf("query address from gid: %w", err)
	}
	if existing != nil {
		return types.AddressFromModel(*existing), nil
	}
	addr, err := platformaddress.InsertAddress(ctx, txn, a)
	if err != nil {
		return types.Address{}, fmt.Errorf("insert address: %w", err)
	}
	return addr, nil
}

func ensureAddressFromFeature(ctx context.Context, txn db.Ex, feature stadia.GeocodeFeature) (types.Address, error) {
	var result types.Address
	if feature.Geometry.Type != "Point" {
		return result, fmt.Errorf("Can't hanlde stadia geometry %s", feature.Geometry.Type)
	}
	existing, err := querypublic.AddressFromGID(ctx, txn, feature.Properties.GID)
	if err != nil {
		return types.Address{}, fmt.Errorf("query address from gid: %w", err)
	}
	if existing != nil {
		return types.AddressFromModel(*existing), nil
	}

	return platformaddress.InsertAddressFeature(ctx, txn, feature)
}
