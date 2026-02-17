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
		Params: tomtom.CalculateRouteParams{
			VersionNumber: 1,
			Locations:     "52.50931,13.42936:52.50274,13.43872",
			ContentType:   "json",
			Traffic:       &traffic,
			TravelMode:    "car",
			RouteType:     "fastest",
		},
	}

	routeResp, err := client.CalculateRoute(routeRequest)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Route distance: %d meters\n", routeResp.Routes[0].Summary.LengthInMeters)

	// Example 2: Calculate reachable range
	/*
	   energyBudget := 43.0
	   rangeRequest := &tomtom.CalculateReachableRangeRequest{
	       Params: tomtom.CalculateReachableRangeParams{
	           VersionNumber:    1,
	           Origin:           "52.50931,13.42936",
	           ContentType:      "json",
	           EnergyBudgetInkWh: &energyBudget,
	           VehicleEngineType: "electric",
	       },
	   }

	   rangeResp, err := client.CalculateReachableRange(rangeRequest)
	   if err != nil {
	       log.Fatal(err)
	   }

	   fmt.Printf("Reachable range includes %d points in the polygon\n", len(rangeResp.Polygon.Exterior))
	*/
}
