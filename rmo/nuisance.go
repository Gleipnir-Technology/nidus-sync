package rmo

import (
	"fmt"
	"net/http"
	"slices"
	"time"

	"github.com/Gleipnir-Technology/bob"
	"github.com/Gleipnir-Technology/bob/dialect/psql"
	"github.com/Gleipnir-Technology/bob/dialect/psql/um"
	"github.com/Gleipnir-Technology/nidus-sync/config"
	"github.com/Gleipnir-Technology/nidus-sync/db"
	"github.com/Gleipnir-Technology/nidus-sync/db/enums"
	"github.com/Gleipnir-Technology/nidus-sync/db/models"
	"github.com/Gleipnir-Technology/nidus-sync/h3utils"
	"github.com/Gleipnir-Technology/nidus-sync/html"
	"github.com/Gleipnir-Technology/nidus-sync/platform/report"
	"github.com/aarondl/opt/omit"
	"github.com/aarondl/opt/omitnull"
	"github.com/rs/zerolog/log"
	"github.com/uber/h3-go/v4"
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
			District:    nil,
			MapboxToken: config.MapboxToken,
			URL:         makeContentURL(nil),
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
			District:    newContentDistrict(district),
			MapboxToken: config.MapboxToken,
			URL:         makeContentURL(nil),
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
	address := r.PostFormValue("address")
	address_country := r.PostFormValue("address-country")
	address_place := r.PostFormValue("address-place")
	address_postcode := r.PostFormValue("address-postcode")
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
	//log.Debug().Bool("is_location_backyard", is_location_backyard).Bool("is_location_frontyard", is_location_frontyard).Bool("is_location_garden", is_location_garden).Bool("is_location_other", is_location_other).Bool("is_location_pool", is_location_pool).Msg("parsed")

	public_id, err := report.GenerateReportID()
	if err != nil {
		respondError(w, "Failed to create report public ID", err, http.StatusInternalServerError)
		return
	}

	txn, err := db.PGInstance.BobDB.BeginTx(ctx, nil)
	if err != nil {
		respondError(w, "Failed to create transaction", err, http.StatusInternalServerError)
		return
	}
	defer txn.Rollback(ctx)

	uploads, err := extractImageUploads(r)
	log.Info().Int("len", len(uploads)).Msg("extracted uploads")
	if err != nil {
		respondError(w, "Failed to extract image uploads", err, http.StatusInternalServerError)
		return
	}
	images, err := saveImageUploads(ctx, txn, uploads)
	if err != nil {
		respondError(w, "Failed to save image uploads", err, http.StatusInternalServerError)
		return
	}
	var organization_id *int32
	var h3cell h3.Cell
	if latlng.Latitude != nil && latlng.Longitude != nil {
		organization_id, err = matchDistrict(ctx, *latlng.Longitude, *latlng.Latitude, uploads)
		if err != nil {
			log.Warn().Err(err).Msg("Failed to match district")
		}
		h3cell, err = h3utils.GetCell(*latlng.Longitude, *latlng.Latitude, 15)
		if err != nil {
			respondError(w, "Failedt o get h3 cell", err, http.StatusInternalServerError)
		}
	}

	setter := models.PublicreportNuisanceSetter{
		AdditionalInfo:      omit.From(additional_info),
		Address:             omit.From(address),
		AddressCountry:      omit.From(address_country),
		AddressPlace:        omit.From(address_place),
		AddressPostcode:     omit.From(address_postcode),
		AddressRegion:       omit.From(address_region),
		AddressStreet:       omit.From(address_street),
		Created:             omit.From(time.Now()),
		Duration:            omit.From(duration),
		H3cell:              omitnull.From(h3cell.String()),
		IsLocationBackyard:  omit.From(is_location_backyard),
		IsLocationFrontyard: omit.From(is_location_frontyard),
		IsLocationGarden:    omit.From(is_location_garden),
		IsLocationOther:     omit.From(is_location_other),
		IsLocationPool:      omit.From(is_location_pool),
		LatlngAccuracyType:  omit.From(latlng.AccuracyType),
		LatlngAccuracyValue: omit.From(latlng.AccuracyValue),
		//Location: omitnull.From(fmt.Sprintf("ST_GeometryFromText(Point(%s %s))", longitude, latitude)),
		Location:          omitnull.FromPtr[string](nil),
		MapZoom:           omit.From(latlng.MapZoom),
		OrganizationID:    omitnull.FromPtr(organization_id),
		PublicID:          omit.From(public_id),
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
	nuisance, err := models.PublicreportNuisances.Insert(&setter).One(ctx, txn)
	if err != nil {
		respondError(w, "Failed to create database record", err, http.StatusInternalServerError)
		return
	}
	if latlng.Latitude != nil && latlng.Longitude != nil {
		_, err = psql.Update(
			um.Table("publicreport.nuisance"),
			um.SetCol("location").To(fmt.Sprintf("ST_GeometryFromText('Point(%f %f)')", *latlng.Longitude, *latlng.Latitude)),
			um.Where(psql.Quote("id").EQ(psql.Arg(nuisance.ID))),
		).Exec(ctx, txn)
		if err != nil {
			respondError(w, "Failed to insert publicreport", err, http.StatusInternalServerError)
			return
		}
	}
	log.Info().Str("public_id", public_id).Int32("id", nuisance.ID).Msg("Created nuisance report")
	if len(images) > 0 {
		setters := make([]*models.PublicreportNuisanceImageSetter, 0)
		for _, image := range images {
			setters = append(setters, &models.PublicreportNuisanceImageSetter{
				ImageID:    omit.From(int32(image.ID)),
				NuisanceID: omit.From(int32(nuisance.ID)),
			})
		}
		_, err = models.PublicreportNuisanceImages.Insert(bob.ToMods(setters...)).Exec(ctx, txn)
		if err != nil {
			respondError(w, "Failed to save reference to images", err, http.StatusInternalServerError)
			return
		}
		log.Info().Int("len", len(images)).Msg("saved uploads")
	}
	txn.Commit(ctx)
	http.Redirect(w, r, fmt.Sprintf("/submit-complete?report=%s", public_id), http.StatusFound)
}
