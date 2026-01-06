package platform

import (
	"time"

	"github.com/Gleipnir-Technology/nidus-sync/db/models"
)

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
