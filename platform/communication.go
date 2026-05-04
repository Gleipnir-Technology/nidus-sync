package platform

import (
	"context"
	"fmt"

	"github.com/Gleipnir-Technology/nidus-sync/db"
	"github.com/Gleipnir-Technology/nidus-sync/db/gen/nidus-sync/public/model"
	querypublic "github.com/Gleipnir-Technology/nidus-sync/db/query/public"
	"github.com/Gleipnir-Technology/nidus-sync/lint"
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
	return communicationMark(ctx, user, comm_id, querypublic.CommunicationMarkInvalid)
}
func CommunicationMarkPendingResponse(ctx context.Context, user User, comm_id int64) error {
	return communicationMark(ctx, user, comm_id, querypublic.CommunicationMarkPendingResponse)
}
func CommunicationMarkPossibleIssue(ctx context.Context, user User, comm_id int64) error {
	return communicationMark(ctx, user, comm_id, querypublic.CommunicationMarkPossibleIssue)
}
func CommunicationMarkPossibleResolved(ctx context.Context, user User, comm_id int64) error {
	return communicationMark(ctx, user, comm_id, querypublic.CommunicationMarkPossibleResolved)
}

type markFunc = func(context.Context, db.Tx, int64, int64, int64) error

func communicationMark(ctx context.Context, user User, comm_id int64, f markFunc) error {
	txn, err := db.BeginTxn(ctx)
	if err != nil {
		return fmt.Errorf("begin txn: %w", err)
	}
	defer lint.LogOnErrRollback(txn.Rollback, ctx, "rollback")
	err = f(ctx, txn, int64(user.Organization.ID), int64(user.ID), comm_id)
	if err != nil {
		return fmt.Errorf("mark: %w", err)
	}
	txn.Commit(ctx)
	return nil
}
