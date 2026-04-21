package types

import (
	"fmt"

	"github.com/gorilla/mux"
	"github.com/rs/zerolog/log"
)

type Evidence interface {
	PopulateURL(*mux.Router) error
}
type EvidenceComplianceReportRequest struct {
	ComplianceReportRequestPublicID string
	URL                             string
}

func (e *EvidenceComplianceReportRequest) PopulateURL(r *mux.Router) error {
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
	log.Debug().Str("url", e.URL).Msg("populated evidence URL")
	return nil
}
