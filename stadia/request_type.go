package stadia

// FocusPoint represents focus point coordinates
type FocusPoint struct {
	Lat *float64 `url:"focus.point.lat,omitempty"`
	Lon *float64 `url:"focus.point.lon,omitempty"`
}

// BoundaryRect represents a bounding rectangle
type BoundaryRect struct {
	MinLon *float64 `url:"boundary.rect.min_lon,omitempty"`
	MaxLon *float64 `url:"boundary.rect.max_lon,omitempty"`
	MinLat *float64 `url:"boundary.rect.min_lat,omitempty"`
	MaxLat *float64 `url:"boundary.rect.max_lat,omitempty"`
}

// BoundaryCircle represents a bounding circle
type BoundaryCircle struct {
	Lat    *float64 `url:"boundary.circle.lat,omitempty"`
	Lon    *float64 `url:"boundary.circle.lon,omitempty"`
	Radius *float64 `url:"boundary.circle.radius,omitempty"`
}
