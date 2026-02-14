package stadia

// GeocodeResponse represents the top-level response from the geocoding API
type GeocodeResponse struct {
	Geocode  GeocodeMeta      `json:"geocoding"`
	Type     string           `json:"type"` // Should be "FeatureCollection"
	BBox     []float64        `json:"bbox"` // [W, S, E, N]
	Features []GeocodeFeature `json:"features"`
}

// GeocodeMeta contains metadata about the geocoding request
type GeocodeMeta struct {
	Attribution string                 `json:"attribution"`
	Query       map[string]interface{} `json:"query,omitempty"`
	Warnings    []string               `json:"warnings,omitempty"`
	Errors      []string               `json:"errors,omitempty"` // v1
	Error       string                 `json:"error,omitempty"`  // v2
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
	GID                   string                 `json:"gid"`
	Layer                 string                 `json:"layer"`
	Sources               []GeocodeSource        `json:"sources"`
	Precision             string                 `json:"precision"`
	Name                  string                 `json:"name"`
	FormattedAddressLines []string               `json:"formatted_address_lines"`
	FormattedAddressLine  string                 `json:"formatted_address_line"`
	CoarseLocation        string                 `json:"coarse_location"`
	AddressComponents     AddressComponents      `json:"address_components,omitempty"`
	Context               GeocodeContext         `json:"context,omitempty"`
	Confidence            float64                `json:"confidence,omitempty"`
	Distance              float64                `json:"distance,omitempty"`
	Addendum              map[string]interface{} `json:"addendum,omitempty"`
}

// GeocodeSource represents a source of geocoding data
type GeocodeSource struct {
	Source   string `json:"source"`
	SourceID string `json:"source_id"`
}

// AddressComponents represents the structured components of an address
type AddressComponents struct {
	Number     string `json:"number,omitempty"`
	Street     string `json:"street,omitempty"`
	Unit       string `json:"unit,omitempty"`
	PostalCode string `json:"postal_code,omitempty"`
}

// GeocodeContext represents the geographic context of a result
type GeocodeContext struct {
	WhosOnFirst WhosOnFirstContext `json:"whosonfirst,omitempty"`
	ISO3166A2   string             `json:"iso_3166_a2,omitempty"`
	ISO3166A3   string             `json:"iso_3166_a3,omitempty"`
}

// WhosOnFirstContext contains geographic hierarchy information
type WhosOnFirstContext struct {
	Country       *ContextPlace `json:"country,omitempty"`
	Region        *ContextPlace `json:"region,omitempty"`
	County        *ContextPlace `json:"county,omitempty"`
	Locality      *ContextPlace `json:"locality,omitempty"`
	Neighbourhood *ContextPlace `json:"neighbourhood,omitempty"`
	Borough       *ContextPlace `json:"borough,omitempty"`
}

// ContextPlace represents a place in the geographic hierarchy
type ContextPlace struct {
	GID  string `json:"gid"`
	Name string `json:"name"`
}
