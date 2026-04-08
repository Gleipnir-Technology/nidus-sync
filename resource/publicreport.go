package resource

import (
	"context"
	nhttp "github.com/Gleipnir-Technology/nidus-sync/http"
	"github.com/Gleipnir-Technology/nidus-sync/platform"
	"net/http"
	//"github.com/rs/zerolog/log"
	"github.com/gorilla/mux"
)

type publicreportR struct {
	router *router
}

type publicreport struct {
	ID       string `json:"id"`
	District string `json:"district"`
	URI      string `json:"uri"`
}

func Publicreport(r *router) *publicreportR {
	return &publicreportR{
		router: r,
	}
}

func (res *publicreportR) ByID(ctx context.Context, r *http.Request, query QueryParams) (*publicreport, *nhttp.ErrorWithStatus) {
	vars := mux.Vars(r)
	public_id := vars["id"]
	if public_id == "" {
		return nil, nhttp.NewBadRequest("You must provid an ID")
	}
	report, err := platform.PublicreportByID(ctx, public_id)
	if err != nil {
		return nil, nhttp.NewError("get report: %w", err)
	}
	district_uri, err := res.router.IDToURI("district.ByIDGet", int(report.OrganizationID))
	if err != nil {
		return nil, nhttp.NewError("district uri: %w", err)
	}
	return &publicreport{
		District: district_uri,
		ID:       report.PublicID,
	}, nil
}
