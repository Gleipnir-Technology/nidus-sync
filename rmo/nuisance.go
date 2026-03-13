package rmo

import (
	"fmt"
	"net/http"
	"slices"
	"time"

	"github.com/Gleipnir-Technology/nidus-sync/db/enums"
	"github.com/Gleipnir-Technology/nidus-sync/db/models"
	"github.com/Gleipnir-Technology/nidus-sync/html"
	"github.com/Gleipnir-Technology/nidus-sync/platform"
	"github.com/Gleipnir-Technology/nidus-sync/platform/report"
	"github.com/aarondl/opt/omit"
	"github.com/aarondl/opt/omitnull"
	"github.com/rs/zerolog/log"
)

type ContentNuisance struct {
	District    *ContentDistrict
	MapboxToken string
	URL         ContentURL
}
type ContentNuisanceSubmitComplete struct {
	District *ContentDistrict
	ReportID string
	URL      ContentURL
}

func getNuisance(w http.ResponseWriter, r *http.Request) {
	html.RenderOrError(
		w,
		"rmo/nuisance.html",
		ContentNuisance{
			District: nil,
			URL:      makeContentURL(nil),
		},
	)
}
func getNuisanceDistrict(w http.ResponseWriter, r *http.Request) {
	district, err := districtBySlug(r)
	if err != nil {
		respondError(w, "Failed to lookup organization", err, http.StatusBadRequest)
		return
	}
	html.RenderOrError(
		w,
		"rmo/nuisance.html",
		ContentNuisance{
			District: newContentDistrict(district),
			URL:      makeContentURL(nil),
		},
	)
}
func getSubmitComplete(w http.ResponseWriter, r *http.Request) {
	report_id := r.URL.Query().Get("report")
	district, err := report.DistrictForReport(r.Context(), report_id)
	if err != nil {
		respondError(w, fmt.Sprintf("Failed to get district for report '%s'", report_id, err), err, http.StatusInternalServerError)
		return
	}
	html.RenderOrError(
		w,
		"rmo/submit-complete.html",
		ContentNuisanceSubmitComplete{
			District: newContentDistrict(district),
			ReportID: report_id,
			URL:      makeContentURL(nil),
		},
	)
}
func postNuisance(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	err := r.ParseMultipartForm(32 << 10) // 32 MB buffer
	if err != nil {
		respondError(w, "Failed to parse form", err, http.StatusBadRequest)
		return
	}
	additional_info := r.PostFormValue("additional-info")
	address_raw := r.PostFormValue("address")
	address_country := r.PostFormValue("address-country")
	address_locality := r.PostFormValue("address-locality")
	address_number := r.PostFormValue("address-number")
	address_postal_code := r.PostFormValue("address-postalcode")
	address_region := r.PostFormValue("address-region")
	address_street := r.PostFormValue("address-street")
	duration_str := postFormValueOrNone(r, "duration")
	source_stagnant := boolFromForm(r, "source-stagnant")
	source_container := boolFromForm(r, "source-container")
	source_description := r.PostFormValue("source-description")
	source_gutters := boolFromForm(r, "source-gutters")
	source_locations := r.Form["source-location"]
	tod_early := boolFromForm(r, "tod-early")
	tod_day := boolFromForm(r, "tod-day")
	tod_evening := boolFromForm(r, "tod-evening")
	tod_night := boolFromForm(r, "tod-night")

	duration := enums.PublicreportNuisancedurationtypeNone
	is_location_frontyard := false
	is_location_backyard := false
	is_location_garden := false
	is_location_pool := false
	is_location_other := false

	latlng, err := parseLatLng(r)

	err = duration.Scan(duration_str)
	if err != nil {
		log.Warn().Err(err).Str("duration_str", duration_str).Msg("Failed to interpret 'duration'")
	}

	//log.Debug().Strs("source_locations", source_locations).Msg("parsing")
	if slices.Contains(source_locations, "backyard") {
		is_location_backyard = true
	}
	if slices.Contains(source_locations, "frontyard") {
		is_location_frontyard = true
	}
	if slices.Contains(source_locations, "garden") {
		is_location_garden = true
	}
	if slices.Contains(source_locations, "other") {
		is_location_other = true
	}
	if slices.Contains(source_locations, "pool-area") {
		is_location_pool = true
	}

	uploads, err := extractImageUploads(r)
	log.Info().Int("len", len(uploads)).Msg("extracted uploads")
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
	setter := models.PublicreportNuisanceSetter{
		AdditionalInfo: omit.From(additional_info),
		//AddressID:              omitnull.From(latlng.Cell.String()),
		AddressRaw:        omit.From(address.Raw),
		AddressCountry:    omit.From(address.Country),
		AddressNumber:     omit.From(address.Number),
		AddressLocality:   omit.From(address.Locality),
		AddressPostalCode: omit.From(address.PostalCode),
		AddressRegion:     omit.From(address.Region),
		AddressStreet:     omit.From(address.Street),
		Created:           omit.From(time.Now()),
		Duration:          omit.From(duration),
		//H3cell:              omitnull.From(latlng.Cell.String()),
		IsLocationBackyard:  omit.From(is_location_backyard),
		IsLocationFrontyard: omit.From(is_location_frontyard),
		IsLocationGarden:    omit.From(is_location_garden),
		IsLocationOther:     omit.From(is_location_other),
		IsLocationPool:      omit.From(is_location_pool),
		LatlngAccuracyType:  omit.From(latlng.AccuracyType),
		LatlngAccuracyValue: omit.From(float32(latlng.AccuracyValue)),
		//Location: omitnull.From(fmt.Sprintf("ST_GeometryFromText(Point(%s %s))", longitude, latitude)),
		Location: omitnull.FromPtr[string](nil),
		MapZoom:  omit.From(latlng.MapZoom),
		//OrganizationID:    omitnull.FromPtr(organization_id),
		//PublicID:          omit.From(public_id),
		ReporterEmail:     omitnull.FromPtr[string](nil),
		ReporterName:      omitnull.FromPtr[string](nil),
		ReporterPhone:     omitnull.FromPtr[string](nil),
		SourceContainer:   omit.From(source_container),
		SourceDescription: omit.From(source_description),
		SourceGutter:      omit.From(source_gutters),
		SourceStagnant:    omit.From(source_stagnant),
		Status:            omit.From(enums.PublicreportReportstatustypeReported),
		TodEarly:          omit.From(tod_early),
		TodDay:            omit.From(tod_day),
		TodEvening:        omit.From(tod_evening),
		TodNight:          omit.From(tod_night),
	}
	public_id, err := platform.NuisanceCreate(ctx, setter, latlng, address, uploads)
	http.Redirect(w, r, fmt.Sprintf("/submit-complete?report=%s", public_id), http.StatusFound)
}
