package types

import (
	"fmt"
)

type Location struct {
	Latitude  float64 `db:"latitude" json:"latitude"`
	Longitude float64 `db:"longitude" json:"longitude"`
}

func (l Location) String() string {
	return fmt.Sprintf("%f %f", l.Longitude, l.Latitude)
}
