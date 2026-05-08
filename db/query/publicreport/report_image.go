package publicreport

import (
	"context"

	"github.com/Gleipnir-Technology/nidus-sync/db"
	"github.com/Gleipnir-Technology/nidus-sync/db/gen/nidus-sync/publicreport/model"
	"github.com/Gleipnir-Technology/nidus-sync/db/gen/nidus-sync/publicreport/table"
	//"github.com/go-jet/jet/v2/postgres"
)

func ReportImageInsert(ctx context.Context, txn db.Ex, m model.ReportImage) (model.ReportImage, error) {
	statement := table.ReportImage.INSERT(table.ReportImage.MutableColumns).
		MODEL(m).
		RETURNING(table.ReportImage.AllColumns)
	return db.ExecuteOneTx[model.ReportImage](ctx, txn, statement)
}
func ReportImagesInsert(ctx context.Context, txn db.Ex, m []model.ReportImage) ([]model.ReportImage, error) {
	statement := table.ReportImage.INSERT(table.ReportImage.MutableColumns).
		MODELS(m).
		RETURNING(table.ReportImage.AllColumns)
	return db.ExecuteManyTx[model.ReportImage](ctx, txn, statement)
}
