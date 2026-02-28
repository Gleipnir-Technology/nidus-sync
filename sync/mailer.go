package sync

import (
	"net/http"
	//"strconv"

	"github.com/Gleipnir-Technology/nidus-sync/config"
	"github.com/Gleipnir-Technology/nidus-sync/db"
	"github.com/Gleipnir-Technology/nidus-sync/db/models"
	"github.com/Gleipnir-Technology/nidus-sync/html"
	//"github.com/google/uuid"
	"github.com/go-chi/chi/v5"
)

type contentMailer struct {
	Config       contentConfig
	DocumentID   string
	LogoURL      string
	Organization *models.Organization
	PoolImageURL string
	QRCodeURL    string
	ReportURL    string
}

// func getMailer(ctx context.Context, r *http.Request, org *models.Organization, user *models.User) (*response[contentMailer], *errorWithStatus) {
func getMailer(w http.ResponseWriter, r *http.Request) {
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
	html.RenderOrError(w, "sync/mailer.html", contentMailer{
		Config:       newContentConfig(),
		DocumentID:   "00000000-0000-0000-0000-000000000000",
		LogoURL:      config.MakeURLNidus("/api/district/%s/logo", org.Slug.GetOr("unset")),
		Organization: org,
		PoolImageURL: config.MakeURLNidus("/api/compliance-request/image/pool/%s", code),
		QRCodeURL:    config.MakeURLNidus("/qr-code/mailer/%s", code),
		ReportURL:    config.MakeURLReport("/mailer/%s", code),
	})
}
