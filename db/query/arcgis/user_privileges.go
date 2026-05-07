package arcgis

import (
	"context"

	"github.com/Gleipnir-Technology/bob"
	"github.com/Gleipnir-Technology/nidus-sync/db"
	"github.com/Gleipnir-Technology/nidus-sync/db/gen/nidus-sync/arcgis/model"
	"github.com/Gleipnir-Technology/nidus-sync/db/gen/nidus-sync/arcgis/table"
	"github.com/go-jet/jet/v2/postgres"
)

func UserPrivilegesDeleteByUserID(ctx context.Context, txn bob.Tx, id string) error {
	statement := table.User.DELETE().
		WHERE(table.User.ID.EQ(postgres.String(id)))
	return db.ExecuteNoneTxBob(ctx, txn, statement)
}
func UserPrivilegeInsert(ctx context.Context, txn bob.Tx, m *model.UserPrivilege) error {
	statement := table.UserPrivilege.INSERT(table.UserPrivilege.MutableColumns).
		MODEL(m).
		RETURNING(table.UserPrivilege.AllColumns)
	return db.ExecuteNoneTxBob(ctx, txn, statement)
}
