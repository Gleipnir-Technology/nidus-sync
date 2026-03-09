package types

import (
	"fmt"

	"github.com/Gleipnir-Technology/nidus-sync/db/enums"
)

type Address struct {
	Country    string `db:"country" json:"country"`
	Locality   string `db:"locality" json:"locality"`
	Number     string `db:"number" json:"number"`
	PostalCode string `db:"postal_code" json:"postal_code"`
	Region     string `db:"region" json:"region"`
	Street     string `db:"street" json:"street"`
	Unit       string `db:"unit" json:"unit"`
}

func (a Address) String() string {
	return fmt.Sprintf("%s %s, %s, %s, %s, %s", a.Number, a.Street, a.Locality, a.Region, a.PostalCode, a.Country)
}
func (a Address) CountryEnum() enums.Countrytype {
	return enums.CountrytypeUsa
}
