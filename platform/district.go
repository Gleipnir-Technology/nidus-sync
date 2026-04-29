package platform

import (
	"context"
	"errors"
	"fmt"

	"github.com/Gleipnir-Technology/bob/dialect/psql"
	"github.com/Gleipnir-Technology/bob/dialect/psql/sm"
	"github.com/Gleipnir-Technology/nidus-sync/db"
	"github.com/Gleipnir-Technology/nidus-sync/db/models"
	"github.com/Gleipnir-Technology/nidus-sync/platform/types"

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
func matchDistrict(ctx context.Context, location *types.Location, images []ImageUpload, address *models.Address) (int32, error) {
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
			return org.ID, nil
		}
	}
	if location != nil && location.Longitude != 0 && location.Latitude != 0 {
		org, err = DistrictForLocation(ctx, location.Longitude, location.Latitude)
		if err != nil {
			return 0, fmt.Errorf("Failed to get district for location: %w", err)
		}
	}
	if address != nil {
		log.Debug().Msg("doing district match via address...")
		location, err = AddressLocation(ctx, address)
		if err != nil {
			return 0, fmt.Errorf("location for address: %w", err)
		}
		org, err = DistrictForLocation(ctx, location.Longitude, location.Latitude)
		if err != nil {
			return 0, fmt.Errorf("Failed to get district for location from address: %w", err)
		}
		log.Debug().Float64("loc.lat", location.Latitude).Float64("loc.lng", location.Longitude).Bool("org", org != nil).Msg("address match")
	}
	if org == nil {
		org, err = DistrictCatchall(ctx)
		if err != nil {
			return 0, fmt.Errorf("get catchall: %w", err)
		}
		log.Debug().Err(err).Int32("id", org.ID).Msg("No district match by report location, images, or address, using catchall")
		return org.ID, nil
	}
	log.Debug().Err(err).Int32("org_id", org.ID).Msg("Found district match for report")
	return org.ID, nil
}
