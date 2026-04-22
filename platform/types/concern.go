package types

import (
	"fmt"

	"github.com/gorilla/mux"
	"github.com/rs/zerolog/log"
)

type Concern interface {
	PopulateURL(*mux.Router) error
}
type ConcernComplianceReportRequest struct {
	ComplianceReportRequestPublicID string `json:"compliance_report_request_public_id"`
	URL                             string `json:"url"`
}

func (e *ConcernComplianceReportRequest) PopulateURL(r *mux.Router) error {
	route_name := "compliance-request.image.pool.ByIDGet"
	handler := r.Get(route_name)
	if handler == nil {
		return fmt.Errorf("failed to get handler '%s'", route_name)
	}
	uri, err := handler.URL("public_id", e.ComplianceReportRequestPublicID)
	if err != nil {
		return fmt.Errorf("failed to create uri from '%s'", e.ComplianceReportRequestPublicID)
	}
	uri.Scheme = "https"
	e.URL = uri.String()
	log.Debug().Str("url", e.URL).Msg("populated concern URL")
	return nil
}
