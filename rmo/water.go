package rmo

import (
	"fmt"
	"net/http"
	"time"

	"github.com/Gleipnir-Technology/bob"
	"github.com/Gleipnir-Technology/bob/dialect/psql"
	"github.com/Gleipnir-Technology/bob/dialect/psql/um"
	"github.com/Gleipnir-Technology/nidus-sync/config"
	"github.com/Gleipnir-Technology/nidus-sync/db"
	"github.com/Gleipnir-Technology/nidus-sync/db/enums"
	"github.com/Gleipnir-Technology/nidus-sync/db/models"
	"github.com/Gleipnir-Technology/nidus-sync/html"
	"github.com/Gleipnir-Technology/nidus-sync/platform/report"
	"github.com/aarondl/opt/omit"
	"github.com/rs/zerolog/log"
)

type ContentPool struct {
	District    *ContentDistrict
	MapboxToken string
	URL         ContentURL
}

func getWater(w http.ResponseWriter, r *http.Request) {
	html.RenderOrError(
		w,
		"rmo/water.html",
		ContentPool{
			District:    nil,
			MapboxToken: config.MapboxToken,
			URL:         makeContentURL(nil),
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
		ContentPool{
			District:    newContentDistrict(district),
			MapboxToken: config.MapboxToken,
			URL:         makeContentURL(district),
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
	address := r.FormValue("address")
	address_country := r.FormValue("address-country")
	address_postcode := r.FormValue("address-postcode")
	address_place := r.FormValue("address-place")
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
		respondError(w, "Failed to parse lat lng for pool report", err, http.StatusInternalServerError)
		return
	}
	public_id, err := report.GenerateReportID()
	if err != nil {
		respondError(w, "Failed to create pool report public ID", err, http.StatusInternalServerError)
		return
	}

	ctx := r.Context()
	tx, err := db.PGInstance.BobDB.BeginTx(ctx, nil)
	if err != nil {
		respondError(w, "Failed to create transaction", err, http.StatusInternalServerError)
		return
	}
	defer tx.Rollback(ctx)

	setter := models.PublicreportPoolSetter{
		AccessComments:  omit.From(access_comments),
		AccessDog:       omit.From(access_dog),
		AccessFence:     omit.From(access_fence),
		AccessGate:      omit.From(access_gate),
		AccessLocked:    omit.From(access_locked),
		AccessOther:     omit.From(access_other),
		Address:         omit.From(address),
		AddressCountry:  omit.From(address_country),
		AddressPostCode: omit.From(address_postcode),
		AddressPlace:    omit.From(address_place),
		AddressStreet:   omit.From(address_street),
		AddressRegion:   omit.From(address_region),
		Comments:        omit.From(comments),
		Created:         omit.From(time.Now()),
		//H3cell: add later
		HasAdult:               omit.From(has_adult),
		HasBackyardPermission:  omit.From(has_backyard_permission),
		HasLarvae:              omit.From(has_larvae),
		HasPupae:               omit.From(has_pupae),
		IsReporterConfidential: omit.From(is_reporter_confidential),
		IsReporterOwner:        omit.From(is_reporter_owner),
		//Location: add later
		MapZoom:       omit.From(latlng.MapZoom),
		OwnerEmail:    omit.From(owner_email),
		OwnerName:     omit.From(owner_name),
		OwnerPhone:    omit.From(owner_phone),
		PublicID:      omit.From(public_id),
		ReporterEmail: omit.From(""),
		ReporterName:  omit.From(""),
		ReporterPhone: omit.From(""),
		Status:        omit.From(enums.PublicreportReportstatustypeReported),
	}
	pool, err := models.PublicreportPools.Insert(&setter).One(ctx, tx)
	if err != nil {
		respondError(w, "Failed to create database record", err, http.StatusInternalServerError)
		return
	}

	geospatial, err := geospatialFromForm(r)
	if err != nil {
		respondError(w, "Failed to handle geospatial data", err, http.StatusInternalServerError)
		return
	}
	if geospatial.Populated {
		_, err = psql.Update(
			um.Table("publicreport.pool"),
			um.SetCol("h3cell").ToArg(geospatial.Cell),
			um.SetCol("location").To(geospatial.GeometryQuery),
			um.Where(psql.Quote("id").EQ(psql.Arg(pool.ID))),
		).Exec(ctx, tx)
		if err != nil {
			respondError(w, "Failed to insert publicreport.pool", err, http.StatusInternalServerError)
			return
		}
	}
	log.Info().Int32("id", pool.ID).Str("public_id", pool.PublicID).Msg("Created pool report")
	uploads, err := extractImageUploads(r)
	if err != nil {
		respondError(w, "Failed to extract image uploads", err, http.StatusInternalServerError)
		return
	}
	images, err := saveImageUploads(r.Context(), tx, uploads)
	setters := make([]*models.PublicreportPoolImageSetter, 0)
	for _, image := range images {
		setters = append(setters, &models.PublicreportPoolImageSetter{
			ImageID: omit.From(int32(image.ID)),
			PoolID:  omit.From(int32(pool.ID)),
		})
	}
	if len(setters) > 0 {
		_, err = models.PublicreportPoolImages.Insert(bob.ToMods(setters...)).Exec(r.Context(), tx)
		if err != nil {
			respondError(w, "Failed to save upload relationships", err, http.StatusInternalServerError)
			return
		}
	}
	tx.Commit(ctx)
	http.Redirect(w, r, fmt.Sprintf("/submit-complete?report=%s", public_id), http.StatusFound)
}
func postWaterDistrict(w http.ResponseWriter, r *http.Request) {
}
