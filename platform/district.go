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
