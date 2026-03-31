package platform

import (
	"context"
	"fmt"

	"github.com/Gleipnir-Technology/bob"
	"github.com/Gleipnir-Technology/bob/dialect/psql"
	"github.com/Gleipnir-Technology/bob/dialect/psql/sm"
	"github.com/Gleipnir-Technology/nidus-sync/db"
	"github.com/Gleipnir-Technology/nidus-sync/db/models"
	//"github.com/rs/zerolog/log"
	"github.com/stephenafamo/scan"
)

type Pool struct {
	Condition string `db:"condition" json:"condition"`
	ID        int32  `db:"id" json:"-"`
}
type UploadPoolError struct {
	Column  uint   `json:"column"`
	Line    uint   `json:"line"`
	Message string `json:"message"`
}

func errorsByLine(ctx context.Context, file *models.FileuploadFile) ([]UploadPoolError, map[int32][]UploadPoolError, error) {
	file_errors := make([]UploadPoolError, 0)
	errors_by_line := make(map[int32][]UploadPoolError, 0)
	error_rows, err := models.FileuploadErrorCSVS.Query(
		models.SelectWhere.FileuploadErrorCSVS.CSVFileID.EQ(file.ID),
	).All(ctx, db.PGInstance.BobDB)
	if err != nil {
		return file_errors, errors_by_line, fmt.Errorf("Failed to lookup errors in csv %d: %w", file.ID, err)
	}
	for _, row := range error_rows {
		e := UploadPoolError{
			Column:  uint(row.Col),
			Line:    uint(row.Line),
			Message: row.Message,
		}
		if row.Line == 0 {
			file_errors = append(file_errors, e)
		} else {
			//log.Info().Int32("line", row.Line).Msg("Found error")
			by_line, ok := errors_by_line[row.Line]
			if !ok {
				errors_by_line[row.Line] = []UploadPoolError{e}
				continue
			}
			by_line = append(by_line, e)
			errors_by_line[row.Line] = by_line
		}
	}
	return file_errors, errors_by_line, nil
}
func poolList(ctx context.Context, org_id int32, pool_ids []int32) ([]*Pool, error) {
	pools, err := bob.All(ctx, db.PGInstance.BobDB, psql.Select(
		sm.Columns(
			"condition",
			"feature_id AS id",
		),
		sm.From(psql.Quote("feature_pool")),
		sm.Where(
			models.FeaturePools.Columns.FeatureID.EQ(psql.Any(pool_ids)),
		),
	), scan.StructMapper[*Pool]())
	if err != nil {
		return nil, fmt.Errorf("query feature_pool: %w", err)
	}
	return pools, nil
}
