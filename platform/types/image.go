package types

import (
	"encoding/json"
	"fmt"
	"math"
	"time"

	"github.com/Gleipnir-Technology/nidus-sync/config"
	"github.com/google/uuid"
	//"github.com/rs/zerolog/log"
)

type Exif struct {
	Created string `json:"created"`
	Make    string `json:"make"`
	Model   string `json:"model"`
}

func (e Exif) MarshalJSON() ([]byte, error) {
	to_marshal := make(map[string]interface{}, 0)
	if e.Created != "" {
		layout := "2006:01:02 15:04:05"

		t, err := time.Parse(layout, e.Created)
		if err != nil {
			fmt.Println("Error parsing date:", err)
			return nil, fmt.Errorf("parse created exif: %w", err)
		}
		to_marshal["created"] = t
	} else {
		to_marshal["created"] = e.Created
	}
	to_marshal["make"] = e.Make
	to_marshal["model"] = e.Model
	return json.Marshal(to_marshal)
}

type Image struct {
	DistanceToReporterMeters *float64  `db:"distance_from_reporter_meters"`
	Exif                     Exif      `db:"-" json:"exif"`
	ExifMake                 string    `db:"exif_make" json:"-"`
	ExifModel                string    `db:"exif_model" json:"-"`
	ExifDateTime             string    `db:"exif_datetime" json:"-"`
	Location                 *Location `db:"location"`
	ReportID                 int32     `db:"report_id" json:"-"`
	URLContent               string    `db:"-" json:"url_content"`
	UUID                     uuid.UUID `db:"uuid"`
}

func (i *Image) MarshalJSON() ([]byte, error) {
	to_marshal := make(map[string]interface{}, 0)
	if i.DistanceToReporterMeters != nil && math.IsNaN(*i.DistanceToReporterMeters) {
		to_marshal["distance_from_reporter_meters"] = nil
	} else {
		to_marshal["distance_from_reporter_meters"] = i.DistanceToReporterMeters
	}
	to_marshal["exif"] = Exif{
		Created: i.ExifDateTime,
		Make:    i.ExifMake,
		Model:   i.ExifModel,
	}
	if math.IsNaN(i.Location.Latitude) || math.IsNaN(i.Location.Longitude) {
		to_marshal["location"] = nil
	} else {
		to_marshal["location"] = i.Location
	}
	//to_marshal["report_id"] = i.ReportID
	to_marshal["url_content"] = config.MakeURLNidus("/api/image/%s/content", i.UUID.String())
	to_marshal["uuid"] = i.UUID

	return json.Marshal(to_marshal)
}
