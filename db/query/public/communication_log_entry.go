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
	statement := table.CommunicationLogEntry.INSERT(table.CommunicationLogEntry.MutableColumns).
		MODEL(m).
		RETURNING(table.CommunicationLogEntry.AllColumns)
	return db.ExecuteOne[model.CommunicationLogEntry](ctx, statement)
}
