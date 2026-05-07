package types

import (
	"context"
	"fmt"

	"github.com/Gleipnir-Technology/bob"
	"github.com/Gleipnir-Technology/bob/dialect/psql"
	"github.com/Gleipnir-Technology/bob/dialect/psql/sm"
	"github.com/Gleipnir-Technology/nidus-sync/db"

	"github.com/Gleipnir-Technology/nidus-sync/db/gen/nidus-sync/public/model"
	"github.com/rs/zerolog/log"
	"github.com/stephenafamo/scan"
)

type Address struct {
	Country    string    `db:"country" json:"country"`
	GID        string    `db:"gid" json:"gid" schema:"gid"`
	ID         *int32    `db:"id" json:"-" schema:"-"`
	Locality   string    `db:"locality" json:"locality"`
	Location   *Location `db:"location" json:"location" schema:"location"`
	Number     string    `db:"number_" json:"number"`
	PostalCode string    `db:"postal_code" json:"postal_code"`
	Raw        string    `db:"raw" json:"raw" schema:"raw"`
	Region     string    `db:"region" json:"region"`
	Street     string    `db:"street" json:"street"`
	Unit       string    `db:"unit" json:"unit"`
}

func (a Address) String() string {
	return fmt.Sprintf("%s %s, %s, %s, %s, %s", a.Number, a.Street, a.Locality, a.Region, a.PostalCode, a.Country)
}
func AddressFromModel(m model.Address) Address {
	//log.Debug().Int32("id", m.ID).Float64("lat", m.LocationLatitude.GetOr(0.0)).Float64("lng", m.LocationLongitude.GetOr(0.0)).Msg("converting address")
	l, err := LocationFromGeom(m.Location)
	if err != nil {
		log.Error().Err(err).Int32("id", m.ID).Msg("getting location for address")
	}
	return Address{
		Country:    m.Country,
		GID:        m.Gid,
		ID:         &m.ID,
		Locality:   m.Locality,
		Location:   &l,
		Number:     m.Number,
		PostalCode: m.PostalCode,
		Raw:        addressToRaw(m),
		Region:     m.Region,
		Street:     m.Street,
		Unit:       m.Unit,
	}
}
func AddressList(ctx context.Context, ids []int32) (map[int32]*Address, error) {
	rows, err := bob.All(ctx, db.PGInstance.BobDB, psql.Select(
		sm.Columns(
			"COALESCE(address.country, 'usa') AS \"country\"",
			"COALESCE(address.gid, '') AS \"gid\"",
			"address.id AS \"id\"",
			"COALESCE(address.locality, '') AS \"locality\"",
			"COALESCE(address.number_, '') AS \"number_\"",
			"COALESCE(address.postal_code, '') AS \"postal_code\"",
			"COALESCE(address.region, '') AS \"region\"",
			"COALESCE(address.street, '') AS \"street\"",
			"COALESCE(address.unit, '') AS \"unit\"",
			// This will work great, up until we add polygons to signal
			"COALESCE(address.location_latitude, 0) AS \"location.latitude\"",
			"COALESCE(address.location_longitude, 0) AS \"location.longitude\"",
		),
		sm.From("address"),
		sm.Where(psql.Quote("address", "id").EQ(psql.Any(ids))),
	), scan.StructMapper[*Address]())
	if err != nil {
		return nil, fmt.Errorf("query addresses: %w", err)
	}
	addresses_by_id := make(map[int32]*Address, len(rows))
	for _, a := range rows {
		addresses_by_id[*a.ID] = a
	}

	return addresses_by_id, err
}
func AddressToRaw(a Address) string {
	return fmt.Sprintf("%s %s, %s, %s", a.Number, a.Street, a.Locality, a.Region)
}
func addressToRaw(m model.Address) string {
	return fmt.Sprintf("%s %s, %s, %s", m.Number, m.Street, m.Locality, m.Region)
}
