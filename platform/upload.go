package platform

import (
	"context"
	"fmt"
	"time"

	"github.com/Gleipnir-Technology/bob"
	"github.com/Gleipnir-Technology/bob/dialect/psql"
	"github.com/Gleipnir-Technology/bob/dialect/psql/sm"
	"github.com/Gleipnir-Technology/bob/dialect/psql/um"
	"github.com/Gleipnir-Technology/nidus-sync/db"
	"github.com/Gleipnir-Technology/nidus-sync/db/models"
	"github.com/stephenafamo/scan"
)

type UploadType = int

const (
	UploadTypePool UploadType = iota
)

type UploadStatus = int

const (
	UploadStatusComplete UploadStatus = iota
)

type UploadSummary struct {
	Created     time.Time `db:"created"`
	Filename    string    `db:"filename"`
	ID          int32     `db:"id"`
	RecordCount int       `db:"recordcount"`
	Status      string    `db:"status"`
	Type        string    `db:"type"`
}

func UploadDiscard(ctx context.Context, org *models.Organization, file_id int32) error {
	_, err := psql.Update(
		um.Table(models.FileuploadFiles.Alias()),
		um.SetCol("status").ToArg("discarded"),
		um.Where(psql.Quote("id").EQ(psql.Arg(file_id))),
		um.Where(psql.Quote("organization_id").EQ(psql.Arg(org.ID))),
	).Exec(ctx, db.PGInstance.BobDB)
	return err
}
func UploadSummaryList(ctx context.Context, org *models.Organization) ([]UploadSummary, error) {
	results := make([]UploadSummary, 0)
	rows, err := bob.All(ctx, db.PGInstance.BobDB, psql.Select(
		sm.Columns(
			// fileupload.csv columns
			//"csv.file_id",
			//"csv.committed",
			"csv.rowcount AS recordcount",
			"csv.type_ AS type",

			// fileupload.file columns
			//"file.content_type",
			"file.created AS created",
			//"file.creator_id",
			//"file.deleted",
			"file.id AS id",
			"file.name AS filename",
			//"file.organization_id",
			"file.status AS status",
			//"file.size_bytes",
			//"file.file_uuid",

		),
		sm.From("fileupload.csv").As("csv"),
		sm.InnerJoin("fileupload.file").As("file").OnEQ(psql.Raw("csv.file_id"), psql.Raw("file.id")),
		sm.Where(psql.Raw("file.organization_id").EQ(psql.Arg(org.ID))),
		sm.OrderBy("created").Desc(),
	), scan.StructMapper[UploadSummary]())
	if err != nil {
		return results, fmt.Errorf("Failed to query pool upload rows: %w", err)
	}
	return rows, nil
}
