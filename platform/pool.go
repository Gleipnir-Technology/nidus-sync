package platform

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/Gleipnir-Technology/nidus-sync/db"
	"github.com/Gleipnir-Technology/nidus-sync/db/enums"
	"github.com/Gleipnir-Technology/nidus-sync/db/models"
	"github.com/Gleipnir-Technology/nidus-sync/platform/types"
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
	Address   types.Address
	Condition string
	Errors    []UploadPoolError
	Status    string
	Tags      map[string]string
}
type Upload struct {
	Created time.Time `db:"created"`
	ID      int32     `db:"id"`
	Status  string    `db:"status"`
}

func GetUploadDetail(ctx context.Context, organization_id int32, file_id int32) (UploadPoolDetail, error) {
	file, err := models.FindFileuploadFile(ctx, db.PGInstance.BobDB, file_id)
	if err != nil {
		return UploadPoolDetail{}, fmt.Errorf("Failed to lookup file %d: %w", file_id, err)
	}
	csv, err := models.FindFileuploadCSV(ctx, db.PGInstance.BobDB, file_id)
	if err != nil {
		return UploadPoolDetail{}, fmt.Errorf("Failed to lookup csv %d: %w", file_id, err)
	}
	switch csv.Type {
	case enums.FileuploadCsvtypeFlyover:
		return getUploadPoollistDetail(ctx, file)
	case enums.FileuploadCsvtypePoollist:
		return getUploadPoollistDetail(ctx, file)
	}
	return UploadPoolDetail{}, errors.New("No idea what to do with upload type")
}

func getUploadPoollistDetail(ctx context.Context, file *models.FileuploadFile) (UploadPoolDetail, error) {
	file_errors, errors_by_line, err := errorsByLine(ctx, file)
	if err != nil {
		return UploadPoolDetail{}, fmt.Errorf("get errors by line: %w", err)
	}
	pool_rows, err := models.FileuploadPools.Query(
		models.SelectWhere.FileuploadPools.CSVFile.EQ(file.ID),
	).All(ctx, db.PGInstance.BobDB)
	if err != nil {
		return UploadPoolDetail{}, fmt.Errorf("Failed to query pools for %d: %w", file.ID, err)
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
		pools = append(pools, UploadPoolRow{
			Address: types.Address{
				Country:    "usa",
				Locality:   r.AddressLocality,
				Number:     r.AddressNumber,
				PostalCode: r.AddressPostalCode,
				Region:     r.AddressRegion,
				Street:     r.AddressStreet,
			},
			Condition: r.Condition.String(),
			Errors:    errors,
			Status:    status,
			Tags:      tags,
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
