package rmo

import (
	"net/http"

	"github.com/Gleipnir-Technology/nidus-sync/config"
	"github.com/Gleipnir-Technology/nidus-sync/db"
	"github.com/Gleipnir-Technology/nidus-sync/db/models"
	"github.com/go-chi/chi/v5"
)

type ContentDistrict struct {
	Name       string
	URLLogo    string
	URLWebsite string
}

func districtBySlug(r *http.Request) (*models.Organization, error) {
	slug := chi.URLParam(r, "slug")
	district, err := models.Organizations.Query(
		models.SelectWhere.Organizations.Slug.EQ(slug),
	).One(r.Context(), db.PGInstance.BobDB)
	return district, err
}
func newContentDistrict(d *models.Organization) *ContentDistrict {
	if d == nil {
		return nil
	}
	return &ContentDistrict{
		Name:       d.Name,
		URLLogo:    config.MakeURLNidus("/api/district/%s/logo", d.Slug.GetOr("unset")),
		URLWebsite: d.Website.GetOr(""),
	}
}
