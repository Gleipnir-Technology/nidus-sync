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

func (res *complianceR) Create(ctx context.Context, r *http.Request, n publicreportForm) (*compliance, *nhttp.ErrorWithStatus) {
	setter_report := models.PublicreportReportSetter{
		//AddressID:              omitnull.From(latlng.Cell.String()),
		AddressCountry:    omit.From(""),
		AddressGid:        omit.From(""),
		AddressNumber:     omit.From(""),
		AddressLocality:   omit.From(""),
		AddressPostalCode: omit.From(""),
		AddressRaw:        omit.From(""),
		AddressRegion:     omit.From(""),
		AddressStreet:     omit.From(""),
		Created:           omit.From(time.Now()),
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
		ReportType:    omit.From(enums.PublicreportReporttypeNuisance),
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
	uri, err := res.router.IDStrToURI("publicreport.ByIDGet", report.PublicID)
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
