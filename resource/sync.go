package resource

import (
	"context"
	"net/http"

	nhttp "github.com/Gleipnir-Technology/nidus-sync/http"
	"github.com/Gleipnir-Technology/nidus-sync/platform"
	"github.com/Gleipnir-Technology/nidus-sync/platform/types"
	//"github.com/aarondl/opt/null"
	"github.com/gorilla/mux"
)

type syncR struct {
	router *mux.Router
}

func Sync(r *mux.Router) *syncR {
	return &syncR{
		router: r,
	}
}

func (res *syncR) List(ctx context.Context, r *http.Request, user platform.User, query QueryParams) ([]*types.Sync, *nhttp.ErrorWithStatus) {
	limit := 20
	if query.Limit != nil {
		limit = *query.Limit
	}
	syncs, err := platform.SyncList(ctx, user, limit)
	if err != nil {
		return nil, nhttp.NewError("list signals: %w", err)
	}
	return syncs, nil
}
