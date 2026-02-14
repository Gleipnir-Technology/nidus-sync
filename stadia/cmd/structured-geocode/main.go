package main

import (
	"log"
	"os"

	"github.com/Gleipnir-Technology/nidus-sync/stadia"
)

func main() {
	key := os.Getenv("STADIA_MAPS_API_KEY")
	if key == "" {
		log.Println("stadia maps api key is empty")
		os.Exit(1)
	}
	client := stadia.NewStadiaMaps(key)
	resp, err := client.StructuredGeocode(stadia.StructuredGeocodeRequest{
		Address:    strPtr("12932 Ave 404"),
		PostalCode: strPtr("93615"),
	})
	if err != nil {
		log.Printf("err: %v\n", err)
		os.Exit(2)
	}
	log.Printf("type: %s", resp.Type)
}

func strPtr(s string) *string {
	return &s
}
