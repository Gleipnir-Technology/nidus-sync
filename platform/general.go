package platform

import (
	"github.com/Gleipnir-Technology/nidus-sync/db/models"
)

func MosquitoSourceQuery() (*[]*MosquitoSource, error) {
	results := make([]*MosquitoSource, 0)
	return &results, nil
}

func ServiceRequestQuery() (*models.FieldseekerServicerequestSlice, error) {
	results := make(models.FieldseekerServicerequestSlice, 0)
	return &results, nil
}

func TrapDataQuery() (*models.FieldseekerTraplocationSlice, error) {
	results := make(models.FieldseekerTraplocationSlice, 0)
	return &results, nil
}
