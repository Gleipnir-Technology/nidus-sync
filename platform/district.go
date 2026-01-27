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

func DistrictForLocation(ctx context.Context, lng float64, lat float64) (*models.ImportDistrict, *models.Organization, error) {
	districts, err := models.ImportDistricts.Query(
		sm.Where(
			psql.F("ST_Contains", psql.Raw("geom_4326"), psql.F("ST_SetSRID", psql.F("ST_MakePoint", psql.Arg(lng), psql.Arg(lat)), psql.Arg(4326))),
		),
	).All(ctx, db.PGInstance.BobDB)

	log.Debug().Int("len", len(districts)).Float64("lng", lng).Float64("lat", lat).Msg("Attempting district match")
	if err != nil {
		return nil, nil, fmt.Errorf("failed to query district: %w", err)
	}
	switch len(districts) {
	case 0:
		return nil, nil, nil
	case 1:
		district := districts[0]
		organizations, err := models.Organizations.Query(
			sm.Where(
				models.Organizations.Columns.ImportDistrictGid.EQ(psql.Arg(district.Gid)),
			),
		).All(ctx, db.PGInstance.BobDB)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to query organization: %w", err)
		}
		switch len(organizations) {
		case 0:
			return nil, nil, nil
		case 1:
			return district, organizations[0], nil
		default:
			return nil, nil, errors.New("too many organizations")
		}
	default:
		return nil, nil, errors.New("too many districts")
	}
}
