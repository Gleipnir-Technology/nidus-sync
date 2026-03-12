package api

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/Gleipnir-Technology/bob"
	"github.com/Gleipnir-Technology/bob/dialect/psql"
	"github.com/Gleipnir-Technology/bob/dialect/psql/sm"
	"github.com/Gleipnir-Technology/nidus-sync/db"
	"github.com/Gleipnir-Technology/nidus-sync/platform"
	"github.com/go-chi/chi/v5"
	"github.com/paulmach/orb/geojson"
	"github.com/rs/zerolog/log"
	"github.com/stephenafamo/scan"
)

func getComplianceRequestImagePool(w http.ResponseWriter, r *http.Request) {
	code := chi.URLParam(r, "public_id")
	if code == "" {
		http.Error(w, "empty public_id", http.StatusBadRequest)
		return
	}

	ctx := r.Context()
	/*
		comp, err := models.ComplianceReportRequests.Query(
			models.Preload.ComplianceReportRequest.Lead(),
			models.SelectWhere.ComplianceReportRequests.PublicID.EQ(code),
		).One(ctx, db.PGInstance.BobDB)
		if err != nil {
			http.Error(w, "no comp", http.StatusInternalServerError)
			return
		}

		lead := comp.R.Lead
		site := lead.R.Site
	*/
	type _Row struct {
		Envelope       string `db:"parcel_envelope"`
		OrganizationID int32  `db:"organization_id"`
	}
	row, err := bob.One(ctx, db.PGInstance.BobDB, psql.Select(
		sm.Columns(
			"ST_AsGeoJSON(ST_Envelope(parcel.geometry)) AS parcel_envelope",
			"organization.id AS organization_id",
		),
		sm.From("compliance_report_request"),
		sm.InnerJoin("lead").OnEQ(
			psql.Quote("compliance_report_request.lead_id"),
			psql.Quote("organization.id"),
		),
		sm.InnerJoin("organization").OnEQ(
			psql.Quote("lead.organization_id"),
			psql.Quote("organization.id"),
		),
		sm.InnerJoin("site").On(
			psql.Quote("lead.site_id").EQ(psql.Quote("site.id")),
		),
		sm.InnerJoin("parcel").OnEQ(
			psql.Quote("site.parcel_id"),
			psql.Quote("parcel.id"),
		),
		sm.Where(psql.Quote("compliance_report_request").EQ(psql.Arg(code))),
	), scan.StructMapper[_Row]())
	org, err := platform.OrganizationByID(ctx, int(row.OrganizationID))
	if err != nil {
		http.Error(w, "org err", http.StatusInternalServerError)
		return
	}
	if org == nil {
		http.Error(w, "no org", http.StatusBadRequest)
		return
	}
	var polygon geojson.Polygon
	err = json.Unmarshal([]byte(row.Envelope), &polygon)
	if err != nil {
		log.Error().Err(err).Msg("unmarshal json")
		http.Error(w, "unmarshal envelope json", http.StatusInternalServerError)
		return
	}
	ring := polygon[0]
	p := ring[0]
	err = writeImage(ctx, w, *org, 19, p[1], p[0])
	if err != nil {
		log.Error().Err(err).Msg("write image")
		http.Error(w, "failed to write image", http.StatusInternalServerError)
		return
	}
}
func writeImage(ctx context.Context, w http.ResponseWriter, org platform.Organization, level uint, lat, lng float64) error {
	img, err := platform.ImageAtPoint(ctx, org, level, lat, lng)
	if err != nil {
		return fmt.Errorf("image at point: %w", err)
	}
	log.Info().Int("size", len(img.Content)).Msg("image")
	w.Header().Set("Content-Type", "image/png")
	w.Header().Set("Content-Length", fmt.Sprintf("%d", len(img.Content)))
	_, err = io.Copy(w, bytes.NewBuffer(img.Content))
	if err != nil {
		return fmt.Errorf("copy bytes: %w", err)
	}
	return nil
}
