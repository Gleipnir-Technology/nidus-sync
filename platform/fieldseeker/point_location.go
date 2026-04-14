package fieldseeker

import (
	"context"
	"fmt"

	//"github.com/Gleipnir-Technology/bob"
	"github.com/Gleipnir-Technology/bob/dialect/psql"
	"github.com/Gleipnir-Technology/bob/dialect/psql/sm"
	"github.com/Gleipnir-Technology/nidus-sync/db"
	"github.com/Gleipnir-Technology/nidus-sync/db/models"
	//"github.com/Gleipnir-Technology/nidus-sync/db/sql"
	"github.com/google/uuid"
)

func PointLocationList(ctx context.Context, point_location_ids []uuid.UUID) (models.FieldseekerPointlocationSlice, error) {
	rows, err := models.FieldseekerPointlocations.Query(
		sm.Where(
			models.FieldseekerPointlocations.Columns.Globalid.EQ(psql.Any(point_location_ids)),
		),
	).All(ctx, db.PGInstance.BobDB)
	if err != nil {
		return nil, fmt.Errorf("query point locations: %w", err)
	}
	return rows, nil
}
