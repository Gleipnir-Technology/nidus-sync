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
	lat := flag.Float64("lat", 0, "The latitude of the point")
	lng := flag.Float64("lng", 0, "The longitude of the point")

	// Parse the flags
	flag.Parse()

	if *lat == 0 || *lng == 0 {
		log.Println("Error: you must specify both lat and lng")
		flag.Usage()
		os.Exit(1)
	}

	key := os.Getenv("STADIA_MAPS_API_KEY")
	if key == "" {
		log.Println("STADIA_MAPS_API_KEY is empty")
		os.Exit(1)
	}

	client := stadia.NewStadiaMaps(key)
	ctx := context.Background()
	req := stadia.RequestReverseGeocode{
		Latitude:  *lat,
		Longitude: *lng,
	}
	resp, err := client.ReverseGeocode(ctx, req)
	if err != nil {
		log.Printf("err: %v\n", err)
		os.Exit(2)
	}
	log.Printf("type: %s, features: %d\n", resp.Type, len(resp.Features))
	for i, feature := range resp.Features {
		log.Printf("feature %d: type %s\n", i, feature.Type)
		log.Printf("\tgeometry %s (%f %f)\n", feature.Geometry.Type, feature.Geometry.Coordinates[0], feature.Geometry.Coordinates[1])
		log.Printf("\tproperties %s\n", feature.Properties.Layer)
	}
}
