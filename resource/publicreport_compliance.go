package resource

import (
	"context"
	"net/http"
	"time"

	"github.com/Gleipnir-Technology/nidus-sync/db/enums"
	"github.com/Gleipnir-Technology/nidus-sync/db/models"
	"github.com/aarondl/opt/omit"
	"github.com/aarondl/opt/omitnull"
	//"github.com/Gleipnir-Technology/nidus-sync/html"
	nhttp "github.com/Gleipnir-Technology/nidus-sync/http"
	"github.com/Gleipnir-Technology/nidus-sync/platform"
	"github.com/Gleipnir-Technology/nidus-sync/platform/types"
	"github.com/gorilla/mux"
	//"github.com/rs/zerolog/log"
)

func Compliance(r *router) *complianceR {
	return &complianceR{
		router: r,
	}
}

type complianceR struct {
	router *router
}
type compliance struct {
	District string `json:"district"`
	ID       string `json:"id"`
	URI      string `json:"uri"`
}

func (res *complianceR) ByID(ctx context.Context, r *http.Request, query QueryParams) (*types.PublicReportCompliance, *nhttp.ErrorWithStatus) {
	vars := mux.Vars(r)
	public_id := vars["id"]
	if public_id == "" {
		return nil, nhttp.NewBadRequest("You must provid an ID")
	}
	report, err := platform.PublicreportByIDCompliance(ctx, public_id)
	if err != nil {
		return nil, nhttp.NewError("get report: %w", err)
	}
	populateDistrictURI(&report.PublicReport, res.router)
	populateReportURI(&report.PublicReport, res.router)
	return report, nil
}
func (res *complianceR) Create(ctx context.Context, r *http.Request, n publicreportComplianceForm) (*compliance, *nhttp.ErrorWithStatus) {
	setter_report := models.PublicreportReportSetter{
		//AddressID:              omitnull.From(latlng.Cell.String()),
		AddressGid: omit.From(""),
		AddressRaw: omit.From(""),
		Created:    omit.From(time.Now()),
		//H3cell:              omitnull.From(latlng.Cell.String()),
		LatlngAccuracyType:  omit.From(enums.PublicreportAccuracytypeBrowser),
		LatlngAccuracyValue: omit.From(float32(0.0)),
		//Location: omitnull.From(fmt.Sprintf("ST_GeometryFromText(Point(%s %s))", longitude, latitude)),
		Location: omitnull.FromPtr[string](nil),
		MapZoom:  omit.From(float32(0.0)),
		//OrganizationID:    omitnull.FromPtr(organization_id),
		//PublicID:          omit.From(public_id),
		ReporterEmail: omit.From(""),
		ReporterName:  omit.From(""),
		ReporterPhone: omit.From(""),
		ReportType:    omit.From(enums.PublicreportReporttypeCompliance),
		Status:        omit.From(enums.PublicreportReportstatustypeReported),
	}
	setter_compliance := models.PublicreportComplianceSetter{
		AccessInstructions: omit.From(""),
		AvailabilityNotes:  omit.From(""),
		Comments:           omit.From(""),
		GateCode:           omit.From(""),
		HasDog:             omitnull.FromPtr[bool](nil),
		PermissionType:     omit.From(enums.PermissionaccesstypeUnselected),
		//ReportID            omit.Val[int32]
		ReportPhoneCanText: omitnull.FromPtr[bool](nil),
		WantsScheduled:     omitnull.FromPtr[bool](nil),
	}
	report, err := platform.PublicReportComplianceCreate(ctx, setter_report, setter_compliance)
	if err != nil {
		return nil, nhttp.NewError("create compliance report: %w", err)
	}
	uri, err := res.router.IDStrToURI("publicreport.compliance.ByIDGet", report.PublicID)
	if err != nil {
		return nil, nhttp.NewError("generate uri: %w", err)
	}
	district_uri, err := res.router.IDToURI("district.ByIDGet", int(report.OrganizationID))
	if err != nil {
		return nil, nhttp.NewError("generate district uri: %w", err)
	}
	return &compliance{
		District: district_uri,
		ID:       report.PublicID,
		URI:      uri,
	}, nil
}

func (res *complianceR) Update(ctx context.Context, r *http.Request, prf publicreportComplianceForm) (*types.PublicReportCompliance, *nhttp.ErrorWithStatus) {
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
	report, err := platform.PublicReportUpdateCompliance(ctx, public_id, report_setter, prf.Address, prf.Location)
	if err != nil {
		return nil, nhttp.NewError("platform update report compliance: %w", err)
	}
	return report, nil
}
