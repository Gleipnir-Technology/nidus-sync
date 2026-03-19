package api

import (
	"context"
	"net/http"

	nhttp "github.com/Gleipnir-Technology/nidus-sync/http"
	"github.com/Gleipnir-Technology/nidus-sync/platform"
	//"github.com/aarondl/opt/null"
)

type contentListSignal struct {
	Signals []platform.Signal `json:"signals"`
}

func listSignal(ctx context.Context, r *http.Request, user platform.User, query queryParams) (*contentListSignal, *nhttp.ErrorWithStatus) {
	limit := 20
	if query.Limit != nil {
		limit = *query.Limit
	}
	signals, err := platform.SignalList(ctx, user, limit)
	if err != nil {
		return nil, nhttp.NewError("list signals: %w", err)
	}
	return &contentListSignal{
		Signals: signals,
	}, nil
}
