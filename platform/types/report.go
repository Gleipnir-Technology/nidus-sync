package types

import (
	"time"
)

type Report struct {
	Address    Address    `db:"address" json:"address"`
	AddressRaw string     `db:"address_raw" json:"address_raw"`
	Created    time.Time  `db:"created" json:"created"`
	ID         int32      `db:"id" json:"-"`
	Images     []Image    `db:"images" json:"images"`
	Location   *Location  `db:"location" json:"location"`
	Log        []LogEntry `db:"-" json:"log"`
	Nuisance   *Nuisance  `db:"nuisance" json:"nuisance"`
	DistrictID *int32     `db:"organization_id" json:"-"`
	District   *string    `db:"-" json:"district"`
	PublicID   string     `db:"public_id" json:"public_id"`
	Reporter   Contact    `db:"reporter" json:"reporter"`
	Status     string     `db:"status" json:"status"`
	Type       string     `db:"report_type" json:"type"`
	URI        string     `db:"-" json:"uri"`
	Water      *Water     `db:"water" json:"water"`
}
