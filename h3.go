package main

import (
	"fmt"

	"github.com/Gleipnir-Technology/go-geojson2h3"
	"github.com/tidwall/geojson"
	"github.com/uber/h3-go/v4"
)

func h3Indexes() []h3.Cell {
	//[] uint64{0x852a134ffffffff})
	/*result := make([]h3.H3Index, 0)
	for _, v := range values {
		result = append(result, v)
	}
	return result*/
	return []h3.Cell{
		0x8629aab2fffffff,
		0x8629a8627ffffff,
		0x8629a8607ffffff,
		0x8629a8717ffffff,
		0x8629a8617ffffff,
		0x8629a8407ffffff,
		0x8629a871fffffff,
		0x8629a845fffffff,
		0x8629aab27ffffff,
		0x8629a84e7ffffff,
	}
}
func h3ToGeoJSON(indexes []h3.Cell) (string, error) {
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

func getCell(x, y float64, resolution int) (h3.Cell, error) {
	latLng := h3.NewLatLng(y, x)
	return h3.LatLngToCell(latLng, resolution)
}
