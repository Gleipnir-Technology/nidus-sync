package platform

import (
	"github.com/Gleipnir-Technology/nidus-sync/db/enums"
)

type LatLng struct {
	Latitude      *float64
	Longitude     *float64
	MapZoom       float32
	AccuracyValue float64
	AccuracyType  enums.PublicreportAccuracytype
}
