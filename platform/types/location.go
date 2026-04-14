package types

import (
	"fmt"

	"github.com/Gleipnir-Technology/nidus-sync/db/models"
	"github.com/Gleipnir-Technology/nidus-sync/h3utils"
	//"github.com/rs/zerolog/log"
	"github.com/uber/h3-go/v4"
)

type Location struct {
	Accuracy  *float32 `db:"accuracy" json:"accuracy" schema:"accuracy"`
	Latitude  float64  `db:"latitude" json:"latitude" schema:"latitude"`
	Longitude float64  `db:"longitude" json:"longitude" schema:"longitude"`
}

func (l Location) String() string {
	return fmt.Sprintf("%f %f", l.Longitude, l.Latitude)
}

func (l Location) Resolution() uint {
	if l.Accuracy != nil {
		return uint(h3utils.MeterAccuracyToH3Resolution(float64(*l.Accuracy)))
	} else {
		return uint(0)
	}
}
func (l Location) H3Cell() (*h3.Cell, error) {
	result, err := h3utils.GetCell(l.Longitude, l.Latitude, int(l.Resolution()))
	return &result, err
}
func (l Location) GeometryQuery() (string, error) {
	return fmt.Sprintf("ST_Point(%f, %f, 4326)", l.Longitude, l.Latitude), nil
}
func LocationFromFS(pl *models.FieldseekerPointlocation) Location {
	return Location{}
}
