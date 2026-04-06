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
	gid := flag.String("gid", "", "The GID to query")

	// Parse the flags
	flag.Parse()

	// Validate required arguments
	if *gid == "" {
		log.Println("Error: -gid is required")
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
	req := stadia.RequestGeocodeByGID{
		GIDs: []string{*gid},
	}
	resp, err := client.GeocodeByGID(ctx, req)
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
