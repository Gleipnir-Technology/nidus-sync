package public

import (
	"context"
	"time"

	"github.com/Gleipnir-Technology/bob"
	"github.com/Gleipnir-Technology/nidus-sync/db"
	"github.com/Gleipnir-Technology/nidus-sync/db/gen/nidus-sync/public/model"
	"github.com/Gleipnir-Technology/nidus-sync/db/gen/nidus-sync/public/table"
	"github.com/go-jet/jet/v2/postgres"
)

func CommunicationInsert(ctx context.Context, txn bob.Tx, m *model.Communication) (*model.Communication, error) {
	m.Created = time.Now()
	statement := table.Communication.INSERT(table.Communication.MutableColumns).
		MODEL(m).
		RETURNING(table.Communication.AllColumns)
	return db.ExecuteOne[model.Communication](ctx, statement)
}
func CommunicationFromID(ctx context.Context, comm_id int64) (*model.Communication, error) {
	statement := table.Communication.SELECT(
		table.Communication.AllColumns,
	).FROM(table.Communication).
		WHERE(table.Communication.ID.EQ(postgres.Int(comm_id)))
	return db.ExecuteOne[model.Communication](ctx, statement)
}
func CommunicationsFromOrganization(ctx context.Context, org_id int64) ([]*model.Communication, error) {
	statement := table.Communication.SELECT(
		table.Communication.AllColumns,
	).FROM(table.Communication).
		WHERE(table.Communication.OrganizationID.EQ(postgres.Int(org_id)))
	return db.ExecuteMany[model.Communication](ctx, statement)
}
func CommunicationMarkInvalid(ctx context.Context, org_id int64, user_id int64, comm_id int64) error {
	statement := table.Communication.UPDATE().
		SET(
			table.Communication.Invalidated.SET(postgres.TimestampT(time.Now())),
			table.Communication.InvalidatedBy.SET(postgres.Int(user_id)),
		).
		WHERE(table.Communication.OrganizationID.EQ(postgres.Int(org_id)).AND(
			table.Communication.ID.EQ(postgres.Int(comm_id))))
	return db.ExecuteNone(ctx, statement)
}
func CommunicationMarkPendingResponse(ctx context.Context, org_id int64, user_id int64, comm_id int64) error {
	statement := table.Communication.UPDATE().
		SET(
			table.Communication.SetPending.SET(postgres.TimestampT(time.Now())),
			table.Communication.SetPendingBy.SET(postgres.Int(user_id)),
		).
		WHERE(table.Communication.OrganizationID.EQ(postgres.Int(org_id)).AND(
			table.Communication.ID.EQ(postgres.Int(comm_id))))
	return db.ExecuteNone(ctx, statement)
}
func CommunicationMarkPossibleIssue(ctx context.Context, org_id int64, user_id int64, comm_id int64) error {
	statement := table.Communication.UPDATE().
		SET(
			table.Communication.SetPossibleIssue.SET(postgres.TimestampT(time.Now())),
			table.Communication.SetPossibleIssueBy.SET(postgres.Int(user_id)),
		).
		WHERE(table.Communication.OrganizationID.EQ(postgres.Int(org_id)).AND(
			table.Communication.ID.EQ(postgres.Int(comm_id))))
	return db.ExecuteNone(ctx, statement)
}
func CommunicationMarkPossibleResolved(ctx context.Context, org_id int64, user_id int64, comm_id int64) error {
	statement := table.Communication.UPDATE().
		SET(
			table.Communication.SetPossibleResolved.SET(postgres.TimestampT(time.Now())),
			table.Communication.SetPossibleResolvedBy.SET(postgres.Int(user_id)),
		).
		WHERE(table.Communication.OrganizationID.EQ(postgres.Int(org_id)).AND(
			table.Communication.ID.EQ(postgres.Int(comm_id))))
	return db.ExecuteNone(ctx, statement)
}
