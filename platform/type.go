package platform

import (
	"time"

	"github.com/Gleipnir-Technology/nidus-sync/db/models"
	"github.com/Gleipnir-Technology/nidus-sync/platform/types"
)

type Location = types.Location

type ClientSync struct {
	Fieldseeker FieldseekerRecordsSync
	Since       time.Time
}

type FieldseekerRecordsSync struct {
	MosquitoSources []MosquitoSource
	ServiceRequests models.FieldseekerServicerequestSlice
	TrapData        models.FieldseekerTraplocationSlice
}

type MosquitoSource struct {
	PointLocation models.FieldseekerPointlocation
	Inspections   models.FieldseekerMosquitoinspectionSlice
	Treatments    models.FieldseekerTreatmentSlice
}
