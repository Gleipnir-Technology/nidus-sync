package main
import (
	"encoding/json"
	"errors"
	//"fmt"

	"github.com/stephenafamo/bob/types"
	"github.com/uber/h3-go/v4"
)
func getPoint(geom types.JSON[json.RawMessage]) (h3.LatLng, error) {
	/*
	I need to figure out how to convert from the weird bob
	type of types.JSON[json.RawMessage] to a X/Y coordinate to
	here, but I need an Internet connection to do that effectively.
	*/
	return h3.LatLng{}, errors.New("need to implement this")
	/*
	msg, err := geom.Value()
	if err != nil {
		return h3.LatLng{}, fmt.Errorf("Can't get underlying JSON value: %w", err)
	}
	x := msg.Get("x")
	y := msg.Get("y")
	return h3.NewLatLng(x, y), nil
	*/
}
