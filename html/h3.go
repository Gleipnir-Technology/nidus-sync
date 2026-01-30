package html

import (
	"fmt"
	"strings"

	"github.com/uber/h3-go/v4"
)

func gisStatement(cb h3.CellBoundary) string {
	var content strings.Builder
	for i, p := range cb {
		if i != 0 {
			content.WriteString(", ")
		}
		content.WriteString(fmt.Sprintf("%f %f", p.Lng, p.Lat))
	}
	// Repeat the first coordinate to close the polygon
	content.WriteString(fmt.Sprintf(", %f %f", cb[0].Lng, cb[0].Lat))
	return fmt.Sprintf("ST_GeomFromText('POLYGON((%s))', 3857)", content.String())
}

func latLngDisplay(ll h3.LatLng) string {
	latDir := "N"
	latVal := ll.Lat
	if ll.Lat < 0 {
		latDir = "S"
		latVal = -ll.Lat
	}

	lngDir := "E"
	lngVal := ll.Lng
	if ll.Lng < 0 {
		lngDir = "W"
		lngVal = -ll.Lng
	}

	return fmt.Sprintf("%.4f° %s, %.4f° %s", latVal, latDir, lngVal, lngDir)
}
