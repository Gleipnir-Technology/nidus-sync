package platform

import (
	"errors"
	"fmt"

	"github.com/Gleipnir-Technology/nidus-sync/db/enums"
	"github.com/Gleipnir-Technology/nidus-sync/h3utils"
	"github.com/rs/zerolog/log"
	"github.com/uber/h3-go/v4"
)

type LatLng struct {
	Latitude      *float64
	Longitude     *float64
	MapZoom       float32
	AccuracyValue float64
	AccuracyType  enums.PublicreportAccuracytype
}

func (l LatLng) Resolution() uint {
	switch l.AccuracyType {
	// These accuracy_type strings come from the Mapbox Geocoding API definition and
	// are far from scientific
	case enums.PublicreportAccuracytypeRooftop:
		return 14
	case enums.PublicreportAccuracytypeParcel:
		return 13
	case enums.PublicreportAccuracytypePoint:
		return 13
	case enums.PublicreportAccuracytypeInterpolated:
		return 12
	case enums.PublicreportAccuracytypeApproximate:
		return 11
	case enums.PublicreportAccuracytypeIntersection:
		return 10
	// This is a special indicator that we got our location from the browser measurements
	case enums.PublicreportAccuracytypeBrowser:
		return uint(h3utils.MeterAccuracyToH3Resolution(l.AccuracyValue))
	default:
		log.Warn().Str("accuracy-type", string(l.AccuracyType)).Msg("unrecognized accuracy type, this indicates either a weird client or misbehaving web page. Defaulting to resolution 13")
		return 13
	}
}
func (l LatLng) H3Cell() (*h3.Cell, error) {
	if l.Longitude == nil || l.Latitude == nil {
		return nil, errors.New("nil lat/lng")
	}
	result, err := h3utils.GetCell(*l.Longitude, *l.Latitude, int(l.Resolution()))
	return &result, err
}
func (l LatLng) GeometryQuery() (string, error) {
	if l.Longitude == nil || l.Latitude == nil {
		return "", errors.New("nil lat/lng")
	}
	return fmt.Sprintf("ST_Point(%f, %f, 4326)", *l.Longitude, *l.Latitude), nil
}
