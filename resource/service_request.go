package resource

import (
	"context"
	"net/http"

	nhttp "github.com/Gleipnir-Technology/nidus-sync/http"
	"github.com/Gleipnir-Technology/nidus-sync/platform"
	"github.com/Gleipnir-Technology/nidus-sync/platform/types"
	//"github.com/aarondl/opt/null"
	//"github.com/gorilla/mux"
)

type serviceRequestR struct {
	router *router
}

func ServiceRequest(r *router) *serviceRequestR {
	return &serviceRequestR{
		router: r,
	}
}

func (res *serviceRequestR) List(ctx context.Context, r *http.Request, user platform.User, query QueryParams) ([]*types.ServiceRequest, *nhttp.ErrorWithStatus) {
	limit := 20
	if query.Limit != nil {
		limit = *query.Limit
	}
	serviceRequests, err := platform.ServiceRequestList(ctx, user, limit)
	if err != nil {
		return nil, nhttp.NewError("list signals: %w", err)
	}
	return serviceRequests, nil
}
