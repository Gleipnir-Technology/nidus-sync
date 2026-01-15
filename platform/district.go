package platform

import (
	"context"
	"fmt"

	"github.com/Gleipnir-Technology/nidus-sync/db"
	"github.com/Gleipnir-Technology/nidus-sync/db/models"

	"github.com/stephenafamo/bob/dialect/psql"
	"github.com/stephenafamo/bob/dialect/psql/sm"
)

func DistrictForLocation(ctx context.Context, lng float64, lat float64) (*models.District, error) {
	rows, err := models.Districts.Query(
		sm.Where(
			psql.F("ST_Contains", psql.Raw("geom_4326"), psql.F("ST_SetSRID", psql.F("ST_MakePoint", psql.Arg(lng), psql.Arg(lat)), psql.Arg(4326))),
		),
	).All(ctx, db.PGInstance.BobDB)
	if err != nil {
		return nil, fmt.Errorf("failed to query district: %w", err)
	}
	switch len(rows) {
	case 0:
		return nil, nil
	case 1:
		return rows[0], nil
	default:
		return nil, nil
	}
}
