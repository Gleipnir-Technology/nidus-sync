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
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/rs/zerolog/log"
)

func PublicReportCompliance(r *router) *complianceR {
	return &complianceR{
		router: r,
	}
}

type complianceR struct {
	router *router
}

type publicreportComplianceForm struct {
	AccessInstructions omit.Val[string]                     `schema:"access_instructions" json:"access_instructions"`
	Address            omit.Val[types.Address]              `schema:"address" json:"address"`
	AvailabilityNotes  omit.Val[string]                     `schema:"availability_notes"  json:"availability_notes"`
	ClientID           uuid.UUID                            `schema:"client_id" json:"client_id"`
	Comments           omit.Val[string]                     `schema:"comments" json:"comments"`
	District           omit.Val[string]                     `schema:"district" json:"district"`
	GateCode           omit.Val[string]                     `schema:"gate_code" json:"gate_code"`
	HasDog             omitnull.Val[bool]                   `schema:"has_dog" json:"has_dog"`
	Location           omit.Val[types.Location]             `schema:"location" json:"location"`
	MailerID           omit.Val[string]                     `schema:"mailer_id" json:"mailer_id"`
	PermissionType     omit.Val[enums.Permissionaccesstype] `schema:"permission_type" json:"permission_type"`
	Reporter           omit.Val[types.Contact]              `schema:"reporter" json:"reporter"`
	ReportPhoneCanSMS  omitnull.Val[bool]                   `schema:"report_phone_can_text"  json:"report_phone_can_text"`
	WantsScheduled     omitnull.Val[bool]                   `schema:"wants_scheduled" json:"wants_scheduled"`
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
func (res *complianceR) Create(ctx context.Context, r *http.Request, n publicreportComplianceForm) (*types.PublicReportCompliance, *nhttp.ErrorWithStatus) {
	if n.District.IsUnset() {
		return nil, nhttp.NewBadRequest("You must provide a district_id")
	}
	district_id, err := res.router.IDFromURI("district.ByIDGet", n.District.MustGet())
	if err != nil || district_id == nil {
		return nil, nhttp.NewBadRequest("parse district ID: %w", err)
	}
	user_agent := r.Header.Get("User-Agent")
	err = platform.EnsureClient(ctx, n.ClientID, user_agent)
	if err != nil {
		return nil, nhttp.NewError("Failed to ensure client: %w", err)
	}
	setter_report := models.PublicreportReportSetter{
		//AddressID:              omitnull.From(latlng.Cell.String()),
		AddressGid: omit.From(""),
		AddressRaw: omit.From(""),
		ClientUUID: omitnull.From(n.ClientID),
		Created:    omit.From(time.Now()),
		//H3cell:              omitnull.From(latlng.Cell.String()),
		LatlngAccuracyType:  omit.From(enums.PublicreportAccuracytypeBrowser),
		LatlngAccuracyValue: omit.From(float32(0.0)),
		//Location: omitnull.From(fmt.Sprintf("ST_GeometryFromText(Point(%s %s))", longitude, latitude)),
		Location: omitnull.FromPtr[string](nil),
		MapZoom:  omit.From(float32(0.0)),
		//OrganizationID:      omit.From[int32](int32(*district_id)),
		PublicID:            n.MailerID,
		ReporterEmail:       omit.From(""),
		ReporterName:        omit.From(""),
		ReporterPhone:       omit.From(""),
		ReporterPhoneCanSMS: omit.From(true),
		ReportType:          omit.From(enums.PublicreportReporttypeCompliance),
		Status:              omit.From(enums.PublicreportReportstatustypeReported),
	}
	setter_compliance := models.PublicreportComplianceSetter{
		AccessInstructions: omit.From(""),
		AvailabilityNotes:  omit.From(""),
		Comments:           omit.From(""),
		GateCode:           omit.From(""),
		HasDog:             omitnull.FromPtr[bool](nil),
		PermissionType:     omit.From(enums.PermissionaccesstypeUnselected),
		//ReportID            omit.Val[int32]
		WantsScheduled: omitnull.FromPtr[bool](nil),
	}
	report, err := platform.PublicReportComplianceCreate(ctx, setter_report, setter_compliance, int32(*district_id))
	if err != nil {
		return nil, nhttp.NewError("create compliance report: %w", err)
	}
	// Return a fully-fleshed-out report object, even though it's a bit more expensive
	result, err := platform.PublicreportByIDCompliance(ctx, report.PublicID)
	if err != nil {
		return nil, nhttp.NewError("get report after creation: %w", err)
	}
	populateDistrictURI(&result.PublicReport, res.router)
	populateReportURI(&result.PublicReport, res.router)
	return result, nil
}
func (res *complianceR) Update(ctx context.Context, r *http.Request, prf publicreportComplianceForm) (*types.PublicReportCompliance, *nhttp.ErrorWithStatus) {
	vars := mux.Vars(r)
	public_id := vars["id"]
	if public_id == "" {
		return nil, nhttp.NewBadRequest("You must provide an ID")
	}
	report_setter := models.PublicreportReportSetter{}
	compliance_setter := models.PublicreportComplianceSetter{}
	var location *types.Location
	if prf.Location.IsValue() {
		l := prf.Location.MustGet()
		location = &l
		if location.Accuracy != nil {
			report_setter.LatlngAccuracyValue = omit.From(*location.Accuracy)
		}
	}
	if prf.Reporter.IsValue() {
		reporter := prf.Reporter.MustGet()
		if reporter.Email != nil {
			report_setter.ReporterEmail = omit.From(*reporter.Email)
		}
		if reporter.Name != nil {
			report_setter.ReporterName = omit.From(*reporter.Name)
		}
		if reporter.Phone != nil {
			report_setter.ReporterPhone = omit.From(*reporter.Phone)
		}
		if reporter.CanSMS != nil {
			report_setter.ReporterPhoneCanSMS = omit.FromPtr(reporter.CanSMS)
		}
	}
	var address *types.Address
	if prf.Address.IsValue() {
		a := prf.Address.MustGet()
		address = &a
	}
	if prf.AccessInstructions.IsValue() {
		compliance_setter.AccessInstructions = prf.AccessInstructions
	}
	if prf.AvailabilityNotes.IsValue() {
		compliance_setter.AvailabilityNotes = prf.AvailabilityNotes
	}
	if prf.Comments.IsValue() {
		compliance_setter.Comments = prf.Comments
	}
	if prf.GateCode.IsValue() {
		compliance_setter.GateCode = prf.GateCode
	}
	if prf.HasDog.IsValue() {
		compliance_setter.HasDog = prf.HasDog
	}
	if prf.PermissionType.IsValue() {
		compliance_setter.PermissionType = prf.PermissionType
	}
	if prf.WantsScheduled.IsValue() {
		compliance_setter.WantsScheduled = prf.WantsScheduled
	}
	log.Debug().
		Bool("access_instructions", prf.AccessInstructions.IsValue()).
		Bool("access_instructions", prf.AccessInstructions.IsValue()).
		Bool("access_instructions", prf.AccessInstructions.IsValue()).
		Bool("access_instructions", prf.AccessInstructions.IsValue()).
		Msg("updating compliance")
	report, err := platform.PublicReportUpdateCompliance(ctx, public_id, &report_setter, &compliance_setter, address, location)
	if err != nil {
		return nil, nhttp.NewError("platform update report compliance: %w", err)
	}
	// Return a fully-fleshed-out report object, even though it's a bit more expensive
	report, err = platform.PublicreportByIDCompliance(ctx, public_id)
	if err != nil {
		return nil, nhttp.NewError("get report after update: %w", err)
	}
	return report, nil
}
