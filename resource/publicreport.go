package resource

import (
	"context"
	"fmt"
	"net/http"

	"github.com/Gleipnir-Technology/nidus-sync/html"
	nhttp "github.com/Gleipnir-Technology/nidus-sync/http"
	"github.com/Gleipnir-Technology/nidus-sync/platform"
	"github.com/Gleipnir-Technology/nidus-sync/platform/types"
	"github.com/gorilla/mux"
	"github.com/rs/zerolog/log"
)

type publicreportR struct {
	router *router
}

func Publicreport(r *router) *publicreportR {
	return &publicreportR{
		router: r,
	}
}

func (res *publicreportR) ByID(ctx context.Context, r *http.Request, query QueryParams) (*types.PublicReport, *nhttp.ErrorWithStatus) {
	vars := mux.Vars(r)
	public_id := vars["id"]
	if public_id == "" {
		return nil, nhttp.NewBadRequest("You must provid an ID")
	}
	report, err := platform.PublicreportByID(ctx, public_id)
	if err != nil {
		return nil, nhttp.NewError("get report: %w", err)
	}
	populateDistrictURI(report, res.router)
	populateReportURI(report, res.router)
	return report, nil
}

type image struct {
	Status string `json:"status"`
}

func (res *publicreportR) ImageCreate(ctx context.Context, r *http.Request, n nuisanceForm) (*image, *nhttp.ErrorWithStatus) {
	vars := mux.Vars(r)
	public_id := vars["id"]
	if public_id == "" {
		return nil, nhttp.NewBadRequest("You must provide an ID")
	}

	uploads, err := html.ExtractImageUploads(r)
	log.Info().Int("len", len(uploads)).Msg("report image uploads")
	if err != nil {
		return nil, nhttp.NewError("Failed to extract image uploads: %w", err)
	}

	platform.PublicReportImageCreate(ctx, public_id, uploads)
	return &image{Status: "ok"}, nil
}

func populateDistrictURI(report *types.PublicReport, r *router) error {
	var district_uri string
	var err error
	if report.DistrictID != nil {
		district_uri, err = r.IDToURI("district.ByIDGet", int(*report.DistrictID))
		if err != nil {
			return nhttp.NewError("district uri: %w", err)
		}
	}
	report.District = &district_uri
	return nil
}
func populateReportURI(report *types.PublicReport, r *router) error {
	var route_name string
	switch report.Type {
	case "compliance":
		route_name = "publicreport.compliance.ByIDGet"
	case "nuisance":
		route_name = "publicreport.nuisance.ByIDGet"
	case "water":
		route_name = "publicreport.water.ByIDGet"
	default:
		return fmt.Errorf("Unrecognized report type '%s'", report.Type)
	}
	uri, err := r.IDStrToURI(route_name, report.PublicID)
	if err != nil {
		return nhttp.NewError("uri: %w", err)
	}
	report.URI = uri
	return nil
}
