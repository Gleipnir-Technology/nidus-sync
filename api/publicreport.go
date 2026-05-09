package api

import (
	"context"
	"fmt"
	"net/http"

	nhttp "github.com/Gleipnir-Technology/nidus-sync/http"
	"github.com/Gleipnir-Technology/nidus-sync/platform"
)

type formPublicreportSignal struct {
	ReportID string `json:"reportID"`
}

func postPublicreportSignal(ctx context.Context, r *http.Request, user platform.User, req formPublicreportSignal) (string, *nhttp.ErrorWithStatus) {
	signal_id, err := platform.SignalCreateFromPublicreport(ctx, user, req.ReportID)
	if err != nil {
		return "", nhttp.NewError("create signal: %w", err)
	}
	return fmt.Sprintf("/signal/%d", *signal_id), nil
}

func postPublicreportInvalid(ctx context.Context, r *http.Request, user platform.User, req formPublicreportSignal) (string, *nhttp.ErrorWithStatus) {
	err := platform.PublicReportInvalid(ctx, user, req.ReportID)
	if err != nil {
		return "", nhttp.NewError("create signal: %w", err)
	}
	return fmt.Sprintf("/publicreport/%s", req.ReportID), nil
}

type formPublicreportMessage struct {
	Message  string `json:"message"`
	ReportID string `json:"reportID"`
}

func postPublicreportMessage(ctx context.Context, r *http.Request, user platform.User, req formPublicreportMessage) (string, *nhttp.ErrorWithStatus) {
	msg_id, err := platform.PublicReportMessageCreate(ctx, user, req.ReportID, req.Message)
	if err != nil {
		return "", nhttp.NewError("failed to create message: %s", err)
	}
	if msg_id == nil {
		return "", nhttp.NewError("nil message id")
	}
	return fmt.Sprintf("/message/%d", *msg_id), nil
}
