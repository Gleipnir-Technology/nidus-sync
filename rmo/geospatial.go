package rmo

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/Gleipnir-Technology/nidus-sync/h3utils"
	"github.com/rs/zerolog/log"
	"github.com/uber/h3-go/v4"
)

type GeospatialData struct {
	Cell          h3.Cell
	GeometryQuery string
	Populated     bool
}

func geospatialFromForm(r *http.Request) (GeospatialData, error) {
	lat := r.FormValue("latitude")
	lng := r.FormValue("longitude")
	accuracy_type := r.FormValue("latlng-accuracy-type")
	accuracy_value := r.FormValue("latlng-accuracy-value")

	if lat == "" || lng == "" {
		return GeospatialData{Populated: false}, nil
	}
	latitude, err := strconv.ParseFloat(lat, 64)
	if err != nil {
		return GeospatialData{Populated: false}, fmt.Errorf("Failed to create parse latitude: %v", err)
	}
	longitude, err := strconv.ParseFloat(lng, 64)
	if err != nil {
		return GeospatialData{Populated: false}, fmt.Errorf("Failed to create parse longitude: %v", err)
	}
	var resolution int
	switch accuracy_type {
	// These accuracy_type strings come from the Mapbox Geocoding API definition and
	// are far from scientific
	case "rooftop":
		resolution = 14
	case "parcel":
		resolution = 13
	case "point":
		resolution = 13
	case "interpolated":
		resolution = 12
	case "approximate":
		resolution = 11
	case "intersection":
		resolution = 10
	// This is a special indicator that we got our location from the browser measurements
	case "meters":
	case "browser":
		accuracy_in_meters, err := strconv.ParseFloat(accuracy_value, 64)
		if err != nil {
			return GeospatialData{Populated: false}, fmt.Errorf("Failed to parse '%s' as an accuracy in meters: %v", accuracy_value, err)
		}
		resolution = h3utils.MeterAccuracyToH3Resolution(accuracy_in_meters)
	default:
		log.Warn().Str("accuracy-type", accuracy_type).Msg("unrecognized accuracy type, this indicates either a weird client or misbehaving web page. Defaulting to resolution 13")
		resolution = 13
	}
	cell, err := h3utils.GetCell(longitude, latitude, resolution)
	return GeospatialData{
		Cell:          cell,
		GeometryQuery: fmt.Sprintf("ST_GeometryFromText('Point(%f %f)')", longitude, latitude),
		Populated:     true,
	}, nil
}
