package types

type Lead struct {
	ComplianceReportRequest   *ComplianceReportRequest `db:"-" json:"compliance_report_request"`
	ComplianceReportRequestID *int32                   `db:"compliance_report_request_id" json:"-"`
	ID                        int32                    `db:"id" json:"id"`
	SiteID                    int32                    `db:"site_id" json:"site_id"`
	Type                      string                   `db:"type" json:"type"`
}
