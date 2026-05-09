package public

import (
	"context"
	"time"

	"github.com/Gleipnir-Technology/nidus-sync/db"
	"github.com/Gleipnir-Technology/nidus-sync/db/gen/nidus-sync/public/model"
	"github.com/Gleipnir-Technology/nidus-sync/db/gen/nidus-sync/public/table"
	//"github.com/Gleipnir-Technology/jet/postgres"
)

func CommunicationLogEntryInsert(ctx context.Context, txn db.Tx, m model.CommunicationLogEntry) (model.CommunicationLogEntry, error) {
	m.Created = time.Now()
	statement := table.Communication.INSERT(table.Communication.MutableColumns).
		MODEL(m).
		RETURNING(table.Communication.AllColumns)
	return db.ExecuteOne[model.CommunicationLogEntry](ctx, statement)
}
