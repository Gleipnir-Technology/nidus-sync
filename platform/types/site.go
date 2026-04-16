package types

import (
	"time"
)

type Site struct {
	Address        Address           `db:"address" json:"address"`
	Created        time.Time         `db:"created" json:"created"`
	CreatorID      int32             `db:"creator_id" json:"creator_id"`
	FileID         int32             `db:"file_id" json:"file_id"`
	ID             int32             `db:"id" json:"id"`
	Notes          string            `db:"notes" json:"notes"`
	OrganizationID int32             `db:"organization_id" json:"organization_id"`
	Owner          *Contact          `db:"owner" json:"owner"`
	ParcelID       *int32            `db:"parcel_id" json:"parcel_id"`
	Resident       *Contact          `db:"resident" json:"resident"`
	ResidentOwned  bool              `db:"resident_owned" json:"resident_owned"`
	Tags           map[string]string `db:"tags" json:"tags"`
	Version        int32             `db:"version" json:"version"`
}
