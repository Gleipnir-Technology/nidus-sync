package platform

import (
	"context"
	"errors"
	"fmt"

	"github.com/Gleipnir-Technology/bob/dialect/psql"
	"github.com/Gleipnir-Technology/bob/dialect/psql/sm"
	"github.com/Gleipnir-Technology/nidus-sync/db"
	"github.com/Gleipnir-Technology/nidus-sync/db/models"

	"github.com/rs/zerolog/log"
)

func DistrictCatchall(ctx context.Context) (*models.Organization, error) {
	return models.Organizations.Query(
		models.SelectWhere.Organizations.IsCatchall.EQ(true),
	).One(ctx, db.PGInstance.BobDB)
}
func DistrictForLocation(ctx context.Context, lng float64, lat float64) (*models.Organization, error) {
	organizations, err := models.Organizations.Query(
		sm.Where(
			psql.F("ST_Contains", psql.Raw("service_area_geometry"), psql.F("ST_SetSRID", psql.F("ST_MakePoint", psql.Arg(lng), psql.Arg(lat)), psql.Arg(4326))),
		),
	).All(ctx, db.PGInstance.BobDB)

	log.Debug().Int("len", len(organizations)).Float64("lng", lng).Float64("lat", lat).Msg("Attempting district match")
	if err != nil {
		return nil, fmt.Errorf("failed to query organization: %w", err)
	}
	switch len(organizations) {
	case 0:
		return nil, nil
	case 1:
		org := organizations[0]
		return org, nil
	default:
		return nil, errors.New("too many organizations")
	}
}
func MatchDistrict(ctx context.Context, longitude, latitude float64, images []ImageUpload) (*int32, error) {
	var err error
	var org *models.Organization
	for _, image := range images {
		if image.Exif == nil {
			continue
		}
		if image.Exif.GPS == nil {
			continue
		}
		org, err = DistrictForLocation(ctx, image.Exif.GPS.Longitude, image.Exif.GPS.Latitude)
		if err != nil {
			log.Warn().Err(err).Msg("Failed to get district for location")
			continue
		}
		if org != nil {
			return &org.ID, nil
		}
	}
	if longitude == 0 || latitude == 0 {
		org, err = DistrictCatchall(ctx)
		if err != nil {
			return nil, fmt.Errorf("get catchall: %w", err)
		}
		log.Debug().Int32("id", org.ID).Msg("No location from images, no latlng for the report itself, using catchall")
		return &org.ID, nil
	}
	org, err = DistrictForLocation(ctx, longitude, latitude)
	if err != nil {
		log.Warn().Err(err).Msg("Failed to get district for location")
		return nil, fmt.Errorf("Failed to get district for location: %w", err)
	}
	if org == nil {
		org, err = DistrictCatchall(ctx)
		if err != nil {
			return nil, fmt.Errorf("get catchall: %w", err)
		}
		log.Debug().Err(err).Float64("lng", longitude).Float64("lat", latitude).Int32("id", org.ID).Msg("No district match by report location, using catchall")
		return &org.ID, nil
	}
	log.Debug().Err(err).Int32("org_id", org.ID).Float64("lng", longitude).Float64("lat", latitude).Msg("Found district match by report location")
	return &org.ID, nil
}
