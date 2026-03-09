package types

type Location struct {
	Latitude  float64 `db:"latitude"`
	Longitude float64 `db:"longitude"`
}
