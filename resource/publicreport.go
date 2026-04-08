package resource

import (
	"context"
	"time"

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

type publicreport struct {
	Address    string         `json:"address"`
	Created    time.Time      `json:"created"`
	District   string         `json:"district"`
	ID         string         `json:"id"`
	ImageCount int            `json:"image_count"`
	Location   types.Location `json:"location"`
	Status     string         `json:"status"`
	Type       string         `json:"type"`
	URI        string         `json:"uri"`
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
	uri, err := res.router.IDStrToURI("publicreport.ByIDGet", report.PublicID)
	if err != nil {
		return nil, nhttp.NewError("uri: %w", err)
	}
	location := types.Location{
		Latitude:  report.LocationLatitude.GetOr(0.0),
		Longitude: report.LocationLongitude.GetOr(0.0),
	}
	return &publicreport{
		District:   district_uri,
		ID:         report.PublicID,
		Address:    report.AddressRaw,
		Created:    report.Created,
		ImageCount: len(report.R.Images),
		Location:   location,
		Status:     report.Status.String(),
		Type:       report.ReportType.String(),
		URI:        uri,
	}, nil
}
