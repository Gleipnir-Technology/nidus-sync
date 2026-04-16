package platform

import (
	"context"

	"github.com/Gleipnir-Technology/bob"
	"github.com/Gleipnir-Technology/bob/dialect/psql"
	"github.com/Gleipnir-Technology/bob/dialect/psql/sm"
	"github.com/Gleipnir-Technology/nidus-sync/db"
	"github.com/Gleipnir-Technology/nidus-sync/platform/types"
	"github.com/stephenafamo/scan"
)

type Address = types.Address

func AddressList(ctx context.Context, ids []int32) ([]*types.Address, error) {
	rows, err := bob.All(ctx, db.PGInstance.BobDB, psql.Select(
		sm.Columns(
			"COALESCE(address.country, 'usa') AS \"country\"",
			"COALESCE(address.gid, '') AS \"gid\"",
			"COALESCE(address.locality, '') AS \"locality\"",
			"COALESCE(address.number_, '') AS \"number\"",
			"COALESCE(address.postal_code, '') AS \"postal_code\"",
			"COALESCE(address.region, '') AS \"region\"",
			"COALESCE(address.street, '') AS \"street\"",
			"COALESCE(address.unit, '') AS \"unit\"",
			// This will work great, up until we add polygons to signal
			"COALESCE(ST_Y(address.location_latitude), 0) AS \"location.latitude\"",
			"COALESCE(ST_X(address.location_longitude), 0) AS \"location.longitude\"",
		),
		sm.From("address"),
		sm.Where(psql.Quote("address", "id").EQ(psql.Arg(ids))),
	), scan.StructMapper[*types.Address]())
	return rows, err
}
