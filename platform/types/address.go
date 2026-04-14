package types

import (
	"fmt"

	"github.com/Gleipnir-Technology/nidus-sync/db/models"
	//"github.com/rs/zerolog/log"
)

type Address struct {
	Country    string    `db:"country" json:"country"`
	GID        string    `db:"gid" json:"gid" schema:"gid"`
	ID         *int32    `db:"id" json:"-" schema:"-"`
	Locality   string    `db:"locality" json:"locality"`
	Location   *Location `db:"location" json:"location" schema:"location"`
	Number     string    `db:"number" json:"number"`
	PostalCode string    `db:"postal_code" json:"postal_code"`
	Raw        string    `db:"raw" json:"raw" schema:"raw"`
	Region     string    `db:"region" json:"region"`
	Street     string    `db:"street" json:"street"`
	Unit       string    `db:"unit" json:"unit"`
}

func (a Address) String() string {
	return fmt.Sprintf("%s %s, %s, %s, %s, %s", a.Number, a.Street, a.Locality, a.Region, a.PostalCode, a.Country)
}
func AddressFromModel(m *models.Address) Address {
	//log.Debug().Int32("id", m.ID).Float64("lat", m.LocationLatitude.GetOr(0.0)).Float64("lng", m.LocationLongitude.GetOr(0.0)).Msg("converting address")
	return Address{
		Country:  m.Country,
		GID:      m.Gid,
		ID:       &m.ID,
		Locality: m.Locality,
		Location: &Location{
			Latitude:  m.LocationLatitude.GetOr(0.0),
			Longitude: m.LocationLongitude.GetOr(0.0),
		},
		Number:     m.Number,
		PostalCode: m.PostalCode,
		Raw:        "",
		Region:     m.Region,
		Street:     m.Street,
		Unit:       m.Unit,
	}
}
