package geomutil

import (
	"github.com/twpayne/go-geom"
)

func PointFromLngLat(lng, lat float64) geom.T {
	return geom.NewPointFlat(geom.XY, []float64{lng, lat})
}
