package resource

import (
	"context"
	"net/http"

	"github.com/aarondl/opt/omit"
	//"github.com/aarondl/opt/omitnull"
	"github.com/Gleipnir-Technology/nidus-sync/db/models"
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

	platform.ReportImageCreate(ctx, public_id, uploads)
	return &image{Status: "ok"}, nil
}

type complianceForm struct {
	Comments *string `schema:"comments"`
}

type publicreportForm struct {
	Address    *types.Address  `schema:"address"`
	ClientID   string          `schema:"client_id"`
	Compliance *complianceForm `schema:"compliance"`
	DistrictID string          `schema:"district"`
	Location   *types.Location `schema:"location"`
	Locator    *Locator        `schema:"locator"`
	Reporter   *types.Contact  `schema:"reporter"`
}

func (res *publicreportR) Update(ctx context.Context, r *http.Request, prf publicreportForm) (*types.Report, *nhttp.ErrorWithStatus) {
	/*
		uploads, err := html.ExtractImageUploads(r)
		log.Info().Int("len", len(uploads)).Msg("extracted compliance uploads")
		if err != nil {
			return nil, nhttp.NewError("Failed to extract image uploads: %w", err)
		}
	*/
	vars := mux.Vars(r)
	public_id := vars["id"]
	if public_id == "" {
		return nil, nhttp.NewBadRequest("You must provide an ID")
	}
	report_setter := models.PublicreportReportSetter{}
	if prf.Location != nil {
		//report_setter.Latitude = omit.From(prf.Location.Latitude)
		//report_setter.Longitude = omit.From(prf.Location.Longitude)
		if prf.Location.Accuracy != nil {
			report_setter.LatlngAccuracyValue = omit.From(*prf.Location.Accuracy)
		}
	}
	if prf.Reporter != nil {
		if prf.Reporter.Email != nil {
			report_setter.ReporterEmail = omit.From(*prf.Reporter.Email)
		}
		if prf.Reporter.Name != nil {
			report_setter.ReporterName = omit.From(*prf.Reporter.Name)
		}
		if prf.Reporter.Phone != nil {
			report_setter.ReporterPhone = omit.From(*prf.Reporter.Phone)
		}
	}
	report, err := platform.PublicReportUpdate(ctx, public_id, report_setter, prf.Address, prf.Location)
	if err != nil {
		return nil, nhttp.NewError("update report: %w", err)
	}
	return report, nil
}
