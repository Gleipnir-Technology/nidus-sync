package h3utils

import (
	"fmt"
	"strings"

	"github.com/Gleipnir-Technology/go-geojson2h3/v2"
	"github.com/tidwall/geojson"
	"github.com/uber/h3-go/v4"
)

/*
func h3ToBoundsGeoJSON(c h3.Cell) (string, error) {
	b, err := h3.CellToBoundary(c)
	if err != nil {
		respondError(w, "Failed to get cell boundary", err, http.StatusInternalServerError)
		return
	}
	features, err := geojson2h3.ToFeatureCollection(b)
	if err != nil {
		return "", fmt.Errorf("Failed to convert boundary to
}
*/

func ToCell(s string) (h3.Cell, error) {
	c := h3.CellFromString(s)
	if !c.IsValid() {
		return c, fmt.Errorf("Invalid cell definition '%s'", s)
	}
	return c, nil
}
func H3ToGeoJSON(indexes []h3.Cell) (interface{}, error) {
	featureCollection, err := geojson2h3.ToFeatureCollection(indexes)
	if err != nil {
		return "", fmt.Errorf("Failed to get feature collection: %w", err)
	}
	return featureCollection.JSON(), nil
}

func main2() {
	resolution := 9
	object, err := geojson.Parse(`{"type":"FeatureCollection","features":[{"type":"Feature","properties":{"shape":"Polygon","name":"Unnamed Layer","category":"default"},"geometry":{"type":"Polygon","coordinates":[[[-73.901303,40.756892],[-73.893924,40.743755],[-73.871476,40.756278],[-73.863378,40.764175],[-73.871444,40.768467],[-73.879852,40.760014],[-73.885515,40.764045],[-73.891522,40.761054],[-73.901303,40.756892]]]},"id":"a6ca1b7e-9ddf-4425-ad07-8a895f7d6ccf"}]}`, nil)
	if err != nil {
		panic(err)
	}

	indexes, err := geojson2h3.ToH3(resolution, object)
	if err != nil {
		panic(err)
	}
	for _, index := range indexes {
		fmt.Printf("h3index: %s\n", index.String())
	}

	featureCollection, err := geojson2h3.ToFeatureCollection(indexes)
	if err != nil {
		panic(err)
	}
	fmt.Println("Polyfill:")
	fmt.Println(featureCollection.JSON())
}

// Given a cell at a smaller resolution remap it to the larger resolution
func scaleCell(cell h3.Cell, resolution int) (h3.Cell, error) {
	r := cell.Resolution()
	if r == resolution {
		return cell, nil
	}
	latLong, err := cell.LatLng()
	if err != nil {
		return 0, fmt.Errorf("Failed to get latlng: %w", err)
	}
	scaled, err := h3.LatLngToCell(latLong, resolution)
	if err != nil {
		return 0, fmt.Errorf("Failed to create latlng: %w", err)
	}
	return scaled, nil
}

func GetCell(x, y float64, resolution int) (h3.Cell, error) {
	latLng := h3.NewLatLng(y, x)
	return h3.LatLngToCell(latLng, resolution)
}

func CellToPostgisGeometry(c h3.Cell) (string, error) {
	boundary, err := h3.CellToBoundary(c)
	if err != nil {
		return "", fmt.Errorf("Failed to get cell boundary: %w", err)
	}
	var sb strings.Builder

	for i, p := range boundary {
		if i > 0 {
			sb.WriteString(",")
		}
		fmt.Fprintf(&sb, "%g %g", p.Lng, p.Lat)
	}
	// add the first point on to the end to close the polygon
	sb.WriteString(",")
	fmt.Fprintf(&sb, "%g %g", boundary[0].Lng, boundary[0].Lat)

	return fmt.Sprintf("POLYGON((%s))", sb.String()), nil
}

// Convert from an accuracy in meters of GPS coordinates to the H3 Resolution that has at least
// the same area. In other words, for a GPS coordinate accuracy of 2m you have pi*(2m)^2 or ~12.5m^2
// of area which corresponds to resolution 13 (average area of 43.87^2) vs resolution 14 (average area 6.26m^2)
// See https://h3geo.org/docs/core-library/restable
func MeterAccuracyToH3Resolution(accuracy_in_meters float64) int {
	area := accuracy_in_meters * accuracy_in_meters * 3.1415
	if area < 0.895 {
		return 15
	} else if area < 6.267 {
		return 14
	} else if area < 43.87 {
		return 13
	} else if area < 307.092 {
		return 12
	} else if area < 2149.643 {
		return 11
	} else if area < 15_047.502 {
		return 10
	} else if area < 105_332.513 {
		return 9
	} else if area < 737_327.598 {
		return 8
	} else if area < 5_161_293.360 {
		return 7
	} else if area < 36_129_062.164 {
		return 6
	} else if area < 252_903_858.182 {
		return 5
	} else if area < 1_770_347_654.491 {
		return 4
	} else if area < 12_393_434_655.088 {
		return 3
	} else if area < 86_801_780_398.997 {
		return 2
	} else if area < 609_788_441_794.134 {
		return 1
	} else {
		return 0
	}
}
