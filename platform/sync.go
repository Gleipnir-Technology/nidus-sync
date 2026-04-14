package platform

import (
	"context"
	"fmt"

	"github.com/Gleipnir-Technology/bob/dialect/psql/sm"
	"github.com/Gleipnir-Technology/nidus-sync/db"
	"github.com/Gleipnir-Technology/nidus-sync/db/models"
	"github.com/Gleipnir-Technology/nidus-sync/platform/types"
)

func SyncList(ctx context.Context, user User, limit int) ([]*types.Sync, error) {
	syncs, err := models.FieldseekerSyncs.Query(
		models.SelectWhere.FieldseekerSyncs.OrganizationID.EQ(user.Organization.ID),
		sm.OrderBy(models.FieldseekerSyncs.Columns.Created).Desc(),
	).All(ctx, db.PGInstance.BobDB)
	if err != nil {
		return nil, fmt.Errorf("query sync: %w", err)
	}
	results := make([]*types.Sync, len(syncs))
	for i, s := range syncs {
		r := types.SyncFromModel(s)
		results[i] = &r
	}
	return results, nil
}
