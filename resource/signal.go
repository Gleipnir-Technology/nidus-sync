package resource

import (
	"context"
	"net/http"

	nhttp "github.com/Gleipnir-Technology/nidus-sync/http"
	"github.com/Gleipnir-Technology/nidus-sync/platform"
	//"github.com/aarondl/opt/null"
	"github.com/gorilla/mux"
)

type signalR struct {
	router *mux.Router
}

func Signal(r *mux.Router) *signalR {
	return &signalR{
		router: r,
	}
}

type contentListSignal struct {
	Signals []*platform.Signal `json:"signals"`
}

func (res *signalR) List(ctx context.Context, r *http.Request, user platform.User, query QueryParams) (*contentListSignal, *nhttp.ErrorWithStatus) {
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
