package types

import (
	"encoding/json"

	"github.com/Gleipnir-Technology/nidus-sync/config"
	"github.com/google/uuid"
)

type Exif struct {
	Created string `json:"created"`
	Make    string `json:"make"`
	Model   string `json:"model"`
}
type Image struct {
	Exif         Exif      `db:"-" json:"exif"`
	ExifMake     string    `db:"exif_make" json:"-"`
	ExifModel    string    `db:"exif_model" json:"-"`
	ExifDateTime string    `db:"exif_datetime" json:"-"`
	Location     Location  `db:"location"`
	NuisanceID   int32     `db:"nuisance_id"`
	URLContent   string    `db:"-" json:"url_content"`
	UUID         uuid.UUID `db:"uuid"`
}

func (i *Image) MarshalJSON() ([]byte, error) {
	to_marshal := make(map[string]interface{}, 0)
	to_marshal["exif"] = Exif{
		Created: i.ExifDateTime,
		Make:    i.ExifMake,
		Model:   i.ExifModel,
	}
	to_marshal["location"] = i.Location
	to_marshal["nuisance_id"] = i.NuisanceID
	to_marshal["url_content"] = config.MakeURLNidus("/api/image/%s/content", i.UUID.String())
	to_marshal["uuid"] = i.UUID

	return json.Marshal(to_marshal)
}
