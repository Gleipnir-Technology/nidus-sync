package sync

import (
	"bytes"
	"fmt"
	"io"
	"net/http"

	"github.com/Gleipnir-Technology/nidus-sync/config"
	"github.com/Gleipnir-Technology/nidus-sync/db"
	"github.com/Gleipnir-Technology/nidus-sync/db/models"
	"github.com/Gleipnir-Technology/nidus-sync/html"
	"github.com/Gleipnir-Technology/nidus-sync/platform"
	"github.com/Gleipnir-Technology/nidus-sync/platform/pdf"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/rs/zerolog/log"
)

type contentMailer struct {
	Config       html.ContentConfig
	DocumentID   string
	LogoURL      string
	Organization *models.Organization
	PoolImageURL string
	QRCodeURL    string
	ReportURL    string
}

func getMailer1(w http.ResponseWriter, r *http.Request) {
	path := "/mailer/mode-1/preview"
	content, err := pdf.GeneratePDF(r.Context(), path)
	if err != nil {
		respondError(w, "generate pdf failure", err, http.StatusInternalServerError)
		return
	}
	err = writePDF(w, content, "mailer-mode-1.pdf")
	if err != nil {
		respondError(w, "copy error", err, http.StatusInternalServerError)
		return
	}
}
func getMailer2(w http.ResponseWriter, r *http.Request) {
	path := "/mailer/mode-2/preview"
	content, err := pdf.GeneratePDF(r.Context(), path)
	if err != nil {
		respondError(w, "generate pdf failure", err, http.StatusInternalServerError)
		return
	}
	err = writePDF(w, content, "mailer-mode-2.pdf")
	if err != nil {
		respondError(w, "copy error", err, http.StatusInternalServerError)
		return
	}
}
func getMailer3(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	code := vars["code"]
	if code == "" {
		http.Error(w, "empty code", http.StatusBadRequest)
		return
	}
	path := fmt.Sprintf("/mailer/mode-3/%s/preview", code)
	content, err := pdf.GeneratePDF(r.Context(), path)
	if err != nil {
		respondError(w, "generate pdf failure", err, http.StatusInternalServerError)
		return
	}
	filename := fmt.Sprintf("compliance-mailer-%s.pdf", code)
	err = writePDF(w, content, filename)
	if err != nil {
		respondError(w, "copy error", err, http.StatusInternalServerError)
		return
	}
}
func writePDF(w http.ResponseWriter, content []byte, filename string) error {
	w.Header().Set("Content-Type", "application/pdf")
	disposition := fmt.Sprintf("attachment; filename=\"%s\"", filename)
	w.Header().Set("Content-Disposition", disposition)
	_, err := io.Copy(w, bytes.NewReader(content))
	return err
}
func getMailer1Preview(w http.ResponseWriter, r *http.Request) {
	html.RenderOrError(w, "sync/mailer-1.html", contentMailer{})
}
func getMailer2Preview(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	org, err := models.FindOrganization(ctx, db.PGInstance.BobDB, 1)
	//org, err := platform.OrganizationByID(ctx, 1)
	if err != nil {
		http.Error(w, "no comp", http.StatusInternalServerError)
		return
	}

	html.RenderOrError(w, "sync/mailer-2.html", contentMailer{
		Config:       html.NewContentConfig(),
		DocumentID:   "abc-123",
		LogoURL:      config.MakeURLNidus("/api/district/delta-mvcd/logo"),
		Organization: org,
		PoolImageURL: config.MakeURLNidus("/mailer/pool/random"),
		QRCodeURL:    config.MakeURLNidus("/api/qr-code/marketing"),
		ReportURL:    "https://nidus.cloud",
	})
}
func getMailer3Preview(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	code := vars["code"]
	if code == "" {
		http.Error(w, "empty code", http.StatusBadRequest)
		return
	}

	ctx := r.Context()
	comp, err := models.ComplianceReportRequests.Query(
		models.Preload.ComplianceReportRequest.Lead(),
		models.SelectWhere.ComplianceReportRequests.PublicID.EQ(code),
	).One(ctx, db.PGInstance.BobDB)

	if err != nil {
		http.Error(w, "no comp", http.StatusInternalServerError)
		return
	}
	lead := comp.R.Lead
	org, err := models.FindOrganization(ctx, db.PGInstance.BobDB, lead.OrganizationID)
	if err != nil {
		http.Error(w, "no comp", http.StatusInternalServerError)
		return
	}
	doc_id := uuid.New()
	html.RenderOrError(w, "sync/mailer-3.html", contentMailer{
		Config:       html.NewContentConfig(),
		DocumentID:   doc_id.String(),
		LogoURL:      config.MakeURLNidus("/api/district/%s/logo", org.Slug.GetOr("unset")),
		Organization: org,
		PoolImageURL: config.MakeURLNidus("/api/compliance-request/image/pool/%s", code),
		QRCodeURL:    config.MakeURLNidus("/api/qr-code/mailer/%s", code),
		ReportURL:    config.MakeURLReport("/mailer/%s", code),
	})
}
func getMailerPoolRandom(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	err := platform.WriteTileRandom(ctx, w)
	if err != nil {
		log.Error().Err(err).Msg("failed to do random tile")
		http.Error(w, "failed to do tile", http.StatusInternalServerError)
		return
	}
}
