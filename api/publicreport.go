package api

import (
	"context"
	"net/http"
	"strconv"

	"github.com/Gleipnir-Technology/nidus-sync/config"
	nhttp "github.com/Gleipnir-Technology/nidus-sync/http"
	"github.com/Gleipnir-Technology/nidus-sync/platform"
)

type formPublicreportSignal struct {
	ReportID string `json:"reportID"`
}
type createdSignal struct {
	ID int32 `json:"id"`
}

func postPublicreportSignal(ctx context.Context, r *http.Request, user platform.User, req formPublicreportSignal) (*createdSignal, *nhttp.ErrorWithStatus) {
	signal_id, err := platform.SignalCreateFromPublicreport(ctx, user, req.ReportID)
	if err != nil {
		return nil, nhttp.NewError("create signal: %w", err)
	}
	return &createdSignal{
		ID: *signal_id,
	}, nil
}

type formPublicreportInvalid struct {
	ReportID string `json:"reportID"`
}
type createdReport struct {
	URI string `json:"uri"`
}

func postPublicreportInvalid(ctx context.Context, r *http.Request, user platform.User, req formPublicreportSignal) (*createdReport, *nhttp.ErrorWithStatus) {
	err := platform.PublicreportInvalid(ctx, user, req.ReportID)
	if err != nil {
		return nil, nhttp.NewError("create signal: %w", err)
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
