package resource

import (
	"context"
	"net/http"
	"time"

	"github.com/Gleipnir-Technology/nidus-sync/db/enums"
	modelpublicreport "github.com/Gleipnir-Technology/nidus-sync/db/gen/nidus-sync/publicreport/model"
	tablepublicreport "github.com/Gleipnir-Technology/nidus-sync/db/gen/nidus-sync/publicreport/table"
	querypublicreport "github.com/Gleipnir-Technology/nidus-sync/db/query/publicreport"

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

type publicReportComplianceForm struct {
	AccessInstructions omit.Val[string]                             `schema:"access_instructions" json:"access_instructions"`
	Address            omit.Val[types.Address]                      `schema:"address" json:"address"`
	AvailabilityNotes  omit.Val[string]                             `schema:"availability_notes"  json:"availability_notes"`
	ClientID           uuid.UUID                                    `schema:"client_id" json:"client_id"`
	Comments           omit.Val[string]                             `schema:"comments" json:"comments"`
	District           omit.Val[string]                             `schema:"district" json:"district"`
	GateCode           omit.Val[string]                             `schema:"gate_code" json:"gate_code"`
	HasDog             omitnull.Val[bool]                           `schema:"has_dog" json:"has_dog"`
	Location           omit.Val[types.Location]                     `schema:"location" json:"location"`
	MailerID           omit.Val[string]                             `schema:"mailer_id" json:"mailer_id"`
	PermissionType     omit.Val[enums.PublicreportPermissionaccess] `schema:"permission_type" json:"permission_type"`
	Reporter           omit.Val[types.Contact]                      `schema:"reporter" json:"reporter"`
	ReportPhoneCanSMS  omitnull.Val[bool]                           `schema:"report_phone_can_text"  json:"report_phone_can_text"`
	Submitted          omitnull.Val[time.Time]                      `schema:"submitted" json:"submitted"`
	WantsScheduled     omitnull.Val[bool]                           `schema:"wants_scheduled" json:"wants_scheduled"`
}

func (res *complianceR) ByID(ctx context.Context, r *http.Request, u platform.User, query QueryParams) (*types.PublicReportCompliance, *nhttp.ErrorWithStatus) {
	return res.byID(ctx, r, false)
}
func (res *complianceR) ByIDPublic(ctx context.Context, r *http.Request, query QueryParams) (*types.PublicReportCompliance, *nhttp.ErrorWithStatus) {
	return res.byID(ctx, r, true)
}
func (res *complianceR) Create(ctx context.Context, r *http.Request, n publicReportComplianceForm) (*types.PublicReportCompliance, *nhttp.ErrorWithStatus) {
	if n.District.IsUnset() && n.MailerID.IsUnset() {
		return nil, nhttp.NewBadRequest("You must provide a district_id or mailer_id")
	}
	user_agent := r.Header.Get("User-Agent")
	err := platform.EnsureClient(ctx, n.ClientID, user_agent)
	if err != nil {
		return nil, nhttp.NewError("Failed to ensure client: %w", err)
	}
	setter_report := modelpublicreport.Report{
		//AddressID:              omitnull.From(...),
		AddressGid: "",
		AddressRaw: "",
		ClientUUID: &n.ClientID,
		Created:    time.Now(),
		//H3cell:              omitnull.From(latlng.Cell.String()),
		LatlngAccuracyType:  modelpublicreport.Accuracytype_Browser,
		LatlngAccuracyValue: float32(0.0),
		//Location: omitnull.From(fmt.Sprintf("ST_GeometryFromText(Point(%s %s))", longitude, latitude)),
		Location: nil,
		MapZoom:  float32(0.0),
		//OrganizationID:      ,
		//PublicID:
		ReporterEmail:       "",
		ReporterName:        "",
		ReporterPhone:       "",
		ReporterPhoneCanSms: true,
		ReportType:          modelpublicreport.Reporttype_Compliance,
		Status:              modelpublicreport.Reportstatustype_Reported,
	}
	setter_compliance := modelpublicreport.Compliance{
		AccessInstructions: "",
		AvailabilityNotes:  "",
		Comments:           "",
		GateCode:           "",
		HasDog:             nil,
		PermissionType:     modelpublicreport.Permissionaccess_Unselected,
		//ReportID            omit.Val[int32]
		WantsScheduled: nil,
	}
	var org_id int32
	if n.District.IsValue() {
		district_str := n.District.MustGet()
		var district_id_ptr *int
		district_id_ptr, err := res.router.IDFromURI("district.ByIDGet", district_str)
		if err != nil || district_id_ptr == nil {
			return nil, nhttp.NewBadRequest("parse district ID: %w", err)
		}
		org_id = int32(*district_id_ptr)
		public_id, err := platform.GenerateReportID()
		if err != nil {
			return nil, nhttp.NewError("generate public ID: %w", err)
		}
		setter_report.PublicID = public_id
	}
	if n.MailerID.IsValue() {
		public_id := n.MailerID.MustGet()
		setter_report.PublicID = public_id

		// If it already exists, just return it
		report, err := platform.PublicReportByIDCompliance(ctx, public_id, true)
		if err != nil {
			return nil, nhttp.NewError("check existing report: %w", err)
		}
		if report != nil {
			return res.complianceHydrate(report, true)
		}

		org_id, err = platform.OrganizationIDForComplianceReportRequest(ctx, public_id)
		if err != nil {
			return nil, nhttp.NewBadRequest("no such mailer")
		}
		address, err := platform.AddressFromComplianceReportRequestID(ctx, public_id)
		if err != nil {
			return nil, nhttp.NewError("get address gid: %w", err)
		}
		setter_report.AddressID = address.ID
		setter_report.AddressGid = address.GID
	}
	report, err := platform.PublicReportComplianceCreate(ctx, setter_report, setter_compliance, org_id)
	if err != nil {
		return nil, nhttp.NewError("create compliance report: %w", err)
	}
	// Return a fully-fleshed-out report object, even though it's a bit more expensive
	result, err := platform.PublicReportByIDCompliance(ctx, report.PublicID, true)
	if err != nil {
		return nil, nhttp.NewError("get report after creation: %w", err)
	}
	return res.complianceHydrate(result, true)
}
func (res *complianceR) Update(ctx context.Context, r *http.Request, prf publicReportComplianceForm) (*types.PublicReportCompliance, *nhttp.ErrorWithStatus) {
	var err error
	vars := mux.Vars(r)
	public_id := vars["id"]
	if public_id == "" {
		return nil, nhttp.NewBadRequest("You must provide an ID")
	}
	report_updater := querypublicreport.NewReportUpdater()
	//report_setter := models.PublicreportReportSetter{}
	compliance_updater := querypublicreport.NewComplianceUpdater()
	//compliance_setter := models.PublicreportComplianceSetter{}
	var location *types.Location
	if prf.Location.IsValue() {
		l := prf.Location.MustGet()
		location = &l
		if location.Accuracy != nil {
			//report_setter.LatlngAccuracyValue = omit.From(*location.Accuracy)
			report_updater.Model.LatlngAccuracyValue = *location.Accuracy
			report_updater.Set(tablepublicreport.Report.LatlngAccuracyValue)
		}
	}
	if prf.Reporter.IsValue() {
		reporter := prf.Reporter.MustGet()
		if reporter.Email != nil {
			//report_setter.ReporterEmail = omit.From(*reporter.Email)
			report_updater.Model.ReporterEmail = *reporter.Email
			report_updater.Set(tablepublicreport.Report.ReporterEmail)
		}
		if reporter.Name != nil {
			//report_setter.ReporterName = omit.From(*reporter.Name)
			report_updater.Model.ReporterName = *reporter.Name
			report_updater.Set(tablepublicreport.Report.ReporterName)
		}
		if reporter.Phone != nil {
			//report_setter.ReporterPhone = omit.From(*reporter.Phone)
			report_updater.Model.ReporterPhone = *reporter.Phone
			report_updater.Set(tablepublicreport.Report.ReporterPhone)
		}
		if reporter.CanSMS != nil {
			//report_setter.ReporterPhoneCanSMS = omit.FromPtr(reporter.CanSMS)
			report_updater.Model.ReporterPhoneCanSms = *reporter.CanSMS
			report_updater.Set(tablepublicreport.Report.ReporterPhoneCanSms)
		}
	}
	var address *types.Address
	if prf.Address.IsValue() {
		a := prf.Address.MustGet()
		address = &a
	}
	if prf.AccessInstructions.IsValue() {
		//compliance_setter.AccessInstructions = prf.AccessInstructions
		compliance_updater.Model.AccessInstructions = prf.AccessInstructions.MustGet()
		compliance_updater.Set(tablepublicreport.Compliance.AccessInstructions)
	}
	if prf.AvailabilityNotes.IsValue() {
		//compliance_setter.AvailabilityNotes = prf.AvailabilityNotes
		compliance_updater.Model.AvailabilityNotes = prf.AvailabilityNotes.MustGet()
		compliance_updater.Set(tablepublicreport.Compliance.AvailabilityNotes)
	}
	if prf.Comments.IsValue() {
		//compliance_setter.Comments = prf.Comments
		compliance_updater.Model.Comments = prf.Comments.MustGet()
		compliance_updater.Set(tablepublicreport.Compliance.Comments)
	}
	if prf.GateCode.IsValue() {
		//compliance_setter.GateCode = prf.GateCode
		compliance_updater.Model.GateCode = prf.GateCode.MustGet()
		compliance_updater.Set(tablepublicreport.Compliance.GateCode)
	}
	if prf.HasDog.IsValue() {
		//compliance_setter.HasDog = prf.HasDog
		has_dog := prf.HasDog.MustGet()
		compliance_updater.Model.HasDog = &has_dog
		compliance_updater.Set(tablepublicreport.Compliance.HasDog)
	}
	if prf.PermissionType.IsValue() {
		//compliance_setter.PermissionType = prf.PermissionType
		var perm_type modelpublicreport.Permissionaccess
		pt := prf.PermissionType.MustGet()
		err = perm_type.Scan(pt)
		if err != nil {
			return nil, nhttp.NewBadRequest("permission type %s can't be scanned: %w", pt, err)
		}
		compliance_updater.Model.PermissionType = perm_type
		compliance_updater.Set(tablepublicreport.Compliance.PermissionType)
	}
	if prf.WantsScheduled.IsValue() {
		//compliance_setter.WantsScheduled = prf.WantsScheduled
		wants_scheduled := prf.WantsScheduled.MustGet()
		compliance_updater.Model.WantsScheduled = &wants_scheduled
		compliance_updater.Set(tablepublicreport.Compliance.WantsScheduled)
	}
	if prf.Submitted.IsValue() {
		log.Debug().Str("submitted", prf.Submitted.MustGet().String()).Msg("got submitted")
		//compliance_setter.Submitted = omitnull.From(time.Now())
		now := time.Now()
		compliance_updater.Model.Submitted = &now
		compliance_updater.Set(tablepublicreport.Compliance.Submitted)
	}
	err = platform.PublicReportUpdateCompliance(ctx, public_id, report_updater, compliance_updater, address, location)
	if err != nil {
		return nil, nhttp.NewError("platform update report compliance: %w", err)
	}
	// Return a fully-fleshed-out report object, even though it's a bit more expensive
	report, err := platform.PublicReportByIDCompliance(ctx, public_id, true)
	if err != nil {
		return nil, nhttp.NewError("get report after update: %w", err)
	}
	return res.complianceHydrate(report, true)
}

func (res *complianceR) byID(ctx context.Context, r *http.Request, is_public bool) (*types.PublicReportCompliance, *nhttp.ErrorWithStatus) {
	vars := mux.Vars(r)
	public_id := vars["id"]
	if public_id == "" {
		return nil, nhttp.NewBadRequest("You must provid an ID")
	}
	report, err := platform.PublicReportByIDCompliance(ctx, public_id, true)
	if err != nil {
		return nil, nhttp.NewError("get report: %w", err)
	}
	return res.complianceHydrate(report, is_public)
}
func (res *complianceR) complianceHydrate(report *types.PublicReportCompliance, is_public bool) (*types.PublicReportCompliance, *nhttp.ErrorWithStatus) {
	if err := populateDistrictURI(&report.PublicReport, res.router); err != nil {
		return nil, nhttp.NewError("populate district URI: %w", err)
	}
	if err := populateReportURI(&report.PublicReport, res.router, is_public); err != nil {
		return nil, nhttp.NewError("populate report URI: %w", err)
	}
	for _, e := range report.Concerns {
		if err := e.PopulateURL(res.router.router); err != nil {
			return nil, nhttp.NewError("populate concern URL: %w", err)
		}
	}
	return report, nil
}
