package types

import (
	"time"
)

type Mailer struct {
	Address                   Address   `json:"address"`
	ComplianceReportRequestID *string   `json:"compliance_report_request_id"`
	Created                   time.Time `json:"created"`
	ID                        int32     `json:"id"`
	Recipient                 string    `json:"recipient"`
	SiteID                    int32     `json:"site_id"`
	Status                    string    `json:"status"`
	URI                       string    `json:"uri"`
}
