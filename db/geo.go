package db

import (
)

type GeoBounds struct {
	East float64
	North float64
	South float64
	West float64
}

type GeoQuery struct {
	Bounds GeoBounds
	Limit int
}

func NewGeoBounds() GeoBounds {
	return GeoBounds{
		East: 180,
		North: 180,
		South: -180,
		West: -180,
	}
}

func NewGeoQuery() GeoQuery {
	return GeoQuery{
		Bounds: NewGeoBounds(),
		Limit: 0,
	}
}
