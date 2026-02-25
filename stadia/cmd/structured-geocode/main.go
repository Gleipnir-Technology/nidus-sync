package main

import (
	"context"
	"flag"
	"log"
	"os"

	"github.com/Gleipnir-Technology/nidus-sync/stadia"
)

func main() {
	// Define command-line flags
	address := flag.String("address", "", "Street address to geocode")
	boundaryRectMaxLat := flag.Float64("boundary-rect-max-lat", 0, "The max lat of the boundary")
	boundaryRectMinLat := flag.Float64("boundary-rect-min-lat", 0, "The min lat of the boundary")
	boundaryRectMaxLon := flag.Float64("boundary-rect-max-lng", 0, "The max lon of the boundary")
	boundaryRectMinLon := flag.Float64("boundary-rect-min-lng", 0, "The min lon of the boundary")
	postalCode := flag.String("postal-code", "", "Postal code")
	focusLat := flag.Float64("focus-lat", 0, "The latitude of the focus point")
	focusLng := flag.Float64("focus-lng", 0, "The longitude of the focus point")

	// Parse the flags
	flag.Parse()

	// Validate required arguments
	if *address == "" {
		log.Println("Error: -address is required")
		flag.Usage()
		os.Exit(1)
	}

	if *postalCode == "" {
		log.Println("Error: -postal-code is required")
		flag.Usage()
		os.Exit(1)
	}
	if focusLat != nil && focusLng == nil {
		log.Println("Error: you must specify both focus-lat and focus-lng together, not just focus-lat")
		flag.Usage()
		os.Exit(1)
	}
	if focusLat == nil && focusLng != nil {
		log.Println("Error: you must specify both focus-lat and focus-lng together, not just focus-lng")
		flag.Usage()
		os.Exit(1)
	}
	if (boundaryRectMaxLat != nil ||
		boundaryRectMinLat != nil ||
		boundaryRectMaxLon != nil ||
		boundaryRectMinLon != nil) && (boundaryRectMaxLat == nil ||
		boundaryRectMinLat == nil ||
		boundaryRectMaxLon == nil ||
		boundaryRectMinLon == nil) {
		log.Println("If you specify one of boundary-rect you need to specify them all")
		os.Exit(1)
	}

	key := os.Getenv("STADIA_MAPS_API_KEY")
	if key == "" {
		log.Println("STADIA_MAPS_API_KEY is empty")
		os.Exit(1)
	}

	client := stadia.NewStadiaMaps(key)
	ctx := context.Background()
	req := stadia.StructuredGeocodeRequest{
		Address:    address,
		PostalCode: postalCode,
	}
	if focusLat != nil && focusLng != nil {
		req.FocusPointLat = focusLat
		req.FocusPointLng = focusLng
	}
	if boundaryRectMaxLat != nil {
		req.BoundaryRectMaxLat = boundaryRectMaxLat
		req.BoundaryRectMinLat = boundaryRectMinLat
		req.BoundaryRectMaxLon = boundaryRectMaxLon
		req.BoundaryRectMinLon = boundaryRectMinLon
	}
	resp, err := client.StructuredGeocode(ctx, req)
	if err != nil {
		log.Printf("err: %v\n", err)
		os.Exit(2)
	}
	log.Printf("type: %s, features: %d\n", resp.Type, len(resp.Features))
	for i, feature := range resp.Features {
		log.Printf("feature %d: type %s\n", i, feature.Type)
		log.Printf("\tgeometry %s\n", feature.Geometry.Type)
		log.Printf("\tproperties %s\n", feature.Properties.Layer)
	}
}
