package types

import (
	"github.com/google/uuid"
)

type Image struct {
	ExifMake     string    `db:"exif_make"`
	ExifModel    string    `db:"exif_model"`
	ExifDateTime string    `db:"exif_datetime"`
	Location     Location  `db:"location"`
	NuisanceID   int32     `db:"nuisance_id"`
	UUID         uuid.UUID `db:"uuid"`
}
