package platform

import (
	"context"

	"github.com/Gleipnir-Technology/nidus-sync/db/gen/nidus-sync/public/model"
	querypublic "github.com/Gleipnir-Technology/nidus-sync/db/query/public"
)

func CommunicationsForOrganization(ctx context.Context, org_id int64) ([]*model.Communication, error) {
	return querypublic.CommunicationsFromOrganization(ctx, org_id)
}
func CommunicationFromID(ctx context.Context, user User, comm_id int64) (*model.Communication, error) {
	comm, err := querypublic.CommunicationFromID(ctx, comm_id)
	if err != nil {
		return nil, err
	}
	if comm.OrganizationID != user.Organization.ID {
		return nil, nil
	}
	return comm, nil
}
func CommunicationMarkInvalid(ctx context.Context, user User, comm_id int64) error {
	return querypublic.CommunicationMarkInvalid(ctx, int64(user.Organization.ID), int64(user.ID), comm_id)
}
func CommunicationMarkPendingResponse(ctx context.Context, user User, comm_id int64) error {
	return querypublic.CommunicationMarkPendingResponse(ctx, int64(user.Organization.ID), int64(user.ID), comm_id)
}
func CommunicationMarkPossibleIssue(ctx context.Context, user User, comm_id int64) error {
	return querypublic.CommunicationMarkPossibleIssue(ctx, int64(user.Organization.ID), int64(user.ID), comm_id)
}
func CommunicationMarkPossibleResolved(ctx context.Context, user User, comm_id int64) error {
	return querypublic.CommunicationMarkPossibleResolved(ctx, int64(user.Organization.ID), int64(user.ID), comm_id)
}
