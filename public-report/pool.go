package publicreport

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/Gleipnir-Technology/nidus-sync/config"
	"github.com/Gleipnir-Technology/nidus-sync/db"
	"github.com/Gleipnir-Technology/nidus-sync/db/enums"
	"github.com/Gleipnir-Technology/nidus-sync/db/models"
	"github.com/Gleipnir-Technology/nidus-sync/htmlpage"
	"github.com/aarondl/opt/omit"
	"github.com/rs/zerolog/log"
	"github.com/stephenafamo/bob/dialect/psql"
	"github.com/stephenafamo/bob/dialect/psql/um"
)

type ContextPool struct {
	MapboxToken string
}
type ContextPoolSubmitComplete struct {
	ReportID string
}

var (
	Pool               = buildTemplate("pool", "base")
	PoolSubmitComplete = buildTemplate("pool-submit-complete", "base")
)

func getPool(w http.ResponseWriter, r *http.Request) {
	htmlpage.RenderOrError(
		w,
		Pool,
		ContextPool{
			MapboxToken: config.MapboxToken,
		},
	)
}
func getPoolSubmitComplete(w http.ResponseWriter, r *http.Request) {
	report := r.URL.Query().Get("report")
	htmlpage.RenderOrError(
		w,
		PoolSubmitComplete,
		ContextPoolSubmitComplete{
			ReportID: report,
		},
	)
}
func postPool(w http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(32 << 10) // 32 MB buffer
	if err != nil {
		respondError(w, "Failed to parse form", err, http.StatusBadRequest)
		return
	}
	access_comments := r.FormValue("access-comments")
	access_gate := boolFromForm(r, "access-gate")
	access_fence := boolFromForm(r, "access-fence")
	access_locked := boolFromForm(r, "access-locked")
	access_dog := boolFromForm(r, "access-dog")
	access_other := boolFromForm(r, "access-other")
	address := r.FormValue("address")
	address_country := r.FormValue("address-country")
	address_postcode := r.FormValue("address-postcode")
	address_place := r.FormValue("address-place")
	address_region := r.FormValue("address-region")
	address_street := r.FormValue("address-street")
	comments := r.FormValue("comments")
	has_adult := boolFromForm(r, "has-adult")
	has_larvae := boolFromForm(r, "has-larvae")
	has_pupae := boolFromForm(r, "has-pupae")
	map_zoom_str := r.FormValue("map-zoom")
	owner_email := r.FormValue("owner-email")
	owner_name := r.FormValue("owner-name")
	owner_phone := r.FormValue("owner-phone")
	reporter_email := r.FormValue("reporter-email")
	reporter_name := r.FormValue("reporter-name")
	reporter_phone := r.FormValue("reporter-phone")
	subscribe := boolFromForm(r, "subscribe")

	map_zoom, err := strconv.ParseFloat(map_zoom_str, 32)
	if err != nil {
		respondError(w, "Failed to parse zoom level", err, http.StatusBadRequest)
		return
	}
	public_id, err := GenerateReportID()
	if err != nil {
		respondError(w, "Failed to create pool report public ID", err, http.StatusInternalServerError)
		return
	}

	setter := models.PublicreportPoolSetter{
		AccessComments:  omit.From(access_comments),
		AccessGate:      omit.From(access_gate),
		AccessFence:     omit.From(access_fence),
		AccessLocked:    omit.From(access_locked),
		AccessDog:       omit.From(access_dog),
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
		HasAdult:  omit.From(has_adult),
		HasLarvae: omit.From(has_larvae),
		HasPupae:  omit.From(has_pupae),
		//Location: add later
		MapZoom:       omit.From(map_zoom),
		OwnerEmail:    omit.From(owner_email),
		OwnerName:     omit.From(owner_name),
		OwnerPhone:    omit.From(owner_phone),
		PublicID:      omit.From(public_id),
		ReporterEmail: omit.From(reporter_email),
		ReporterName:  omit.From(reporter_name),
		ReporterPhone: omit.From(reporter_phone),
		Status:        omit.From(enums.PublicreportReportstatustypeReported),
		Subscribe:     omit.From(subscribe),
	}
	pool, err := models.PublicreportPools.Insert(&setter).One(r.Context(), db.PGInstance.BobDB)
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
		).Exec(r.Context(), db.PGInstance.BobDB)
		if err != nil {
			respondError(w, "Failed to insert publicreport.pool", err, http.StatusInternalServerError)
			return
		}
	}
	log.Info().Int32("id", pool.ID).Str("public_id", pool.PublicID).Msg("Created pool report")
	photoSetters := make([]*models.PublicreportPoolPhotoSetter, 0)
	uploads, err := extractPhotoUploads(r)
	if err != nil {
		respondError(w, "Failed to extract photo uploads", err, http.StatusInternalServerError)
		return
	}
	for _, u := range uploads {
		photoSetters = append(photoSetters, &models.PublicreportPoolPhotoSetter{
			Filename: omit.From(u.Filename),
			Size:     omit.From(u.Size),
			UUID:     omit.From(u.UUID),
		})
	}
	err = pool.InsertPoolPhotos(r.Context(), db.PGInstance.BobDB, photoSetters...)
	if err != nil {
		respondError(w, "Failed to create photo records", err, http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, fmt.Sprintf("/pool-submit-complete?report=%s", public_id), http.StatusFound)
}
