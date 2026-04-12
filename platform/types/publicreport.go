package types

import (
	"time"
)

type PublicReport struct {
	Address    Address    `db:"address" json:"address"`
	Created    time.Time  `db:"created" json:"created"`
	ID         int32      `db:"id" json:"-"`
	Images     []Image    `db:"images" json:"images"`
	Location   *Location  `db:"location" json:"location"`
	Log        []LogEntry `db:"-" json:"log"`
	DistrictID *int32     `db:"organization_id" json:"-"`
	District   *string    `db:"-" json:"district"`
	PublicID   string     `db:"public_id" json:"public_id"`
	Reporter   Contact    `db:"reporter" json:"reporter"`
	Status     string     `db:"status" json:"status"`
	Type       string     `db:"report_type" json:"type"`
	URI        string     `db:"-" json:"uri"`
}
