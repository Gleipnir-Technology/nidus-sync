package public

import (
	"context"

	//"github.com/Gleipnir-Technology/bob"
	"github.com/Gleipnir-Technology/nidus-sync/db"
	//"github.com/Gleipnir-Technology/nidus-sync/db/gen/nidus-sync/public/enum"
	"github.com/Gleipnir-Technology/jet/postgres"
	"github.com/Gleipnir-Technology/nidus-sync/db/gen/nidus-sync/public/model"
	"github.com/Gleipnir-Technology/nidus-sync/db/gen/nidus-sync/public/table"
)

func CommunicationInsert(ctx context.Context, txn db.Tx, m model.Communication) (model.Communication, error) {
	statement := table.Communication.INSERT(table.Communication.MutableColumns).
		MODEL(m).
		RETURNING(table.Communication.AllColumns)
	return db.ExecuteOneTx[model.Communication](ctx, txn, statement)
}
func CommunicationFromID(ctx context.Context, comm_id int64) (model.Communication, error) {
	statement := table.Communication.SELECT(
		table.Communication.AllColumns,
	).FROM(table.Communication).
		WHERE(table.Communication.ID.EQ(postgres.Int(comm_id)))
	return db.ExecuteOne[model.Communication](ctx, statement)
}
func CommunicationsFromOrganization(ctx context.Context, org_id int64) ([]model.Communication, error) {
	statement := table.Communication.SELECT(
		table.Communication.AllColumns,
	).FROM(table.Communication).
		WHERE(table.Communication.OrganizationID.EQ(postgres.Int(org_id))).
		ORDER_BY(table.Communication.Created.DESC())
	return db.ExecuteMany[model.Communication](ctx, statement)
}
func CommunicationSetStatus(ctx context.Context, txn db.Tx, org_id int64, comm_id int64, status model.Communicationstatus) error {
	statement := table.Communication.UPDATE().
		SET(
			table.Communication.Status.SET(postgres.String(status.String())),
		).
		WHERE(table.Communication.OrganizationID.EQ(postgres.Int(org_id)).AND(
			table.Communication.ID.EQ(postgres.Int(comm_id))))
	return db.ExecuteNoneTx(ctx, txn, statement)
}
