package rmo

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/Gleipnir-Technology/nidus-sync/config"
	"github.com/Gleipnir-Technology/nidus-sync/db"
	"github.com/Gleipnir-Technology/nidus-sync/db/enums"
	"github.com/Gleipnir-Technology/nidus-sync/db/models"
	"github.com/Gleipnir-Technology/nidus-sync/html"
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
	ReportID string
}

var (
	Nuisance               = buildTemplate("nuisance", "base")
	NuisanceSubmitComplete = buildTemplate("nuisance-submit-complete", "base")
)

func getNuisance(w http.ResponseWriter, r *http.Request) {
	html.RenderOrError(
		w,
		Nuisance,
		ContentNuisance{
			District:    nil,
			MapboxToken: config.MapboxToken,
			URL:         makeContentURL(),
		},
	)
}
func getNuisanceSubmitComplete(w http.ResponseWriter, r *http.Request) {
	report := r.URL.Query().Get("report")
	html.RenderOrError(
		w,
		NuisanceSubmitComplete,
		ContentNuisanceSubmitComplete{
			ReportID: report,
		},
	)
}
func postNuisance(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		respondError(w, "Failed to parse form", err, http.StatusBadRequest)
		return
	}
	tod_early := boolFromForm(r, "tod-early")
	tod_day := boolFromForm(r, "tod-day")
	tod_evening := boolFromForm(r, "tod-evening")
	tod_night := boolFromForm(r, "tod-night")

	source_stagnant := boolFromForm(r, "source-stagnant")
	source_container := boolFromForm(r, "source-container")
	source_roof := boolFromForm(r, "source-container")

	request_call := boolFromForm(r, "request-call")

	duration_str := postFormValueOrNone(r, "duration")
	var duration enums.PublicreportNuisancedurationtype
	err = duration.Scan(duration_str)
	if err != nil {
		respondError(w, fmt.Sprintf("Failed to interpret 'duration' of '%s'", duration_str), err, http.StatusBadRequest)
		return
	}

	inspection_type_str := postFormValueOrNone(r, "inspection-type")
	var inspection_type enums.PublicreportNuisanceinspectiontype
	err = inspection_type.Scan(inspection_type_str)
	if err != nil {
		respondError(w, fmt.Sprintf("Failed to interpret 'inspection-type' of '%s'", inspection_type_str), err, http.StatusBadRequest)
		return
	}

	source_location_str := postFormValueOrNone(r, "source-location")
	var source_location enums.PublicreportNuisancelocationtype
	err = source_location.Scan(source_location_str)
	if err != nil {
		respondError(w, fmt.Sprintf("Failed to interpret 'source-location' of '%s'", source_location_str), err, http.StatusBadRequest)
		return
	}
	preferred_date_range_str := postFormValueOrNone(r, "preferred-date-range")
	var preferred_date_range enums.PublicreportNuisancepreferreddaterangetype
	err = preferred_date_range.Scan(preferred_date_range_str)
	if err != nil {
		respondError(w, fmt.Sprintf("Failed to interpret 'preferred-date-range' of '%s'", preferred_date_range_str), err, http.StatusBadRequest)
		return
	}
	preferred_time_str := postFormValueOrNone(r, "preferred-time")
	var preferred_time enums.PublicreportNuisancepreferredtimetype
	err = preferred_time.Scan(preferred_time_str)
	if err != nil {
		respondError(w, fmt.Sprintf("Failed to interpret 'preferred-time' of '%s'", preferred_time_str), err, http.StatusBadRequest)
		return
	}

	severity_str := r.PostFormValue("severity")
	severity, err := strconv.ParseInt(severity_str, 10, 16)
	if err != nil {
		respondError(w, fmt.Sprintf("Failed to interpret 'severity' of '%s' as an integer", severity_str), err, http.StatusBadRequest)
		return
	}

	source_description := r.PostFormValue("source-description")
	address := r.PostFormValue("address")
	name := r.PostFormValue("name")
	phone := r.PostFormValue("phone")
	email := r.PostFormValue("email")
	additional_info := r.PostFormValue("additional-info")

	public_id, err := GenerateReportID()
	if err != nil {
		respondError(w, "Failed to create quick report public ID", err, http.StatusInternalServerError)
		return
	}

	log.Info().Str("address", address).Str("name", name).Msg("Got report")
	setter := models.PublicreportNuisanceSetter{
		AdditionalInfo:     omit.From(additional_info),
		Created:            omit.From(time.Now()),
		Duration:           omit.From(duration),
		Email:              omit.From(email),
		InspectionType:     omit.From(inspection_type),
		Location:           omitnull.FromPtr[string](nil),
		PreferredDateRange: omit.From(preferred_date_range),
		PreferredTime:      omit.From(preferred_time),
		PublicID:           omit.From(public_id),
		RequestCall:        omit.From(request_call),
		Severity:           omit.From(int16(severity)),
		SourceContainer:    omit.From(source_container),
		SourceDescription:  omit.From(source_description),
		SourceRoof:         omit.From(source_roof),
		SourceLocation:     omit.From(source_location),
		SourceStagnant:     omit.From(source_stagnant),
		Status:             omit.From(enums.PublicreportReportstatustypeReported),
		TimeOfDayDay:       omit.From(tod_day),
		TimeOfDayEarly:     omit.From(tod_early),
		TimeOfDayEvening:   omit.From(tod_evening),
		TimeOfDayNight:     omit.From(tod_night),
		ReporterAddress:    omit.From(address),
		ReporterEmail:      omit.From(email),
		ReporterName:       omit.From(name),
		ReporterPhone:      omit.From(phone),
	}
	nuisance, err := models.PublicreportNuisances.Insert(&setter).One(r.Context(), db.PGInstance.BobDB)
	if err != nil {
		respondError(w, "Failed to create database record", err, http.StatusInternalServerError)
		return
	}
	log.Info().Str("public_id", public_id).Int32("id", nuisance.ID).Msg("Created nuisance report")
	http.Redirect(w, r, fmt.Sprintf("/nuisance-submit-complete?report=%s", public_id), http.StatusFound)
}
