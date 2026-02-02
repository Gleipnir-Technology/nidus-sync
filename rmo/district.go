package rmo

import (
	"net/http"

	"github.com/Gleipnir-Technology/bob/dialect/psql/sm"
	"github.com/Gleipnir-Technology/nidus-sync/config"
	"github.com/Gleipnir-Technology/nidus-sync/db"
	"github.com/Gleipnir-Technology/nidus-sync/db/models"
	"github.com/Gleipnir-Technology/nidus-sync/html"
	"github.com/go-chi/chi/v5"
)

type ContentDistrict struct {
	Name       string
	URLLogo    string
	URLRMO     string
	URLWebsite string
}
type ContentDistrictList struct {
	Districts []ContentDistrict
	URL       ContentURL
}

var (
	DistrictListT = buildTemplate("district-list", "base")
)

func districtBySlug(r *http.Request) (*models.Organization, error) {
	slug := chi.URLParam(r, "slug")
	district, err := models.Organizations.Query(
		models.SelectWhere.Organizations.Slug.EQ(slug),
	).One(r.Context(), db.PGInstance.BobDB)
	return district, err
}
func getDistrictList(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	rows, err := models.Organizations.Query(
		models.SelectWhere.Organizations.ImportDistrictGid.IsNotNull(),
		sm.OrderBy("name"),
	).All(ctx, db.PGInstance.BobDB)
	if err != nil {
		respondError(w, "failed to query for districts", err, http.StatusInternalServerError)
		return
	}
	districts := make([]ContentDistrict, 0)
	for _, row := range rows {
		districts = append(districts, *newContentDistrict(row))
	}
	html.RenderOrError(
		w,
		DistrictListT,
		ContentDistrictList{
			Districts: districts,
			URL:       makeContentURL(nil),
		},
	)

}
func newContentDistrict(d *models.Organization) *ContentDistrict {
	if d == nil {
		return nil
	}
	return &ContentDistrict{
		Name:       d.Name,
		URLLogo:    config.MakeURLNidus("/api/district/%s/logo", d.Slug.GetOr("unset")),
		URLRMO:     config.MakeURLReport("/district/%s", d.Slug.GetOr("unset")),
		URLWebsite: d.Website.GetOr(""),
	}
}
