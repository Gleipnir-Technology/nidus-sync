package platform

import (
	"context"
	"errors"
	"fmt"

	"github.com/Gleipnir-Technology/nidus-sync/db"
	querypublic "github.com/Gleipnir-Technology/nidus-sync/db/query/public"
	"github.com/Gleipnir-Technology/nidus-sync/platform/types"
)

type Address = types.Address

func AddressFromComplianceReportRequestID(ctx context.Context, public_id string) (*types.Address, error) {
	row, err := querypublic.AddressFromComplianceReportRequestID(ctx, db.PGInstance.PGXPool, public_id)
	if err != nil {
		if errors.Is(err, db.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("query address from compliance report request: %w", err)
	}
	result := types.AddressFromModel(row)
	return &result, nil
}

func AddressLocation(ctx context.Context, address types.Address) (*types.Location, error) {
	address_id := int64(*address.ID)
	addr, err := querypublic.AddressFromID(ctx, db.PGInstance.PGXPool, address_id)
	if err != nil {
		return nil, fmt.Errorf("query address: %w", err)
	}
	l, err := types.LocationFromGeom(addr.Location)
	if err != nil {
		return nil, fmt.Errorf("location from geom: %w", err)
	}
	return &l, nil
}

func AddressInsert(ctx context.Context) (*types.Address, error) {
	return nil, nil
}
