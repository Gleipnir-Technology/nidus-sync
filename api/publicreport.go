package api

import (
	"context"
	"net/http"

	"github.com/Gleipnir-Technology/nidus-sync/config"
	nhttp "github.com/Gleipnir-Technology/nidus-sync/http"
	"github.com/Gleipnir-Technology/nidus-sync/platform"
)

type formPublicreportLead struct {
	ReportID string `json:"reportID"`
}

func postPublicreportLead(ctx context.Context, r *http.Request, user platform.User, req formPublicreportLead) (*createdLead, *nhttp.ErrorWithStatus) {
	lead_id, err := platform.LeadCreateFromPublicreport(ctx, user, req.ReportID)
	if err != nil {
		return nil, nhttp.NewError("create lead: %w", err)
	}
	return &createdLead{
		ID: *lead_id,
	}, nil
}

type formPublicreportInvalid struct {
	ReportID string `json:"reportID"`
}
type createdReport struct {
	URI string `json:"uri"`
}

func postPublicreportInvalid(ctx context.Context, r *http.Request, user platform.User, req formPublicreportLead) (*createdReport, *nhttp.ErrorWithStatus) {
	err := platform.PublicreportInvalid(ctx, user, req.ReportID)
	if err != nil {
		return nil, nhttp.NewError("create lead: %w", err)
	}
	return &createdReport{
		URI: config.MakeURLNidus("/publicreport/%s", req.ReportID),
	}, nil
}
