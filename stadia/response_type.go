package stadia

/*
	"address_components": {
	    "number": "3397",
	    "postal_code": "84065",
	    "street": "West Chatel Drive"
	},
*/
type AddressComponents struct {
	Number     string `json:"number"`
	PostalCode string `json:"postal_code"`
	Street     string `json:"street"`
}
type Country struct {
	Abbreviation string `json:"abbreviation"`
	GID          string `json:"gid"`
	Name         string `json:"name"`
}
type County struct {
	Abbreviation string `json:"abbreviation"`
	GID          string `json:"gid"`
	Name         string `json:"name"`
}
type Locality struct {
	GID  string `json:"gid"`
	Name string `json:"name"`
}
type Region struct {
	Abbreviation string `json:"abbreviation"`
	GID          string `json:"gid"`
	Name         string `json:"name"`
}

/*
	"country": {
	    "abbreviation": "USA",
	    "gid": "whosonfirst:country:85633793",
	    "name": "United States"
	},

	"county": {
	    "abbreviation": "SL",
	    "gid": "whosonfirst:county:102082877",
	    "name": "Salt Lake County"
	},

	"locality": {
	    "gid": "whosonfirst:locality:101728073",
	    "name": "Riverton"
	},

	"region": {
	    "abbreviation": "UT",
	    "gid": "whosonfirst:region:85688567",
	    "name": "Utah"
	}
*/
type ContextWhosOnFirst struct {
	Country  Country  `json:"country"`
	County   County   `json:"county"`
	Locality Locality `json:"locality"`
	Region   Region   `json:"region"`
}

/*
	"context": {
	    "iso_3166_a2": "US",
	    "iso_3166_a3": "USA",
	    "whosonfirst": {...}
	    }
	}
*/
type Context struct {
	ISO3166A2   string             `json:"iso_3166_a2"`
	ISO3166A3   string             `json:"iso_3166_a3"`
	WhosOnFirst ContextWhosOnFirst `json:"whosonfirst,omitempty"`
}

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
	Geometry   *GeocodeGeometry  `json:"geometry"`
	Properties GeocodeProperties `json:"properties"`
}

// GeocodeGeometry represents the GeoJSON geometry
type GeocodeGeometry struct {
	Type        string    `json:"type"` // "Point", "Polygon", etc.
	Coordinates []float64 `json:"coordinates"`
}

// GeocodeProperties contains the properties of a geocoding result
type GeocodeProperties struct {
	Addendum              map[string]interface{} `json:"addendum,omitempty"`
	AddressComponents     AddressComponents      `json:"address_components,omitempty"`
	Accuracy              string                 `json:"accuracy"`                // 'point'
	CoarseLocation        *string                `json:"coarse_location"`         // 'Riverton, UT, USA'
	Confidence            float64                `json:"confidence"`              // 1
	Context               Context                `json:"context,omitempty"`       // bunch of stuff
	Country               string                 `json:"country"`                 // 'United States'
	CountryA              string                 `json:"country_a"`               // 'USA'
	CountryCode           string                 `json:"country_code"`            // 'US'
	CountryGID            string                 `json:"country_gid"`             // 'whosonfirst:country:85633793'
	County                string                 `json:"county"`                  // "Tulare County"
	CountyA               string                 `json:"county_a"`                // 'TL'
	CountyGID             string                 `json:"county_gid"`              // 'whosonfirst:county:102082895'
	FormattedAddressLine  string                 `json:"formatted_address_line"`  // '123 Main Street, Riverton, Utah 84065, United States of America'
	FormattedAddressLines []string               `json:"formatted_address_lines"` // '123 Main Street', 'Riverton, Utah 84065', 'United States of America'
	GID                   string                 `json:"gid"`                     // 'openaddresses:address:us/ca/tulare-addresses-county:fe9dfab3d45c4550'
	HouseNumber           string                 `json:"housenumber"`             // '1234'
	ID                    string                 `json:"id"`                      // us/ca/tulare-addresses-county:fe9dfab3d45c4550
	Label                 string                 `json:"label"`                   // 1234 Main St, Dinuba, CA, USA
	Layer                 string                 `json:"layer"`                   // 'address'
	Locality              string                 `json:"locality"`                // 'Dinuba'
	LocalityGID           string                 `json:"locality_gid"`            // 'whosonfirst:locality:85922491'
	MatchType             string                 `json:"match_type"`              // 'exact'
	Name                  string                 `json:"name"`                    // '1234 Main St'
	PostalCode            string                 `json:"postalcode"`              // '93618'
	Precision             string                 `json:"precision"`               // 'centroid'
	Region                string                 `json:"region"`                  // 'California'
	RegionA               string                 `json:"region_a"`                // 'CA'
	RegionGID             string                 `json:"region_gid"`              // 'whosonfirst:region:85688637'
	Source                string                 `json:"source"`                  // 'openaddresses'
	Sources               []GeocodeSource        `json:"sources"`
	SourceID              string                 `json:"source_id"` // 'us/ca/tulare-addresses-county:fe9dfab3d45c4550'
	Street                string                 `json:"street"`    // 'Main Street'
}

// GeocodeSource represents a source of geocoding data
type GeocodeSource struct {
	FixitURL string `json:"fixit_url"`
	Source   string `json:"source"`
	SourceID string `json:"source_id"`
}

func (gf GeocodeFeature) CountryCode() string {
	if gf.Properties.CountryCode != "" {
		return gf.Properties.CountryCode
	}
	if gf.Properties.Context.ISO3166A3 != "" {
		return gf.Properties.Context.ISO3166A3
	}
	if gf.Properties.Context.WhosOnFirst.Country.Abbreviation != "" {
		return gf.Properties.Context.WhosOnFirst.Country.Abbreviation
	}
	return ""
}
func (gf GeocodeFeature) Locality() string {
	if gf.Properties.Locality != "" {
		return gf.Properties.Locality
	}
	if gf.Properties.Context.WhosOnFirst.Locality.Name != "" {
		return gf.Properties.Context.WhosOnFirst.Locality.Name
	}
	return ""
}
func (gf GeocodeFeature) Number() string {
	if gf.Properties.AddressComponents.Number != "" {
		return gf.Properties.AddressComponents.Number
	}
	if gf.Properties.HouseNumber != "" {
		return gf.Properties.HouseNumber
	}
	return ""
}
func (gf GeocodeFeature) PostalCode() string {
	if gf.Properties.PostalCode != "" {
		return gf.Properties.PostalCode
	}
	if gf.Properties.AddressComponents.PostalCode != "" {
		return gf.Properties.AddressComponents.PostalCode
	}
	return ""
}
func (gf GeocodeFeature) Region() string {
	if gf.Properties.Region != "" {
		return gf.Properties.Region
	}
	if gf.Properties.Context.WhosOnFirst.Region.Name != "" {
		return gf.Properties.Context.WhosOnFirst.Region.Name
	}
	return ""
}
func (gf GeocodeFeature) Street() string {
	if gf.Properties.Street != "" {
		return gf.Properties.Street
	}
	if gf.Properties.AddressComponents.Street != "" {
		return gf.Properties.AddressComponents.Street
	}
	return ""
}
