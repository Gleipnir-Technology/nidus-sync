package types

import (
	"time"

	"github.com/Gleipnir-Technology/nidus-sync/db/models"
)

type Sync struct {
	Created          time.Time `json:"created"`
	ID               int32     `json:"id"`
	OrganizationID   int32     `json:"organization_id"`
	RecordsCreated   int32     `json:"records_created"`
	RecordsUnchanged int32     `json:"records_unchanged"`
	RecordsUpdated   int32     `json:"records_updated"`
}

func SyncFromModel(m *models.FieldseekerSync) Sync {
	//log.Debug().Int32("id", m.ID).Float64("lat", m.LocationLatitude.GetOr(0.0)).Float64("lng", m.LocationLongitude.GetOr(0.0)).Msg("converting address")
	return Sync{
		Created:          m.Created,
		ID:               m.ID,
		OrganizationID:   m.OrganizationID,
		RecordsCreated:   m.RecordsCreated,
		RecordsUnchanged: m.RecordsUnchanged,
		RecordsUpdated:   m.RecordsUpdated,
	}
}
