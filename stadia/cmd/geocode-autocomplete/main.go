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
	query := flag.String("query", "", "Street address query to autocomplete")
	boundaryRectMaxLat := flag.Float64("boundary-rect-max-lat", 0, "The max lat of the boundary")
	boundaryRectMinLat := flag.Float64("boundary-rect-min-lat", 0, "The min lat of the boundary")
	boundaryRectMaxLon := flag.Float64("boundary-rect-max-lng", 0, "The max lon of the boundary")
	boundaryRectMinLon := flag.Float64("boundary-rect-min-lng", 0, "The min lon of the boundary")
	focusLat := flag.Float64("focus-lat", 0, "The latitude of the focus point")
	focusLng := flag.Float64("focus-lng", 0, "The longitude of the focus point")

	// Parse the flags
	flag.Parse()

	// Validate required arguments
	if *query == "" {
		log.Println("Error: -query is required")
		flag.Usage()
		os.Exit(1)
	}

	if *focusLat != 0 && *focusLng == 0 {
		log.Println("Error: you must specify both focus-lat and focus-lng together, not just focus-lat")
		flag.Usage()
		os.Exit(1)
	}
	if *focusLat == 0 && *focusLng != 0 {
		log.Println("Error: you must specify both focus-lat and focus-lng together, not just focus-lng")
		flag.Usage()
		os.Exit(1)
	}
	if (*boundaryRectMaxLat != 0 ||
		*boundaryRectMinLat != 0 ||
		*boundaryRectMaxLon != 0 ||
		*boundaryRectMinLon != 0) && (*boundaryRectMaxLat == 0 ||
		*boundaryRectMinLat == 0 ||
		*boundaryRectMaxLon == 0 ||
		*boundaryRectMinLon == 0) {
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
	req := stadia.RequestGeocodeAutocomplete{
		Text: *query,
	}
	if *focusLat != 0 && *focusLng != 0 {
		req.FocusPointLat = focusLat
		req.FocusPointLng = focusLng
	}
	if *boundaryRectMaxLat != 0 {
		req.BoundaryRectMaxLat = boundaryRectMaxLat
		req.BoundaryRectMinLat = boundaryRectMinLat
		req.BoundaryRectMaxLon = boundaryRectMaxLon
		req.BoundaryRectMinLon = boundaryRectMinLon
	}
	resp, err := client.GeocodeAutocomplete(ctx, req)
	if err != nil {
		log.Printf("err: %v\n", err)
		os.Exit(2)
	}
	log.Printf("type: %s, features: %d\n", resp.Type, len(resp.Features))
	for i, feature := range resp.Features {
		log.Printf("feature %d: type %s\n", i, feature.Type)
		if feature.Geometry == nil {
			log.Printf("\tno geometry")
		} else {
			log.Printf("\tgeometry %s\n", feature.Geometry.Type) //, feature.Geometry.Coordinates[0], feature.Geometry.Coordinates[1])
		}
		log.Printf("\tproperties %s\n", feature.Properties.Layer)
		switch feature.Properties.Layer {
		case "address":
			log.Printf("\t\t%s", feature.Properties.Name)
			if feature.Properties.CoarseLocation != nil {
				log.Printf("\t\t%s", *feature.Properties.CoarseLocation)
			}
			log.Printf("\t\t%s", feature.Properties.Precision)
			log.Printf("\t\t%s", feature.Properties.Layer)
			log.Printf("\t\t%s", feature.Properties.GID)
		}
	}
}
