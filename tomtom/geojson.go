package tomtom

import (
	"fmt"
	"strings"
)

// Convert a slice of points to GeoJSON
func PolylineToGeoJSON(polyline []Point) string {
	var sb strings.Builder

	sb.WriteString(`{"type":"LineString","coordinates":[`)

	for i, point := range polyline {
		// IMPORTANT: GeoJSON uses [longitude, latitude] order!
		sb.WriteString(fmt.Sprintf("[%g,%g]", point.Longitude, point.Latitude))

		// Add comma if not the last point
		if i < len(polyline)-1 {
			sb.WriteString(",")
		}
	}

	sb.WriteString("]}")

	return sb.String()
}
