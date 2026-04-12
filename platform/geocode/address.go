package geocode

import (
	"context"
	"fmt"
	"time"

	"github.com/Gleipnir-Technology/bob"
	"github.com/Gleipnir-Technology/bob/dialect/psql"
	"github.com/Gleipnir-Technology/bob/dialect/psql/dialect"
	"github.com/Gleipnir-Technology/bob/dialect/psql/im"
	//bobtypes "github.com/Gleipnir-Technology/bob/types"
	"github.com/Gleipnir-Technology/nidus-sync/db/models"
	"github.com/Gleipnir-Technology/nidus-sync/h3utils"
	"github.com/Gleipnir-Technology/nidus-sync/platform/types"
	"github.com/Gleipnir-Technology/nidus-sync/stadia"
	"github.com/rs/zerolog/log"
	"github.com/stephenafamo/scan"
)

type _rowWithID struct {
	ID int32 `db:"id"`
}

// Ensure the provided address exists. If it doesn't add it to the database.
func EnsureAddress(ctx context.Context, txn bob.Executor, a types.Address) (*models.Address, error) {
	address, err := models.Addresses.Query(
		models.SelectWhere.Addresses.Gid.EQ(a.GID),
	).One(ctx, txn)
	if err == nil {
		return address, nil
	}
	id, err := insertAddress(ctx, txn, a)
	if err != nil {
		return nil, fmt.Errorf("insert address: %w", err)
	}
	return &models.Address{
		Country:    a.Country,
		Created:    time.Now(),
		Gid:        a.GID,
		H3cell:     "",
		ID:         *id,
		Locality:   a.Locality,
		Location:   "",
		PostalCode: a.PostalCode,
		Street:     a.Street,
		Unit:       a.Unit,
		Region:     a.Region,
		Number:     a.Number,
	}, nil
}

func ensureAddressFromFeature(ctx context.Context, txn bob.Executor, feature stadia.GeocodeFeature) (int32, error) {
	if feature.Geometry.Type != "Point" {
		return 0, fmt.Errorf("Can't hanlde stadia geometry %s", feature.Geometry.Type)
	}
	lat := feature.Geometry.Coordinates[1]
	lng := feature.Geometry.Coordinates[0]
	cell, err := h3utils.GetCell(lng, lat, 15)
	if err != nil {
		return 0, fmt.Errorf("failed to convert lat %f lng %f to h3 cell", lat, lng)
	}
	query := addressQuery()
	query.Apply(
		im.Values(
			psql.Arg(feature.CountryCode()),
			psql.Arg(time.Now()),
			psql.Arg(feature.Properties.GID),
			psql.Arg(cell.String()),
			psql.Raw("DEFAULT"),
			psql.Arg(feature.Locality()),
			psql.F("ST_Point", lng, lat, 4326),
			psql.Arg(feature.Number()),
			psql.Arg(feature.PostalCode()),
			psql.Arg(feature.Region()),
			psql.Arg(feature.Street()),
			psql.Raw("''"),
		),
	)
	row, err := bob.One(ctx, txn, query, scan.StructMapper[_rowWithID]())
	log.Info().Int32("id", row.ID).Msg("inserted address")
	if err != nil {
		return 0, fmt.Errorf("insert: %w", err)
	}
	return row.ID, nil
}
func insertAddress(ctx context.Context, txn bob.Executor, address types.Address) (*int32, error) {
	lng := address.Location.Longitude
	lat := address.Location.Latitude
	cell, err := h3utils.GetCell(lng, lat, 15)
	if err != nil {
		return nil, fmt.Errorf("failed to convert lat %f lng %f to h3 cell", lat, lng)
	}
	query := addressQuery()
	query.Apply(
		im.Values(
			psql.Arg(address.Country),
			psql.Arg(time.Now()),
			psql.Arg(address.GID),
			psql.Arg(cell),
			psql.Raw("DEFAULT"),
			psql.Arg(address.Locality),
			psql.F("ST_Point", address.Location.Longitude, address.Location.Latitude, 4326),
			psql.Arg(address.Number),
			psql.Arg(address.PostalCode),
			psql.Arg(address.Region),
			psql.Arg(address.Street),
			psql.Raw("''"),
		),
	)
	row, err := bob.One(ctx, txn, query, scan.StructMapper[_rowWithID]())
	if err != nil {
		return nil, fmt.Errorf("insert: %w", err)
	}
	return &row.ID, nil
}
func insertAddresses(ctx context.Context, txn bob.Executor, features []stadia.GeocodeFeature) ([]int32, error) {
	query := addressQuery()
	for _, feature := range features {
		lng := feature.Geometry.Coordinates[0]
		lat := feature.Geometry.Coordinates[1]
		cell, err := h3utils.GetCell(lng, lat, 15)
		if err != nil {
			return nil, fmt.Errorf("failed to convert lat %f lng %f to h3 cell", lat, lng)
		}
		query.Apply(
			im.Values(
				psql.Arg(feature.CountryCode()),
				psql.Arg(time.Now()),
				psql.Arg(feature.Properties.GID),
				psql.Arg(cell.String()),
				psql.Raw("DEFAULT"),
				psql.Arg(feature.Locality()),
				psql.F("ST_Point", lng, lat, 4326),
				psql.Arg(feature.Number()),
				psql.Arg(feature.PostalCode()),
				psql.Arg(feature.Region()),
				psql.Arg(feature.Street()),
				psql.Raw("''"),
			),
		)
	}
	rows, err := bob.All(ctx, txn, query, scan.StructMapper[_rowWithID]())
	if err != nil {
		return nil, fmt.Errorf("insert: %w", err)
	}
	results := make([]int32, len(rows))
	for _, row := range rows {
		results = append(results, row.ID)
	}
	return results, nil
}
func addressQuery() bob.BaseQuery[*dialect.InsertQuery] {
	return psql.Insert(
		im.Into("address", "country", "created", "gid", "h3cell", "id", "locality", "location", "number_", "postal_code", "region", "street", "unit"),
		im.Returning("id"),
	)
}
