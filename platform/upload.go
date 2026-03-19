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
	"github.com/Gleipnir-Technology/nidus-sync/db/enums"
	"github.com/Gleipnir-Technology/nidus-sync/db/models"
	"github.com/Gleipnir-Technology/nidus-sync/platform/background"
	"github.com/Gleipnir-Technology/nidus-sync/platform/file"
	"github.com/aarondl/opt/omit"
	"github.com/aarondl/opt/omitnull"
	"github.com/rs/zerolog/log"
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

func NewUpload(ctx context.Context, u User, upload file.FileUpload, t enums.FileuploadCsvtype) (*Upload, error) {
	txn, err := db.PGInstance.BobDB.BeginTx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("Failed to begin transaction: %w", err)
	}
	defer txn.Rollback(ctx)

	file, err := models.FileuploadFiles.Insert(&models.FileuploadFileSetter{
		ContentType:    omit.From(upload.ContentType),
		Created:        omit.From(time.Now()),
		CreatorID:      omit.From(int32(u.ID)),
		Deleted:        omitnull.FromPtr[time.Time](nil),
		Name:           omit.From(upload.Name),
		OrganizationID: omit.From(u.Organization.ID()),
		Status:         omit.From(enums.FileuploadFilestatustypeUploaded),
		SizeBytes:      omit.From(int32(upload.SizeBytes)),
		FileUUID:       omit.From(upload.UUID),
	}).One(ctx, txn)
	if err != nil {
		return nil, fmt.Errorf("Failed to create file upload: %w", err)
	}
	_, err = models.FileuploadCSVS.Insert(&models.FileuploadCSVSetter{
		Committed: omitnull.FromPtr[time.Time](nil),
		FileID:    omit.From(file.ID),
		Rowcount:  omit.From(int32(0)),
		Type:      omit.From(t),
	}).One(ctx, txn)
	if err != nil {
		return nil, fmt.Errorf("Failed to create csv: %w", err)
	}
	log.Info().Int32("id", file.ID).Msg("Created new pool CSV upload")
	err = background.NewCSVImport(ctx, txn, file.ID)
	if err != nil {
		return nil, fmt.Errorf("background job create: %w", err)
	}
	txn.Commit(ctx)
	return &Upload{
		ID: file.ID,
	}, nil
}
func UploadCommit(ctx context.Context, org Organization, file_id int32, committer User) error {
	txn, err := db.PGInstance.BobDB.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("Failed to begin transaction: %w", err)
	}
	defer txn.Rollback(ctx)

	_, err = psql.Update(
		um.Table(models.FileuploadFiles.Alias()),
		um.SetCol("status").ToArg("committing"),
		um.SetCol("committer").ToArg(committer.ID),
		um.Where(psql.Quote("id").EQ(psql.Arg(file_id))),
		um.Where(psql.Quote("organization_id").EQ(psql.Arg(org.ID))),
	).Exec(ctx, txn)
	err = background.NewCSVCommit(ctx, txn, file_id)
	if err != nil {
		return fmt.Errorf("update upload: %w", err)
	}
	err = txn.Commit(ctx)

	return err
}
func UploadDiscard(ctx context.Context, org Organization, file_id int32) error {
	_, err := psql.Update(
		um.Table(models.FileuploadFiles.Alias()),
		um.SetCol("status").ToArg("discarded"),
		um.Where(psql.Quote("id").EQ(psql.Arg(file_id))),
		um.Where(psql.Quote("organization_id").EQ(psql.Arg(org.ID))),
	).Exec(ctx, db.PGInstance.BobDB)
	return err
}
func UploadSummaryList(ctx context.Context, org Organization) ([]UploadSummary, error) {
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
		sm.Where(psql.Quote("file", "organization_id").EQ(psql.Arg(org.ID()))),
		sm.OrderBy("created").Desc(),
	), scan.StructMapper[UploadSummary]())
	if err != nil {
		return results, fmt.Errorf("Failed to query pool upload rows: %w", err)
	}
	return rows, nil
}
