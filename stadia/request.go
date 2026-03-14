package stadia

type RequestGeocode interface {
	SetBoundaryRect(xmin, ymin, xmax, ymax float64)
	SetFocusPoint(x, y float64)
}
