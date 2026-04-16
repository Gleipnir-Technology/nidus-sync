package stadia

import (
	"context"
	"fmt"
	"math"
	"strconv"

	"github.com/rs/zerolog/log"
)

type RequestTileRaster struct {
	Latitude  float64
	Longitude float64
	//Style string
	Zoom uint
}

func (s *StadiaMaps) TileRaster(ctx context.Context, req RequestTileRaster) ([]byte, error) {
	// https://docs.stadiamaps.com/raster/
	//url := "https://{urlBase}/tiles/{style}/{z}/{x}/{y}{r}.png"
	//url := "https://{urlBase}/data/imagery/{z}/{x}/{y}{r}.png"
	url := "https://{urlBase}/tiles/alidade_satellite/{z}/{x}/{y}.jpg"

	y, x := LatLngToTile(req.Zoom, req.Latitude, req.Longitude)
	//var api_error Error
	resp, err := s.client.R().
		SetContext(ctx).
		//SetPathParam("style", req.Style).
		//SetPathParam("r", "").
		SetPathParam("x", strconv.Itoa(int(x))).
		SetPathParam("y", strconv.Itoa(int(y))).
		SetPathParam("z", strconv.Itoa(int(req.Zoom))).
		SetPathParam("urlBase", s.urlBaseTiles).
		SetQueryParam("api_key", s.APIKey).
		Get(url)
	if err != nil {
		return nil, fmt.Errorf("autocomplete get: %w", err)
	}

	if !resp.IsSuccess() {
		return nil, parseError(resp)
	}
	content_type := resp.Header().Get("Content-Type")
	log.Debug().Str("content_type", content_type).Send()
	return resp.Bytes(), nil
}

// LatLngToTile converts GPS coordinates to ArcGIS tile coordinates
func LatLngToTile(level uint, lat, lng float64) (row, column uint) {
	// Get number of tiles per dimension at this zoom level
	numTiles := math.Pow(2, float64(level))

	// Convert longitude to tile column
	// Range: -180 to 180 degrees maps to 0 to numTiles
	column = uint(math.Floor((lng + 180.0) / 360.0 * numTiles))

	// Convert latitude to tile row using Mercator projection
	// First convert lat to radians
	latRad := lat * math.Pi / 180.0

	// Apply Mercator projection formula
	// This maps latitude from -85.0511 to 85.0511 degrees to 0 to numTiles
	mercatorY := 0.5 - math.Log(math.Tan(latRad)+1/math.Cos(latRad))/(2*math.Pi)
	row = uint(math.Floor(mercatorY * numTiles))

	// Ensure values are within valid range
	if column < 0 {
		column = 0
	} else if column >= uint(numTiles) {
		column = uint(numTiles) - 1
	}

	if row < 0 {
		row = 0
	} else if row >= uint(numTiles) {
		row = uint(numTiles) - 1
	}

	return row, column
}
