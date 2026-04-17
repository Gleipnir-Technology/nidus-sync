package types

type ComplianceReportRequest struct {
	ID       int32  `db:"id" json:"id"`
	PublicID string `db:"public_id" json:"public_id"`
}
