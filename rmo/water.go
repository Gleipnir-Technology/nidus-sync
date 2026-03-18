package rmo

import (
	"fmt"
	"net/http"
	"time"

	"github.com/Gleipnir-Technology/nidus-sync/db/enums"
	"github.com/Gleipnir-Technology/nidus-sync/db/models"
	"github.com/Gleipnir-Technology/nidus-sync/html"
	"github.com/Gleipnir-Technology/nidus-sync/platform"
	"github.com/aarondl/opt/omit"
)

type ContentWater struct {
	District *ContentDistrict
	URL      ContentURL
}

func getWater(w http.ResponseWriter, r *http.Request) {
	html.RenderOrError(
		w,
		"rmo/water.html",
		ContentWater{
			District: nil,
			URL:      makeContentURL(nil),
		},
	)
}
func getWaterDistrict(w http.ResponseWriter, r *http.Request) {
	district, err := districtBySlug(r)
	if err != nil {
		respondError(w, "Failed to lookup organization", err, http.StatusBadRequest)
		return
	}
	html.RenderOrError(
		w,
		"rmo/water.html",
		ContentWater{
			District: newContentDistrict(district),
			URL:      makeContentURL(district),
		},
	)
}
func postWater(w http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(32 << 10) // 32 MB buffer
	if err != nil {
		respondError(w, "Failed to parse form", err, http.StatusBadRequest)
		return
	}

	access_comments := r.FormValue("access-comments")
	access_dog := boolFromForm(r, "access-dog")
	access_fence := boolFromForm(r, "access-fence")
	access_gate := boolFromForm(r, "access-gate")
	access_locked := boolFromForm(r, "access-locked")
	access_other := boolFromForm(r, "access-other")
	address_raw := r.FormValue("address")
	address_country := r.FormValue("address-country")
	address_locality := r.FormValue("address-locality")
	address_number := r.FormValue("address-number")
	address_postal_code := r.FormValue("address-postalcode")
	address_region := r.FormValue("address-region")
	address_street := r.FormValue("address-street")
	comments := r.FormValue("comments")
	has_adult := boolFromForm(r, "has-adult")
	has_backyard_permission := boolFromForm(r, "backyard-permission")
	has_larvae := boolFromForm(r, "has-larvae")
	has_pupae := boolFromForm(r, "has-pupae")
	is_reporter_confidential := boolFromForm(r, "reporter-confidential")
	is_reporter_owner := boolFromForm(r, "property-ownership")
	owner_email := r.FormValue("owner-email")
	owner_name := r.FormValue("owner-name")
	owner_phone := r.FormValue("owner-phone")

	latlng, err := parseLatLng(r)
	if err != nil {
		respondError(w, "Failed to parse lat lng for water report", err, http.StatusInternalServerError)
		return
	}

	ctx := r.Context()

	uploads, err := extractImageUploads(r)
	if err != nil {
		respondError(w, "Failed to extract image uploads", err, http.StatusInternalServerError)
		return
	}

	address := platform.Address{
		Country:    address_country,
		Locality:   address_locality,
		Number:     address_number,
		PostalCode: address_postal_code,
		Raw:        address_raw,
		Region:     address_region,
		Street:     address_street,
		Unit:       "",
	}
	setter_report := models.PublicreportReportSetter{
		AddressRaw:        omit.From(address_raw),
		AddressCountry:    omit.From(address_country),
		AddressLocality:   omit.From(address_locality),
		AddressNumber:     omit.From(address_number),
		AddressPostalCode: omit.From(address_postal_code),
		AddressStreet:     omit.From(address_street),
		AddressRegion:     omit.From(address_region),
		Created:           omit.From(time.Now()),
		//H3cell:       omitnull.From(geospatial.Cell.String()),
		LatlngAccuracyType:  omit.From(latlng.AccuracyType),
		LatlngAccuracyValue: omit.From(float32(latlng.AccuracyValue)),
		//Location: add later
		MapZoom: omit.From(latlng.MapZoom),
		//OrganizationID: omitnull.FromPtr(organization_id),
		//PublicID:       omit.From(public_id),
		ReporterEmail: omit.From(""),
		ReporterName:  omit.From(""),
		ReporterPhone: omit.From(""),
		ReportType:    omit.From(enums.PublicreportReporttypeWater),
		Status:        omit.From(enums.PublicreportReportstatustypeReported),
	}
	setter_water := models.PublicreportWaterSetter{
		AccessComments:         omit.From(access_comments),
		AccessDog:              omit.From(access_dog),
		AccessFence:            omit.From(access_fence),
		AccessGate:             omit.From(access_gate),
		AccessLocked:           omit.From(access_locked),
		AccessOther:            omit.From(access_other),
		Comments:               omit.From(comments),
		HasAdult:               omit.From(has_adult),
		HasBackyardPermission:  omit.From(has_backyard_permission),
		HasLarvae:              omit.From(has_larvae),
		HasPupae:               omit.From(has_pupae),
		IsReporterConfidential: omit.From(is_reporter_confidential),
		IsReporterOwner:        omit.From(is_reporter_owner),
		OwnerEmail:             omit.From(owner_email),
		OwnerName:              omit.From(owner_name),
		OwnerPhone:             omit.From(owner_phone),
		//ReportID               omit.Val[int32]
	}
	report, err := platform.ReportWaterCreate(ctx, setter_report, setter_water, latlng, address, uploads)
	if err != nil {
		respondError(w, "Failed to save new report", err, http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, fmt.Sprintf("/submit-complete?report=%s", report.PublicID), http.StatusFound)
}
func postWaterDistrict(w http.ResponseWriter, r *http.Request) {
}
