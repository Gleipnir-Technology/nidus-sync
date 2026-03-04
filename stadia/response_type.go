package stadia

// GeocodeResponse represents the top-level response from the geocoding API
type GeocodeResponse struct {
	BBox         []float64        `json:"bbox"` // [W, S, E, N]
	ErrorMessage string           `json:"error,omitempty"`
	Features     []GeocodeFeature `json:"features"`
	Geocode      GeocodeMeta      `json:"geocoding"`
	Type         string           `json:"type"` // Should be "FeatureCollection"
}

// GeocodeMeta contains metadata about the geocoding request
type GeocodeMeta struct {
	Attribution string                 `json:"attribution"`
	Error       string                 `json:"error,omitempty"`  // v2
	Errors      []string               `json:"errors,omitempty"` // v1
	Query       map[string]interface{} `json:"query,omitempty"`
	Warnings    []string               `json:"warnings,omitempty"`
}

// GeocodeFeature represents a GeoJSON feature in the response
type GeocodeFeature struct {
	Type       string            `json:"type"` // Should be "Feature"
	Geometry   GeocodeGeometry   `json:"geometry"`
	Properties GeocodeProperties `json:"properties"`
}

// GeocodeGeometry represents the GeoJSON geometry
type GeocodeGeometry struct {
	Type        string    `json:"type"` // "Point", "Polygon", etc.
	Coordinates []float64 `json:"coordinates"`
}

// GeocodeProperties contains the properties of a geocoding result
type GeocodeProperties struct {
	Addendum    map[string]interface{} `json:"addendum,omitempty"`
	Accuracy    string                 `json:"accuracy"`     // 'point'
	Confidence  float64                `json:"confidence"`   // 1
	Country     string                 `json:"country"`      // 'United States'
	CountryA    string                 `json:"country_a"`    // 'USA'
	CountryCode string                 `json:"country_code"` // 'US'
	CountryGID  string                 `json:"country_gid"`  // 'whosonfirst:country:85633793'
	County      string                 `json:"county"`       // "Tulare County"
	CountyA     string                 `json:"county_a"`     // 'TL'
	CountyGID   string                 `json:"county_gid"`   // 'whosonfirst:county:102082895'
	GID         string                 `json:"gid"`          // 'openaddresses:address:us/ca/tulare-addresses-county:fe9dfab3d45c4550'
	HouseNumber string                 `json:"housenumber"`  // '1234'
	ID          string                 `json:"id"`           // us/ca/tulare-addresses-county:fe9dfab3d45c4550
	Label       string                 `json:"label"`        // 1234 Main St, Dinuba, CA, USA
	Layer       string                 `json:"layer"`        // 'address'
	Locality    string                 `json:"locality"`     // 'Dinuba'
	LocalityGID string                 `json:"locality_gid"` // 'whosonfirst:locality:85922491'
	MatchType   string                 `json:"match_type"`   // 'exact'
	Name        string                 `json:"name"`         // '1234 Main St'
	PostalCode  string                 `json:"postalcode"`   // '93618'
	Region      string                 `json:"region"`       // 'California'
	RegionA     string                 `json:"region_a"`     // 'CA'
	RegionGID   string                 `json:"region_gid"`   // 'whosonfirst:region:85688637'
	Source      string                 `json:"source"`       // 'openaddresses'
	SourceID    string                 `json:"source"`       // 'us/ca/tulare-addresses-county:fe9dfab3d45c4550'
	Street      string                 `json:"street"`       // 'Main Street'
}

// GeocodeSource represents a source of geocoding data
type GeocodeSource struct {
	Source   string `json:"source"`
	SourceID string `json:"source_id"`
}
