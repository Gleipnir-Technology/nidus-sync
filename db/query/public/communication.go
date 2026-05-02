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
func CommunicationsFromOrganization(ctx context.Context, org_id int64) ([]*model.Communication, error) {
	statement := table.Communication.SELECT(
		table.Communication.AllColumns,
	).FROM(table.Communication).
		WHERE(table.Communication.OrganizationID.EQ(postgres.Int(org_id)))
	return db.ExecuteMany[model.Communication](ctx, statement)
}
