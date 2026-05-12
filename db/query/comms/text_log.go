package comms

import (
	"context"

	"github.com/Gleipnir-Technology/jet/postgres"
	"github.com/Gleipnir-Technology/nidus-sync/db"
	"github.com/Gleipnir-Technology/nidus-sync/db/gen/nidus-sync/comms/model"
	"github.com/Gleipnir-Technology/nidus-sync/db/gen/nidus-sync/comms/table"
)

func TextLogFromID(ctx context.Context, id int64) (model.TextLog, error) {
	statement := table.TextLog.SELECT(
		table.TextLog.AllColumns,
	).FROM(table.TextLog).
		WHERE(table.TextLog.ID.EQ(postgres.Int(id)))
	return db.ExecuteOne[model.TextLog](ctx, statement)
}
