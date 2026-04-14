package platform

import (
	"context"
	"fmt"

	//"github.com/Gleipnir-Technology/bob/dialect/psql/sm"
	"github.com/Gleipnir-Technology/nidus-sync/db"
	"github.com/Gleipnir-Technology/nidus-sync/db/models"
	"github.com/Gleipnir-Technology/nidus-sync/platform/types"
)

func ServiceRequestList(ctx context.Context, user User, limit int) ([]*types.ServiceRequest, error) {
	syncs, err := models.FieldseekerServicerequests.Query(
		models.SelectWhere.FieldseekerServicerequests.OrganizationID.EQ(user.Organization.ID),
		//sm.OrderBy(models.FieldseekerServicerequests.Columns.Created).Desc(),
	).All(ctx, db.PGInstance.BobDB)
	if err != nil {
		return nil, fmt.Errorf("query sync: %w", err)
	}
	results := make([]*types.ServiceRequest, len(syncs))
	for i, s := range syncs {
		r := types.ServiceRequestFromModel(s)
		results[i] = &r
	}
	return results, nil
}
