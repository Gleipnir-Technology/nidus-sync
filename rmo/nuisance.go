package rmo

import (
	"fmt"
	"net/http"
	"strconv"
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

var (
	NuisanceT       = buildTemplate("nuisance", "base")
	SubmitCompleteT = buildTemplate("submit-complete", "base")
)

func getNuisance(w http.ResponseWriter, r *http.Request) {
	html.RenderOrError(
		w,
		NuisanceT,
		ContentNuisance{
			District:    nil,
			MapboxToken: config.MapboxToken,
			URL:         makeContentURL(),
		},
	)
}
func getSubmitComplete(w http.ResponseWriter, r *http.Request) {
	report := r.URL.Query().Get("report")
	html.RenderOrError(
		w,
		SubmitCompleteT,
		ContentNuisanceSubmitComplete{
			District: nil,
			ReportID: report,
			URL:      makeContentURL(),
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
	address := r.PostFormValue("address")
	lat := r.FormValue("latitude")
	lng := r.FormValue("longitude")
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

	latitude, err := strconv.ParseFloat(lat, 64)
	if err != nil {
		respondError(w, "Failed to create parse latitude", err, http.StatusBadRequest)
		return
	}
	longitude, err := strconv.ParseFloat(lng, 64)
	if err != nil {
		respondError(w, "Failed to create parse longitude", err, http.StatusBadRequest)
		return
	}
	public_id, err := report.GenerateReportID()
	if err != nil {
		respondError(w, "Failed to create quick report public ID", err, http.StatusInternalServerError)
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
	organization_id, err := matchDistrict(ctx, longitude, latitude, uploads)
	if err != nil {
		log.Warn().Err(err).Msg("Failed to match district")
	}
	c, err := h3utils.GetCell(longitude, latitude, 15)
	if err != nil {
		respondError(w, "Failedt o get h3 cell", err, http.StatusInternalServerError)
	}

	setter := models.PublicreportNuisanceSetter{
		AdditionalInfo: omit.From(additional_info),
		Address:        omit.From(address),
		Created:        omit.From(time.Now()),
		Duration:       omit.From(duration),
		H3cell:         omitnull.From(c.String()),
		//Location: omitnull.From(fmt.Sprintf("ST_GeometryFromText(Point(%s %s))", longitude, latitude)),
		Location:          omitnull.FromPtr[string](nil),
		OrganizationID:    omitnull.FromPtr(organization_id),
		PublicID:          omit.From(public_id),
		ReporterEmail:     omitnull.FromPtr[string](nil),
		ReporterName:      omitnull.FromPtr[string](nil),
		ReporterPhone:     omitnull.FromPtr[string](nil),
		SourceContainer:   omit.From(source_container),
		SourceDescription: omit.From(source_description),
		SourceGutter:      omit.From(source_gutters),
		SourceLocation:    omit.From(source_location),
		SourceStagnant:    omit.From(source_stagnant),
		Status:            omit.From(enums.PublicreportReportstatustypeReported),
	}
	nuisance, err := models.PublicreportNuisances.Insert(&setter).One(ctx, txn)
	if err != nil {
		respondError(w, "Failed to create database record", err, http.StatusInternalServerError)
		return
	}
	_, err = psql.Update(
		um.Table("publicreport.nuisance"),
		um.SetCol("location").To(fmt.Sprintf("ST_GeometryFromText('Point(%f %f)')", longitude, latitude)),
		um.Where(psql.Quote("id").EQ(psql.Arg(nuisance.ID))),
	).Exec(ctx, txn)
	if err != nil {
		respondError(w, "Failed to insert publicreport", err, http.StatusInternalServerError)
		return
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
