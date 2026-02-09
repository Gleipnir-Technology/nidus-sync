package platform

import (
	"context"
	"fmt"
	"time"

	"github.com/Gleipnir-Technology/bob"
	"github.com/Gleipnir-Technology/bob/dialect/psql"
	"github.com/Gleipnir-Technology/bob/dialect/psql/sm"
	"github.com/Gleipnir-Technology/nidus-sync/background"
	"github.com/Gleipnir-Technology/nidus-sync/db"
	"github.com/Gleipnir-Technology/nidus-sync/db/enums"
	"github.com/Gleipnir-Technology/nidus-sync/db/models"
	"github.com/Gleipnir-Technology/nidus-sync/userfile"
	"github.com/aarondl/opt/omit"
	"github.com/aarondl/opt/omitnull"
	"github.com/rs/zerolog/log"
	"github.com/stephenafamo/scan"
)

type PoolUpload struct {
	Created time.Time `db:"created"`
	ID      int32     `db:"id"`
	Status  string    `db:"status"`
}

func NewPoolUpload(ctx context.Context, u *models.User, upload userfile.FileUpload) (PoolUpload, error) {
	txn, err := db.PGInstance.BobDB.BeginTx(ctx, nil)
	if err != nil {
		return PoolUpload{}, fmt.Errorf("Failed to begin transaction: %w", err)
	}
	defer txn.Rollback(ctx)

	file, err := models.FileuploadFiles.Insert(&models.FileuploadFileSetter{
		ContentType:    omit.From(upload.ContentType),
		Created:        omit.From(time.Now()),
		CreatorID:      omit.From(u.ID),
		Deleted:        omitnull.FromPtr[time.Time](nil),
		Name:           omit.From(upload.Name),
		OrganizationID: omit.From(u.OrganizationID),
		Status:         omit.From(enums.FileuploadFilestatustypeUploaded),
		SizeBytes:      omit.From(int32(upload.SizeBytes)),
		FileUUID:       omit.From(upload.UUID),
	}).One(ctx, txn)
	if err != nil {
		return PoolUpload{}, fmt.Errorf("Failed to create file upload: %w", err)
	}
	_, err = models.FileuploadCSVS.Insert(&models.FileuploadCSVSetter{
		Committed: omitnull.FromPtr[time.Time](nil),
		FileID:    omit.From(file.ID),
		Rowcount:  omit.From(int32(0)),
		Type:      omit.From(enums.FileuploadCsvtypePoollist),
	}).One(ctx, txn)
	if err != nil {
		return PoolUpload{}, fmt.Errorf("Failed to create csv: %w", err)
	}
	log.Info().Int32("id", file.ID).Msg("Created new pool CSV upload")
	txn.Commit(ctx)
	background.ProcessUpload(file.ID)
	return PoolUpload{
		ID: file.ID,
	}, nil
}
func PoolUploadList(ctx context.Context, organization_id int32) ([]PoolUpload, error) {
	results := make([]PoolUpload, 0)
	rows, err := bob.All(ctx, db.PGInstance.BobDB, psql.Select(
		sm.Columns(
			// fileupload.csv columns
			//"csv.file_id",
			//"csv.committed",
			//"csv.rowcount",
			//"csv.type_",

			// fileupload.file columns
			//"file.content_type",
			"file.created AS created",
			//"file.creator_id",
			//"file.deleted",
			"file.id AS id",
			//"file.name",
			//"file.organization_id",
			"file.status AS status",
			//"file.size_bytes",
			//"file.file_uuid",
		),
		sm.From("fileupload.csv").As("csv"),
		sm.InnerJoin("fileupload.file").As("file").OnEQ(psql.Raw("csv.file_id"), psql.Raw("file.id")),
		sm.Where(psql.Raw("file.organization_id").EQ(psql.Arg(organization_id))),
	), scan.StructMapper[PoolUpload]())
	if err != nil {
		return results, fmt.Errorf("Failed to query pool upload rows: %w", err)
	}
	return rows, nil
}
