package platform
import (
	"time"

	"github.com/Gleipnir-Technology/nidus-sync/db/models"
)

type ClientSync struct {
	Fieldseeker FieldseekerRecordsSync
	Since time.Time
}

type FieldseekerRecordsSync struct {
	MosquitoSources models.FieldseekerPointlocationSlice
}
