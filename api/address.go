package api

type Address struct {
	Country    string `db:"country" json:"country"`
	Locality   string `db:"locality" json:"locality"`
	Number     string `db:"number" json:"number"`
	PostalCode string `db:"postal_code" json:"postal_code"`
	Region     string `db:"region" json:"region"`
	Street     string `db:"street" json:"street"`
	Unit       string `db:"unit" json:"unit"`
}
