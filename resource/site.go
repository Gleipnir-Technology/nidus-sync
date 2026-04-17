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

type siteR struct {
	router *router
}

func Site(r *router) *siteR {
	return &siteR{
		router: r,
	}
}

func (res *siteR) List(ctx context.Context, r *http.Request, user platform.User, query QueryParams) ([]*types.Site, *nhttp.ErrorWithStatus) {
	limit := 1000
	if query.Limit != nil {
		limit = *query.Limit
	}
	sites, err := platform.SiteList(ctx, user, limit)
	if err != nil {
		return nil, nhttp.NewError("list signals: %w", err)
	}
	return sites, nil
}
