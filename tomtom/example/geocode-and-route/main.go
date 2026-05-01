package main

import (
	"log"

	"github.com/Gleipnir-Technology/nidus-sync/tomtom"
)

func main() {
	client := tomtom.NewClient()

	// Example 1: Geocode a series of addresses
	waypoints := []string{
		"1737 W Houston Ave, Visalia, CA 93291",
		"1138 W Prescott Ave, Visalia, CA 93291",
		"3228 W Clinton Ct, Visalia, CA 93291",
		"3800 N Mendonca St, Visalia, CA 93291",
	}
	coords := make([]tomtom.Point, 0)
	for _, wp := range waypoints {
		geocode, err := client.Geocode(wp)
		if err != nil {
			log.Fatal("Failed to geocode '%s': %w", wp, err)
		}
		if len(geocode.Results) == 0 {
			log.Fatal("Failed to get any results for '%s'", wp)
		}
		result := geocode.Results[0]
		coords = append(coords, result.Position.AsPoint())
		log.Printf("Geocoded %s to %f, %f", wp, result.Position.Longitude, result.Position.Latitude)
	}
	// Example 2: Calculate a simple route through them
	traffic := false
	routeRequest := &tomtom.CalculateRouteRequest{
		Locations:  coords,
		Traffic:    &traffic,
		TravelMode: tomtom.TravelModeCar,
		RouteType:  tomtom.RouteTypeFastest,
	}

	//client.SetDebug(true)
	routeResp, err := client.CalculateRoute(routeRequest)
	if err != nil {
		log.Fatal(err)
	}

	all_points := make([]tomtom.Point, 0)
	all_stops := make([]tomtom.Point, 0)
	for i, route := range routeResp.Routes {
		s := route.Summary
		log.Printf("%d: %d meters, %d seconds, %s traffic delay", i, s.LengthInMeters, s.TravelTimeInSeconds, s.TrafficDelayInSeconds)
		for _, leg := range route.Legs {
			all_stops = append(all_stops, leg.Points[0])
			all_points = append(all_points, leg.Points...)
		}
	}
	lines := tomtom.PolylineToGeoJSON(all_points)
	log.Printf("%s", lines)
	stops := tomtom.PolylineToGeoJSON(all_stops)
	log.Printf("%s", stops)
}
