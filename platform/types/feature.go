package types

type Feature struct {
	ID       int32    `db:"id" json:"id"`
	Location Location `db:"location" json:"location"`
	SiteID   int32    `db:"site_id" json:"-"`
	Type     string   `db:"-" json:"type"`
}
