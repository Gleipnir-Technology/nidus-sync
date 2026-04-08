package resource

import (
	"context"

	nhttp "github.com/Gleipnir-Technology/nidus-sync/http"
	"github.com/Gleipnir-Technology/nidus-sync/platform"
	"github.com/Gleipnir-Technology/nidus-sync/platform/types"
	"net/http"
	//"github.com/rs/zerolog/log"
	"github.com/gorilla/mux"
)

type publicreportR struct {
	router *router
}

func Publicreport(r *router) *publicreportR {
	return &publicreportR{
		router: r,
	}
}

func (res *publicreportR) ByID(ctx context.Context, r *http.Request, query QueryParams) (*types.Report, *nhttp.ErrorWithStatus) {
	vars := mux.Vars(r)
	public_id := vars["id"]
	if public_id == "" {
		return nil, nhttp.NewBadRequest("You must provid an ID")
	}
	report, err := platform.PublicreportByID(ctx, public_id)
	if err != nil {
		return nil, nhttp.NewError("get report: %w", err)
	}
	var district_uri string
	if report.DistrictID != nil {
		district_uri, err = res.router.IDToURI("district.ByIDGet", int(*report.DistrictID))
		if err != nil {
			return nil, nhttp.NewError("district uri: %w", err)
		}
	}
	uri, err := res.router.IDStrToURI("publicreport.ByIDGet", report.PublicID)
	if err != nil {
		return nil, nhttp.NewError("uri: %w", err)
	}
	report.District = &district_uri
	report.URI = uri
	return report, nil
}
