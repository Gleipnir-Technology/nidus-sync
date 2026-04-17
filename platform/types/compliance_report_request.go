package types

import (
	"github.com/Gleipnir-Technology/nidus-sync/db/models"
)

type ComplianceReportRequest struct {
	ID       int32  `db:"id" json:"id"`
	PublicID string `db:"public_id" json:"public_id"`
}

func ComplianceReportRequestFromModel(crr *models.ComplianceReportRequest) *ComplianceReportRequest {
	return &ComplianceReportRequest{
		ID:       crr.ID,
		PublicID: crr.PublicID,
	}
}
