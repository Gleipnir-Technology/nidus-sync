package platform

import (
	"context"

	"github.com/Gleipnir-Technology/nidus-sync/db/gen/nidus-sync/public/model"
	querypublic "github.com/Gleipnir-Technology/nidus-sync/db/query/public"
)

func CommunicationsForOrganization(ctx context.Context, org_id int64) ([]*model.Communication, error) {
	return querypublic.CommunicationsFromOrganization(ctx, org_id)
}
