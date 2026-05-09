package platform

import (
	"context"
	"fmt"
	"time"

	"github.com/Gleipnir-Technology/nidus-sync/db"
	//"github.com/Gleipnir-Technology/nidus-sync/db/gen/nidus-sync/public/enum"
	"github.com/Gleipnir-Technology/nidus-sync/db/gen/nidus-sync/public/model"
	querypublic "github.com/Gleipnir-Technology/nidus-sync/db/query/public"
	"github.com/Gleipnir-Technology/nidus-sync/lint"
	//"github.com/Gleipnir-Technology/jet/postgres"
)

func CommunicationsForOrganization(ctx context.Context, org_id int64) ([]model.Communication, error) {
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
	return &comm, nil
}
func CommunicationMarkInvalid(ctx context.Context, user User, comm_id int32) error {
	return communicationMark(ctx, user, comm_id, model.Communicationstatus_Invalid, model.Communicationlogentry_StatusInvalidated)
}
func CommunicationMarkPendingResponse(ctx context.Context, user User, comm_id int32) error {
	return communicationMark(ctx, user, comm_id, model.Communicationstatus_Pending, model.Communicationlogentry_StatusPending)
}
func CommunicationMarkPossibleIssue(ctx context.Context, user User, comm_id int32) error {
	return communicationMark(ctx, user, comm_id, model.Communicationstatus_PossibleIssue, model.Communicationlogentry_StatusPossibleIssue)
}
func CommunicationMarkPossibleResolved(ctx context.Context, user User, comm_id int32) error {
	return communicationMark(ctx, user, comm_id, model.Communicationstatus_PossibleResolved, model.Communicationlogentry_StatusPossibleResolved)
}

func communicationMark(ctx context.Context, user User, comm_id int32, status model.Communicationstatus, log_type model.Communicationlogentry) error {
	txn, err := db.BeginTxn(ctx)
	if err != nil {
		return fmt.Errorf("begin txn: %w", err)
	}
	defer lint.LogOnErrRollback(txn.Rollback, ctx, "rollback")
	err = querypublic.CommunicationSetStatus(ctx, txn, int64(user.Organization.ID), int64(comm_id), status)
	if err != nil {
		return fmt.Errorf("mark: %w", err)
	}
	user_id := int32(user.ID)
	log_entry := model.CommunicationLogEntry{
		CommunicationID: comm_id,
		Created:         time.Now(),
		Type:            log_type,
		User:            &user_id,
	}
	querypublic.CommunicationLogEntryInsert(ctx, txn, log_entry)
	txn.Commit(ctx)
	return nil
}
