package platform

import (
	"context"
	"errors"
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
	"github.com/Gleipnir-Technology/nidus-sync/platform/types"
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

type Upload struct {
	Created     time.Time      `db:"created" json:"created"`
	Error       string         `db:"error" json:"error"`
	Filename    string         `db:"filename" json:"filename"`
	ID          int32          `db:"id" json:"id"`
	RecordCount int            `db:"recordcount" json:"recordcount"`
	Status      string         `db:"status" json:"status"`
	Type        string         `db:"type" json:"type"`
	CSVPool     *CSVPoolDetail `json:"csv_pool"`
}

type CSVPoolDetailCount struct {
	Existing int `json:"existing"`
	New      int `json:"new"`
	Outside  int `json:"outside"`
}
type CSVPoolDetail struct {
	Count  CSVPoolDetailCount `json:"count"`
	Errors []UploadPoolError  `json:"errors"`
	Pools  []UploadPoolRow    `json:"pools"`
}
type UploadPoolRow struct {
	Address   types.Address     `json:"address"`
	Condition string            `json:"condition"`
	Errors    []UploadPoolError `json:"errors"`
	Status    string            `json:"status"`
	Tags      map[string]string `json:"tags"`
}

func GetUploadDetail(ctx context.Context, organization_id int32, file_id int32) (*Upload, error) {
	file, err := models.FindFileuploadFile(ctx, db.PGInstance.BobDB, file_id)
	if err != nil {
		return nil, fmt.Errorf("Failed to lookup file %d: %w", file_id, err)
	}
	csv, err := models.FindFileuploadCSV(ctx, db.PGInstance.BobDB, file_id)
	if err != nil {
		return nil, fmt.Errorf("Failed to lookup csv %d: %w", file_id, err)
	}
	switch csv.Type {
	case enums.FileuploadCsvtypeFlyover:
		return getUploadDetailPool(ctx, file)
	case enums.FileuploadCsvtypePoollist:
		return getUploadDetailPool(ctx, file)
	}
	return nil, errors.New("No idea what to do with upload type")
}

func NewUpload(ctx context.Context, u User, upload file.Upload, t enums.FileuploadCsvtype) (*int32, error) {
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
		Error:          omit.From(""),
		Name:           omit.From(upload.Name),
		OrganizationID: omit.From(u.Organization.ID),
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
	return &file.ID, nil
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
	if err != nil {
		return fmt.Errorf("update upload: %w", err)
	}
	err = background.NewCSVCommit(ctx, txn, file_id)
	if err != nil {
		return fmt.Errorf("background csv commit: %w", err)
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
func UploadList(ctx context.Context, org Organization) ([]Upload, error) {
	results := make([]Upload, 0)
	rows, err := bob.All(ctx, db.PGInstance.BobDB, psql.Select(
		sm.Columns(
			// fileupload.csv columns
			//"csv.file_id AS file_id",
			//"csv.committed",
			"csv.rowcount AS recordcount",
			"csv.type_ AS type",

			// fileupload.file columns
			//"file.content_type",
			"file.created AS created",
			//"file.creator_id",
			//"file.deleted",
			"file.error AS error",
			"file.id AS id",
			"file.name AS filename",
			//"file.organization_id",
			"file.status AS status",
			//"file.size_bytes",
			//"file.file_uuid",
			// Aggregate data
		),
		sm.From("fileupload.csv").As("csv"),
		sm.InnerJoin("fileupload.file").As("file").OnEQ(psql.Raw("csv.file_id"), psql.Raw("file.id")),
		sm.Where(psql.Quote("file", "organization_id").EQ(psql.Arg(org.ID))),
		sm.OrderBy("created").Desc(),
	), scan.StructMapper[Upload]())
	if err != nil {
		return results, fmt.Errorf("Failed to query pool upload rows: %w", err)
	}
	return rows, nil
}
func getUploadDetailPool(ctx context.Context, file *models.FileuploadFile) (*Upload, error) {
	file_errors, errors_by_line, err := errorsByLine(ctx, file)
	if err != nil {
		return nil, fmt.Errorf("get errors by line: %w", err)
	}
	pool_rows, err := models.FileuploadPools.Query(
		models.SelectWhere.FileuploadPools.CSVFile.EQ(file.ID),
	).All(ctx, db.PGInstance.BobDB)
	if err != nil {
		return nil, fmt.Errorf("Failed to query pools for %d: %w", file.ID, err)
	}
	address_ids := make([]int32, 0)
	for _, r := range pool_rows {
		if r.AddressID.IsValue() {
			address_ids = append(address_ids, r.AddressID.MustGet())
		}
	}
	addresses, err := AddressList(ctx, address_ids)
	if err != nil {
		return nil, fmt.Errorf("get address list: %w", err)
	}
	addresses_by_id := make(map[int32]*types.Address, len(address_ids))
	for _, a := range addresses {
		addresses_by_id[*a.ID] = a
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
		} else {
			count_existing = count_existing + 1
			status = "existing"
		}
		if !r.IsInDistrict {
			count_outside++
			status = "outside"
		}
		tags := db.ConvertFromPGData(r.Tags)
		// add 2 here because our file lines are 1-indexed and we skip the header line, but we are ranging 0-indexed
		errors, ok := errors_by_line[r.LineNumber]
		if !ok {
			errors = []UploadPoolError{}
		}
		var address *types.Address
		if r.AddressID.IsValue() {
			address = addresses_by_id[r.AddressID.MustGet()]
		} else {
			address = &types.Address{
				Country:    "usa",
				Locality:   r.AddressLocality,
				Number:     r.AddressNumber,
				PostalCode: r.AddressPostalCode,
				Region:     r.AddressRegion,
				Street:     r.AddressStreet,
			}
		}
		pools = append(pools, UploadPoolRow{
			Address:   *address,
			Condition: r.Condition.String(),
			Errors:    errors,
			Status:    status,
			Tags:      tags,
		})
	}
	return &Upload{
		Created:     file.Created,
		Error:       file.Error,
		Filename:    file.Name,
		ID:          file.ID,
		RecordCount: len(pool_rows),
		CSVPool: &CSVPoolDetail{
			Count: CSVPoolDetailCount{
				Existing: count_existing,
				Outside:  count_outside,
				New:      count_new,
			},
			Errors: file_errors,
			Pools:  pools,
		},
		Status: file.Status.String(),
	}, nil
}
