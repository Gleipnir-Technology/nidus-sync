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

type UploadPoolDetail struct {
	CountExisting int
	CountNew      int
	CountOutside  int
	Created       time.Time
	Errors        []UploadPoolError
	ID            int32
	Name          string
	Pools         []UploadPoolRow
	Status        string
}
type UploadPoolError struct {
	Column  uint
	Line    uint
	Message string
}
type UploadPoolRow struct {
	City       string
	Condition  string
	Errors     []UploadPoolError
	PostalCode string
	Status     string
	Street     string
	Tags       map[string]string
}
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
func GetUploadPoolDetail(ctx context.Context, organization_id int32, file_id int32) (UploadPoolDetail, error) {
	file, err := models.FindFileuploadFile(ctx, db.PGInstance.BobDB, file_id)
	if err != nil {
		return UploadPoolDetail{}, fmt.Errorf("Failed to lookup file %d: %w", file_id, err)
	}
	/*
		csv, err := models.FindFileuploadCSV(ctx, db.PGInstance.BobDB, file_id)
		if err != nil {
			return UploadPoolDetail{}, fmt.Errorf("Failed to lookup csv %d: %w", file_id, err)
		}
	*/
	error_rows, err := models.FileuploadErrorCSVS.Query(
		models.SelectWhere.FileuploadErrorCSVS.CSVFileID.EQ(file_id),
	).All(ctx, db.PGInstance.BobDB)
	if err != nil {
		return UploadPoolDetail{}, fmt.Errorf("Failed to lookup errors in csv %d: %w", file_id, err)
	}
	file_errors := make([]UploadPoolError, 0)
	errors_by_line := make(map[int32][]UploadPoolError, 0)
	for _, row := range error_rows {
		e := UploadPoolError{
			Column:  uint(row.Col),
			Line:    uint(row.Line),
			Message: row.Message,
		}
		if row.Line == 0 {
			file_errors = append(file_errors, e)
		} else {
			log.Info().Int32("line", row.Line).Msg("Found error")
			by_line, ok := errors_by_line[row.Line]
			if !ok {
				errors_by_line[row.Line] = []UploadPoolError{e}
				continue
			}
			by_line = append(by_line, e)
			errors_by_line[row.Line] = by_line
		}
	}

	pool_rows, err := models.FileuploadPools.Query(
		models.SelectWhere.FileuploadPools.CSVFile.EQ(file_id),
	).All(ctx, db.PGInstance.BobDB)
	if err != nil {
		return UploadPoolDetail{}, fmt.Errorf("Failed to query pools for %d: %w", file_id, err)
	}

	pools := make([]UploadPoolRow, 0)
	count_existing := 0
	count_new := 0
	count_outside := 0
	status := "unknown"
	for _, r := range pool_rows {
		if r.IsNew {
			count_new = count_new + 1
			status = "new"
		} else if !r.IsInDistrict {
			count_outside = count_outside + 1
			status = "outside"
		} else {
			count_existing = count_existing + 1
			status = "existing"
		}
		tags := db.ConvertFromPGData(r.Tags)
		// add 2 here because our file lines are 1-indexed and we skip the header line, but we are ranging 0-indexed
		errors, ok := errors_by_line[r.LineNumber]
		if !ok {
			errors = []UploadPoolError{}
		} else {
			log.Info().Int32("line", r.LineNumber).Int32("id", r.ID).Msg("Found errors in errors_by_line")
		}
		pools = append(pools, UploadPoolRow{
			City:       r.AddressCity,
			Condition:  r.Condition.String(),
			Errors:     errors,
			PostalCode: r.AddressPostalCode,
			Status:     status,
			Street:     r.AddressStreet,
			Tags:       tags,
		})
	}
	return UploadPoolDetail{
		CountExisting: count_existing,
		CountOutside:  count_outside,
		CountNew:      count_new,
		Errors:        file_errors,
		Name:          file.Name,
		Pools:         pools,
		Status:        file.Status.String(),
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
