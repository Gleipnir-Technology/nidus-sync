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
	lat := flag.Float64("lat", 0, "The latitude of the tile")
	lng := flag.Float64("lng", 0, "The longitude of the tile")
	zoom := flag.Uint("zoom", 16, "The zoom level")

	// Parse the flags
	flag.Parse()

	if *lat == 0 {
		log.Println("Error: you must specify -lat")
		flag.Usage()
		os.Exit(1)
	}
	if *lng == 0 {
		log.Println("Error: you must specify -lng")
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
	req := stadia.RequestTileRasterLatLng{
		Latitude:  *lat,
		Longitude: *lng,
		Zoom:      *zoom,
	}
	data, err := client.TileRasterLatLng(ctx, req)
	if err != nil {
		log.Printf("err: %v\n", err)
		os.Exit(2)
	}
	err = os.WriteFile("tile.raw", data, 0666)
	if err != nil {
		log.Printf("err: %v\n", err)
		os.Exit(2)
	}
	log.Printf("wrote tile.raw")
}
