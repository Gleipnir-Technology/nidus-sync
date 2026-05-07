package publicreport

import (
	"context"

	//"github.com/Gleipnir-Technology/bob"
	"github.com/Gleipnir-Technology/nidus-sync/db"
	//"github.com/Gleipnir-Technology/nidus-sync/db/gen/nidus-sync/public/enum"
	"github.com/Gleipnir-Technology/nidus-sync/db/gen/nidus-sync/publicreport/model"
	"github.com/Gleipnir-Technology/nidus-sync/db/gen/nidus-sync/publicreport/table"
	//"github.com/go-jet/jet/v2/postgres"
)

func ImageExifInserts(ctx context.Context, txn db.Ex, image_exifs []model.ImageExif) ([]model.ImageExif, error) {
	statement := table.ImageExif.
		INSERT(table.ImageExif.MutableColumns).
		MODELS(image_exifs).
		RETURNING(table.ImageExif.AllColumns)
	return db.ExecuteManyTx[model.ImageExif](ctx, txn, statement)
}
