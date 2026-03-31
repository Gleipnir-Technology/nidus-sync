package types

import (
	"time"
)

type Site struct {
	Address        Address           `json:"address"`
	Created        time.Time         `json:"created"`
	CreatorID      int32             `json:"creator_id"`
	FileID         int32             `json:"file_id"`
	ID             int32             `json:"id"`
	Notes          string            `json:"notes"`
	OrganizationID int32             `json:"organization_id"`
	Owner          *Contact          `json:"owner"`
	ParcelID       *int32            `json:"parcel_id"`
	Resident       *Contact          `json:"resident"`
	ResidentOwned  bool              `json:"resident_owned"`
	Tags           map[string]string `json:"tags"`
	Version        int32             `json:"version"`
}
