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

func (res *publicreportR) ByID(ctx context.Context, w http.ResponseWriter, r *http.Request) *nhttp.ErrorWithStatus {
	vars := mux.Vars(r)
	public_id := vars["id"]
	if public_id == "" {
		return nhttp.NewBadRequest("You must provide an ID")
	}
	report_type, err := platform.PublicReportTypeByID(ctx, public_id)
	if err != nil {
		return nhttp.NewError("get report '%s': %w", public_id, err)
	}
	path, err := reportURI(res.router, report_type, public_id)
	if err != nil {
		return nhttp.NewError("get uri '%s': %w", public_id, err)
	}
	http.Redirect(w, r, path, http.StatusFound)
	return nil
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
	uri, err := reportURI(r, report.Type, report.PublicID)
	if err != nil {
		return fmt.Errorf("report uri: %w", err)
	}
	report.URI = uri
	return nil
}
func reportURI(r *router, report_type string, public_id string) (string, error) {
	var route_name string
	switch report_type {
	case "compliance":
		route_name = "publicreport.compliance.ByIDGet"
	case "nuisance":
		route_name = "publicreport.nuisance.ByIDGet"
	case "water":
		route_name = "publicreport.water.ByIDGet"
	default:
		return "", fmt.Errorf("Unrecognized report type '%s'", report_type)
	}
	uri, err := r.IDStrToURI(route_name, public_id)
	if err != nil {
		return "", fmt.Errorf("id str to uri '%s' '%s': %w", route_name, public_id, err)
	}
	return uri, nil
}
