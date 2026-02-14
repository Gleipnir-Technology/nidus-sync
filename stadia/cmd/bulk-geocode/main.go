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
	requests := make([]stadia.BulkGeocodeQuery, 0)
	requests = append(requests, stadia.StructuredGeocodeRequest{
		Address:    strPtr("12932 Ave 404"),
		PostalCode: strPtr("93615"),
	})
	requests = append(requests, stadia.StructuredGeocodeRequest{
		Address:    strPtr("1187 N Arno Rd"),
		PostalCode: strPtr("93618"),
	})
	resp, err := client.BulkGeocode(requests)
	if err != nil {
		log.Printf("err: %v\n", err)
		os.Exit(2)
	}
	for _, r := range resp {
		log.Printf("Status: %s", r.Status)
	}
}

func strPtr(s string) *string {
	return &s
}
