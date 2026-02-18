package main

import (
	"fmt"
	"log"

	"github.com/Gleipnir-Technology/nidus-sync/tomtom"
)

func main() {
	client := tomtom.NewClient()

	// Example 1: Calculate a simple route
	traffic := false
	routeRequest := &tomtom.CalculateRouteRequest{
		Locations: []tomtom.Point{
			tomtom.P(52.50931, 13.42936),
			tomtom.P(52.50274, 13.43872),
		},
		Traffic:    &traffic,
		TravelMode: tomtom.TravelModeCar,
		RouteType:  tomtom.RouteTypeFastest,
	}

	routeResp, err := client.CalculateRoute(routeRequest)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Route distance: %d meters\n", routeResp.Routes[0].Summary.LengthInMeters)
}
