package api

import (
	"context"
	"net/http"
	"strconv"

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

type formPublicreportMessage struct {
	Message  string `json:"message"`
	ReportID string `json:"reportID"`
}
type createdMessage struct {
	URI string `json:"uri"`
}

func postPublicreportMessage(ctx context.Context, r *http.Request, user platform.User, req formPublicreportMessage) (*createdMessage, *nhttp.ErrorWithStatus) {
	msg_id, err := platform.PublicReportMessageCreate(ctx, user, req.ReportID, req.Message)
	if err != nil {
		return nil, nhttp.NewError("failed to create message: %s", err)
	}
	if msg_id == nil {
		return nil, nhttp.NewError("nil message id")
	}
	return &createdMessage{
		URI: config.MakeURLNidus("/message/%s", strconv.Itoa(int(*msg_id))),
	}, nil
}
