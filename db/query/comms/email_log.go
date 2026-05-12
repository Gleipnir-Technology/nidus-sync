package comms

import (
	"context"

	//"github.com/Gleipnir-Technology/bob"
	"github.com/Gleipnir-Technology/nidus-sync/db"
	//"github.com/Gleipnir-Technology/nidus-sync/db/gen/nidus-sync/public/enum"
	"github.com/Gleipnir-Technology/jet/postgres"
	"github.com/Gleipnir-Technology/nidus-sync/db/gen/nidus-sync/comms/model"
	"github.com/Gleipnir-Technology/nidus-sync/db/gen/nidus-sync/comms/table"
)

/*
func CommunicationInsert(ctx context.Context, txn db.Tx, m model.Communication) (model.Communication, error) {
	statement := table.Communication.INSERT(table.Communication.MutableColumns).
		MODEL(m).
		RETURNING(table.Communication.AllColumns)
	return db.ExecuteOneTx[model.Communication](ctx, txn, statement)
}
func CommunicationSetStatus(ctx context.Context, txn db.Tx, org_id int64, comm_id int64, status model.Communicationstatus) error {
	statement := table.Communication.UPDATE().
		SET(
			table.Communication.Status.SET(postgres.NewEnumValue(status.String())),
		).
		WHERE(table.Communication.OrganizationID.EQ(postgres.Int(org_id)).AND(
			table.Communication.ID.EQ(postgres.Int(comm_id))))
	return db.ExecuteNoneTx(ctx, txn, statement)
}
*/

func EmailLogFromID(ctx context.Context, id int64) (model.EmailLog, error) {
	statement := table.EmailLog.SELECT(
		table.EmailLog.AllColumns,
	).FROM(table.EmailLog).
		WHERE(table.EmailLog.ID.EQ(postgres.Int(id)))
	return db.ExecuteOne[model.EmailLog](ctx, statement)
}
func EmailLogsFromAddress(ctx context.Context, address string) ([]model.EmailLog, error) {
	statement := table.EmailLog.SELECT(
		table.EmailLog.AllColumns,
	).FROM(table.EmailLog).
		WHERE(table.EmailLog.Source.EQ(postgres.String(address)).OR(
			table.EmailLog.Destination.EQ(postgres.String(address))))
	return db.ExecuteMany[model.EmailLog](ctx, statement)
}
