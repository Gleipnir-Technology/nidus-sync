package geomutil

import (
	"errors"
	"github.com/twpayne/go-geom"
)

func AsPoint(g geom.T) (geom.Point, error) {
	p, ok := g.(*geom.Point)
	if !ok {
		return geom.Point{}, errors.New("not a point")
	}
	if p == nil {
		return geom.Point{}, errors.New("nil point")
	}
	return *p, nil
}
func PointFromLngLat(lng, lat float64) geom.T {
	return geom.NewPointFlat(geom.XY, []float64{lng, lat})
}
