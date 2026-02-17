package rmo

import (
	"context"
	"fmt"
	"net/http"

	"github.com/Gleipnir-Technology/bob/dialect/psql/sm"
	"github.com/Gleipnir-Technology/nidus-sync/config"
	"github.com/Gleipnir-Technology/nidus-sync/db"
	"github.com/Gleipnir-Technology/nidus-sync/db/models"
	"github.com/Gleipnir-Technology/nidus-sync/html"
	"github.com/Gleipnir-Technology/nidus-sync/platform"
	"github.com/go-chi/chi/v5"
	"github.com/rs/zerolog/log"
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
		"rmo/district-list.html",
		ContentDistrictList{
			Districts: districts,
			URL:       makeContentURL(nil),
		},
	)

}
func matchDistrict(ctx context.Context, longitude, latitude *float64, images []ImageUpload) (*int32, error) {
	var err error
	var org *models.Organization
	for _, image := range images {
		if image.Exif.GPS == nil {
			continue
		}
		org, err = platform.DistrictForLocation(ctx, image.Exif.GPS.Longitude, image.Exif.GPS.Latitude)
		if err != nil {
			log.Warn().Err(err).Msg("Failed to get district for location")
			continue
		}
		if org != nil {
			return &org.ID, nil
		}
	}
	if longitude == nil || latitude == nil {
		log.Debug().Msg("No location from images, no latlng for the report itself, cannot match")
		return nil, nil
	}
	org, err = platform.DistrictForLocation(ctx, *longitude, *latitude)
	if err != nil {
		log.Warn().Err(err).Msg("Failed to get district for location")
		return nil, fmt.Errorf("Failed to get district for location: %w", err)
	}
	if org == nil {
		log.Debug().Err(err).Float64("lng", *longitude).Float64("lat", *latitude).Msg("No district match by report location")
		return nil, nil
	}
	log.Debug().Err(err).Int32("org_id", org.ID).Float64("lng", *longitude).Float64("lat", *latitude).Msg("Found district match by report location")
	return &org.ID, nil
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
