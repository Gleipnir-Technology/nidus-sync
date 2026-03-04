package rmo

import (
	"context"
	"net/http"

	"github.com/Gleipnir-Technology/bob"
	"github.com/Gleipnir-Technology/bob/dialect/psql"
	"github.com/Gleipnir-Technology/bob/dialect/psql/sm"
	"github.com/Gleipnir-Technology/nidus-sync/config"
	"github.com/Gleipnir-Technology/nidus-sync/db"
	"github.com/Gleipnir-Technology/nidus-sync/html"
	nhttp "github.com/Gleipnir-Technology/nidus-sync/http"
	"github.com/go-chi/chi/v5"
	"github.com/rs/zerolog/log"
	"github.com/stephenafamo/scan"
	//"github.com/Gleipnir-Technology/nidus-sync/config"
)

type address struct {
	Country          string `db:"country"`
	Locality         string `db:"locality"`
	LocationGeoJSON  string `db:"location_geo_json"`
	Number           int32  `db:"number_"`
	OrganizationSlug string `db:"slug"`
	PostalCode       string `db:"postal_code"`
	Region           string `db:"region"`
	Street           string `db:"street"`
}
type contentMailer struct {
	Address  address
	PublicID string
	URLLogo  string
}

func getMailer(ctx context.Context, r *http.Request) (*html.Response[contentMailer], *nhttp.ErrorWithStatus) {
	public_id := chi.URLParam(r, "public_id")
	if public_id == "" {
		return nil, nhttp.NewErrorStatus(http.StatusBadRequest, "No 'public_id' in the url params")
	}

	/*
		compliance_request, err := models.ComplianceReportRequests.Query(
			models.Preload.ComplianceReportRequest.Site(),
			models.SelectWhere.ComplianceReportRequests.PublicID.EQ(public_id),
		).One(ctx, db.PGInstance.BobDB)
		if err != nil {
			respondError(w, "failed to get compliance request", err, http.StatusBadRequest)
		}
		site := compliance_request.
	*/
	report, err := bob.One(ctx, db.PGInstance.BobDB, psql.Select(
		sm.Columns(
			"address.number_",
			"address.street",
			"address.locality",
			"ST_AsGeoJSON(address.geom) AS location_geo_json",
			"address.region",
			"address.postal_code",
			"address.country",
			"organization.slug",
		),
		sm.From("compliance_report_request").As("crr"),
		sm.InnerJoin("site").OnEQ(psql.Raw("crr.site_id"), psql.Raw("site.id")),
		sm.InnerJoin("organization").OnEQ(psql.Raw("site.organization_id"), psql.Raw("organization.id")),
		sm.InnerJoin("address").OnEQ(psql.Raw("site.address_id"), psql.Raw("address.id")),
		sm.Where(psql.Raw("crr.public_id").EQ(psql.Arg(public_id))),
	), scan.StructMapper[address]())
	if err != nil {
		log.Warn().Err(err).Msg("failed to get compliance report")
		return nil, nhttp.NewErrorStatus(http.StatusNotFound, "No compliance report with that public ID")
	}
	return html.NewResponse(
		"rmo/mailer/root.html", contentMailer{
			Address:  report,
			PublicID: public_id,
			URLLogo:  config.MakeURLNidus("/api/district/%s/logo", report.OrganizationSlug),
		},
	), nil

}
func getMailerConfirm(ctx context.Context, r *http.Request) (*html.Response[contentMailer], *nhttp.ErrorWithStatus) {
	public_id := chi.URLParam(r, "public_id")
	if public_id == "" {
		return nil, nhttp.NewErrorStatus(http.StatusBadRequest, "No 'public_id' in the url params")
	}
	return html.NewResponse(
		"rmo/mailer/confirm.html", contentMailer{
			PublicID: public_id,
		},
	), nil
}
func getMailerContribute(ctx context.Context, r *http.Request) (*html.Response[contentMailer], *nhttp.ErrorWithStatus) {
	public_id := chi.URLParam(r, "public_id")
	if public_id == "" {
		return nil, nhttp.NewErrorStatus(http.StatusBadRequest, "No 'public_id' in the url params")
	}
	return html.NewResponse(
		"rmo/mailer/contribute.html", contentMailer{
			PublicID: public_id,
		},
	), nil
}
func getMailerEvidence(ctx context.Context, r *http.Request) (*html.Response[contentMailer], *nhttp.ErrorWithStatus) {
	public_id := chi.URLParam(r, "public_id")
	if public_id == "" {
		return nil, nhttp.NewErrorStatus(http.StatusBadRequest, "No 'public_id' in the url params")
	}
	return html.NewResponse(
		"rmo/mailer/evidence.html", contentMailer{
			PublicID: public_id,
		},
	), nil
}
func getMailerSchedule(ctx context.Context, r *http.Request) (*html.Response[contentMailer], *nhttp.ErrorWithStatus) {
	public_id := chi.URLParam(r, "public_id")
	if public_id == "" {
		return nil, nhttp.NewErrorStatus(http.StatusBadRequest, "No 'public_id' in the url params")
	}
	return html.NewResponse(
		"rmo/mailer/schedule.html", contentMailer{
			PublicID: public_id,
		},
	), nil
}
func getMailerUpdate(ctx context.Context, r *http.Request) (*html.Response[contentMailer], *nhttp.ErrorWithStatus) {
	public_id := chi.URLParam(r, "public_id")
	if public_id == "" {
		return nil, nhttp.NewErrorStatus(http.StatusBadRequest, "No 'public_id' in the url params")
	}
	return html.NewResponse(
		"rmo/mailer/update.html", contentMailer{
			PublicID: public_id,
		},
	), nil
}

type formMailerConfirm struct{}

func postMailerConfirm(ctx context.Context, r *http.Request, form formMailerConfirm) (string, *nhttp.ErrorWithStatus) {
	log.Info().Msg("Fake confirm location")
	public_id := chi.URLParam(r, "public_id")
	if public_id == "" {
		return "", nhttp.NewErrorStatus(http.StatusBadRequest, "No 'public_id' in the url params")
	}
	return config.MakeURLReport("/mailer/%s/evidence", public_id), nil
}
