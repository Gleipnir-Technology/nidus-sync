package types

type Lead struct {
	ComplianceReportRequests []*ComplianceReportRequest `db:"-" json:"compliance_report_requests"`
	ID                       int32                      `db:"id" json:"id"`
	SiteID                   int32                      `db:"site_id" json:"site_id"`
	Type                     string                     `db:"type" json:"type"`
}
