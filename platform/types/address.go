package types

import (
	"fmt"

	"github.com/Gleipnir-Technology/nidus-sync/db/enums"
	"github.com/Gleipnir-Technology/nidus-sync/db/models"
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
func (a Address) CountryEnum() enums.Countrytype {
	return enums.CountrytypeUsa
}
func AddressFromModel(m *models.Address) Address {
	return Address{
		Country:  m.Country.String(),
		GID:      m.Gid,
		ID:       &m.ID,
		Locality: m.Locality,
		Location: &Location{
			Latitude:  m.LocationY.GetOr(0.0),
			Longitude: m.LocationX.GetOr(0.0),
		},
		Number:     m.Number,
		PostalCode: m.PostalCode,
		Raw:        "",
		Region:     m.Region,
		Street:     m.Street,
		Unit:       m.Unit,
	}
}
