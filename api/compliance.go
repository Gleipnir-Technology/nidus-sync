package api

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"

	"github.com/Gleipnir-Technology/nidus-sync/db"
	"github.com/Gleipnir-Technology/nidus-sync/db/models"
	"github.com/Gleipnir-Technology/nidus-sync/platform"
	"github.com/Gleipnir-Technology/nidus-sync/platform/imagetile"
	"github.com/go-chi/chi/v5"
	"github.com/rs/zerolog/log"
)

func getComplianceRequestImagePool(w http.ResponseWriter, r *http.Request) {
	code := chi.URLParam(r, "public_id")
	if code == "" {
		http.Error(w, "empty public_id", http.StatusBadRequest)
		return
	}

	ctx := r.Context()
	comp, err := models.ComplianceReportRequests.Query(
		models.Preload.ComplianceReportRequest.Site(),
		models.SelectWhere.ComplianceReportRequests.PublicID.EQ(code),
	).One(ctx, db.PGInstance.BobDB)
	if err != nil {
		http.Error(w, "no comp", http.StatusInternalServerError)
		return
	}

	site := comp.R.Site
	org, err := models.FindOrganization(ctx, db.PGInstance.BobDB, site.OrganizationID)
	if err != nil {
		http.Error(w, "no org", http.StatusInternalServerError)
		return
	}
	envelope, err := platform.ParcelEnvelope(ctx, site.ParcelID)
	if err != nil {
		log.Error().Err(err).Msg("parcel envelop failure")
		http.Error(w, "parcel env", http.StatusInternalServerError)
		return
	}
	log.Info().Int("len", len(*envelope)).Msg("got envelope")
	ring := (*envelope)[0]
	log.Info().Int("len", len(ring)).Msg("got ring")
	p := ring[0]
	log.Info().Int("len", len(p)).Msg("got point")
	writeImage(ctx, w, org, 14, p[0], p[1])
}
func writeImage(ctx context.Context, w http.ResponseWriter, org *models.Organization, level uint, x, y float64) {
	img, err := imagetile.ImageAtPoint(ctx, org, level, x, y)
	w.Header().Set("Content-Type", "image/png")
	w.Header().Set("Content-Length", fmt.Sprintf("%d", len(img)))
	_, err = io.Copy(w, bytes.NewBuffer(img))
	if err != nil {
		http.Error(w, "failed copy", http.StatusInternalServerError)
		return
	}
}
