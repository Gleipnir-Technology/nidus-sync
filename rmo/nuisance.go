package rmo

import (
	"fmt"
	"net/http"
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
	address := r.PostFormValue("address")
	source_stagnant := boolFromForm(r, "source-stagnant")
	source_container := boolFromForm(r, "source-container")
	source_gutters := boolFromForm(r, "source-gutters")

	duration_str := postFormValueOrNone(r, "duration")
	var duration enums.PublicreportNuisancedurationtype
	err = duration.Scan(duration_str)
	if err != nil {
		respondError(w, fmt.Sprintf("Failed to interpret 'duration' of '%s'", duration_str), err, http.StatusBadRequest)
		return
	}

	source_location_str := postFormValueOrNone(r, "source-location")
	var source_location enums.PublicreportNuisancelocationtype
	err = source_location.Scan(source_location_str)
	if err != nil {
		respondError(w, fmt.Sprintf("Failed to interpret 'source-location' of '%s'", source_location_str), err, http.StatusBadRequest)
		return
	}

	source_description := r.PostFormValue("source-description")
	additional_info := r.PostFormValue("additional-info")

	public_id, err := GenerateReportID()
	if err != nil {
		respondError(w, "Failed to create quick report public ID", err, http.StatusInternalServerError)
		return
	}

	setter := models.PublicreportNuisanceSetter{
		AdditionalInfo:    omit.From(additional_info),
		Address:           omit.From(address),
		Created:           omit.From(time.Now()),
		Duration:          omit.From(duration),
		Location:          omitnull.FromPtr[string](nil),
		PublicID:          omit.From(public_id),
		SourceContainer:   omit.From(source_container),
		SourceDescription: omit.From(source_description),
		SourceGutter:      omit.From(source_gutters),
		SourceLocation:    omit.From(source_location),
		SourceStagnant:    omit.From(source_stagnant),
		Status:            omit.From(enums.PublicreportReportstatustypeReported),
		ReporterEmail:     omitnull.FromPtr[string](nil),
		ReporterName:      omitnull.FromPtr[string](nil),
		ReporterPhone:     omitnull.FromPtr[string](nil),
	}
	nuisance, err := models.PublicreportNuisances.Insert(&setter).One(r.Context(), db.PGInstance.BobDB)
	if err != nil {
		respondError(w, "Failed to create database record", err, http.StatusInternalServerError)
		return
	}
	log.Info().Str("public_id", public_id).Int32("id", nuisance.ID).Msg("Created nuisance report")
	http.Redirect(w, r, fmt.Sprintf("/nuisance-submit-complete?report=%s", public_id), http.StatusFound)
}
