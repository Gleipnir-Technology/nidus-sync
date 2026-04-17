package resource

import (
	"context"
	"net/http"
	"strconv"

	nhttp "github.com/Gleipnir-Technology/nidus-sync/http"
	"github.com/Gleipnir-Technology/nidus-sync/platform"
	"github.com/Gleipnir-Technology/nidus-sync/platform/types"
	//"github.com/aarondl/opt/null"
	"github.com/gorilla/mux"
)

type siteR struct {
	router *router
}

func Site(r *router) *siteR {
	return &siteR{
		router: r,
	}
}

func (res *siteR) ByIDGet(ctx context.Context, r *http.Request, user platform.User, query QueryParams) (*types.Site, *nhttp.ErrorWithStatus) {
	vars := mux.Vars(r)
	id_str := vars["id"]
	id, err := strconv.Atoi(id_str)
	if err != nil {
		return nil, nhttp.NewBadRequest("'%s' is not a valid site ID: %w", id_str, err)
	}
	site, err := platform.SiteByID(ctx, user, int32(id))
	if err != nil {
		return nil, nhttp.NewError("site by id: %w", err)
	}
	return site, nil
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
	for _, site := range sites {
		uri, err := res.router.IDToURI("site.ByIDGet", int(site.ID))
		if err != nil {
			return nil, nhttp.NewError("set uri: %w", err)
		}
		site.URI = uri
	}
	return sites, nil
}
