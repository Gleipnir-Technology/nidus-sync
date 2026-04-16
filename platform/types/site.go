package types

import (
	"time"

	"github.com/Gleipnir-Technology/nidus-sync/db/models"
)

type Site struct {
	Address        Address           `db:"address" json:"address"`
	Created        time.Time         `db:"created" json:"created"`
	CreatorID      int32             `db:"creator_id" json:"creator_id"`
	FileID         int32             `db:"file_id" json:"file_id"`
	ID             int32             `db:"id" json:"id"`
	Notes          string            `db:"notes" json:"notes"`
	OrganizationID int32             `db:"organization_id" json:"organization_id"`
	Owner          Contact           `db:"owner" json:"owner"`
	ParcelID       *int32            `db:"parcel_id" json:"parcel_id"`
	Resident       *Contact          `db:"resident" json:"resident"`
	ResidentOwned  *bool             `db:"resident_owned" json:"resident_owned"`
	Tags           map[string]string `db:"tags" json:"tags"`
	Version        int32             `db:"version" json:"version"`
}

func SiteFromModel(s *models.Site) Site {
	owner_phone := s.OwnerPhoneE164.GetOr("")
	var resident_owned *bool
	if s.ResidentOwned.IsValue() {
		b := s.ResidentOwned.MustGet()
		resident_owned = &b
	}
	return Site{
		Created:   s.Created,
		CreatorID: s.CreatorID,
		//FileID: s.FileID,
		ID:             s.ID,
		Notes:          s.Notes,
		OrganizationID: s.OrganizationID,
		Owner: Contact{
			Name:  &s.OwnerName,
			Phone: &owner_phone,
		},
		ResidentOwned: resident_owned,
		//ParcelID: s.ParcelID,
	}
}
