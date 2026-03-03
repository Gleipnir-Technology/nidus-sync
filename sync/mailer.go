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
	"github.com/Gleipnir-Technology/nidus-sync/platform/pdf"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
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

func getMailer(w http.ResponseWriter, r *http.Request) {
	code := chi.URLParam(r, "code")
	if code == "" {
		http.Error(w, "empty code", http.StatusBadRequest)
		return
	}

	content, err := pdf.GeneratePDF(r.Context(), code)
	if err != nil {
		respondError(w, "generate pdf failure", err, http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/pdf")
	disposition := fmt.Sprintf("attachment; filename=\"compliance-mailer-%s.pdf\"", code)
	w.Header().Set("Content-Disposition", disposition)
	_, err = io.Copy(w, bytes.NewReader(content))
	if err != nil {
		respondError(w, "copy error", err, http.StatusInternalServerError)
		return
	}
}
func getMailerPreview(w http.ResponseWriter, r *http.Request) {
	code := chi.URLParam(r, "code")
	if code == "" {
		http.Error(w, "empty code", http.StatusBadRequest)
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
		http.Error(w, "no comp", http.StatusInternalServerError)
		return
	}
	doc_id := uuid.New()
	html.RenderOrError(w, "sync/mailer.html", contentMailer{
		Config:       html.NewContentConfig(),
		DocumentID:   doc_id.String(),
		LogoURL:      config.MakeURLNidus("/api/district/%s/logo", org.Slug.GetOr("unset")),
		Organization: org,
		PoolImageURL: config.MakeURLNidus("/api/compliance-request/image/pool/%s", code),
		QRCodeURL:    config.MakeURLNidus("/qr-code/mailer/%s", code),
		ReportURL:    config.MakeURLReport("/mailer/%s", code),
	})
}
