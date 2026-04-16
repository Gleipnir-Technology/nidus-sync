package types

type Feature struct {
	ID       int32    `db:"id" json:"id"`
	Location Location `db:"location" json:"location"`
	Type     string   `db:"-" json:"type"`
}
