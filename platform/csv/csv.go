package csv

import (
	"context"
	//"encoding/csv"
	"fmt"
	//"io"
	"strconv"
	"strings"
	//"sync"
	//"time"

	"github.com/Gleipnir-Technology/bob"
	"github.com/Gleipnir-Technology/bob/dialect/psql"
	"github.com/Gleipnir-Technology/bob/dialect/psql/um"
	//"github.com/Gleipnir-Technology/nidus-sync/config"
	"github.com/Gleipnir-Technology/nidus-sync/db"
	"github.com/Gleipnir-Technology/nidus-sync/db/enums"
	"github.com/Gleipnir-Technology/nidus-sync/db/models"
	"github.com/Gleipnir-Technology/nidus-sync/platform/geocode"
	//"github.com/Gleipnir-Technology/nidus-sync/h3utils"
	//"github.com/Gleipnir-Technology/nidus-sync/platform/geom"
	//"github.com/Gleipnir-Technology/nidus-sync/platform/text"
	//"github.com/Gleipnir-Technology/nidus-sync/stadia"
	//"github.com/Gleipnir-Technology/nidus-sync/userfile"
	"github.com/aarondl/opt/omit"
	//"github.com/aarondl/opt/omitnull"
	"github.com/rs/zerolog/log"
)

type csvParserFunc[T any] = func(context.Context, bob.Tx, *models.FileuploadFile, *models.FileuploadCSV) ([]T, error)
type csvProcessorFunc[T any] = func(context.Context, bob.Tx, *models.FileuploadFile, *models.FileuploadCSV, []T) error

func JobCommit(ctx context.Context, file_id int32) error {
	txn, err := db.PGInstance.BobDB.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("Failed to start transaction: %w", err)
	}

	f, err := models.FindFileuploadFile(ctx, txn, file_id)
	if err != nil {
		return fmt.Errorf("Failed to get csv file %d from DB: %w", file_id, err)
	}
	org, err := models.FindOrganization(ctx, txn, f.OrganizationID)
	if err != nil {
		return fmt.Errorf("Failed to get org %d from DB: %w", f.OrganizationID, err)
	}

	rows, err := models.FileuploadPools.Query(
		models.SelectWhere.FileuploadPools.CSVFile.EQ(file_id),
	).All(ctx, txn)
	if err != nil {
		return fmt.Errorf("Failed to get all rows of file %d: %w", file_id, err)
	}
	for _, row := range rows {
		a := geocode.Address{
			Country:    enums.CountrytypeUsa,
			Locality:   row.AddressLocality,
			Number:     row.AddressNumber,
			PostalCode: row.AddressPostalCode,
			Region:     row.AddressRegion,
			Street:     row.AddressStreet,
			Unit:       "",
		}
		address, err := geocode.EnsureAddress(ctx, txn, org, a)
		if err != nil {
			return fmt.Errorf("ensure address: %w", err)
		}
		log.Info().Int32("id", address.ID).Msg("made address")
	}
	return nil
}
func JobImport(ctx context.Context, file_id int32, type_ enums.FileuploadCsvtype) error {
	var err error
	switch type_ {
	case enums.FileuploadCsvtypePoollist:
		err = importCSV(ctx, file_id, parseCSVPoollist, processCSVPoollist)
	case enums.FileuploadCsvtypeFlyover:
		err = importCSV(ctx, file_id, parseCSVFlyover, processCSVFlyover)
	}
	if err != nil {
		_, err := psql.Update(
			um.Table("fileupload.file"),
			um.SetCol("status").ToArg("error"),
			um.Where(psql.Quote("id").EQ(psql.Arg(file_id))),
		).Exec(ctx, db.PGInstance.BobDB)
		if err != nil {
			log.Error().Err(err).Msg("Failed to set upload to error status")
		}
	}
	return err
}

func importCSV[T any](ctx context.Context, file_id int32, parser csvParserFunc[T], processor csvProcessorFunc[T]) error {
	// Not done in the transaction so the state shows up immediately
	_, err := psql.Update(
		um.Table("fileupload.file"),
		um.SetCol("status").ToArg("parsing"),
		um.Where(psql.Quote("id").EQ(psql.Arg(file_id))),
	).Exec(ctx, db.PGInstance.BobDB)
	if err != nil {
		return fmt.Errorf("Failed to set file %d to processing: %w", file_id, err)
	}

	file, c, err := loadFileAndCSV(ctx, file_id)
	if err != nil {
		return fmt.Errorf("load file and csv: %w", err)
	}
	txn, err := db.PGInstance.BobDB.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("Failed to start transaction: %w", err)
	}
	defer txn.Rollback(ctx)
	parsed, err := parser(ctx, txn, file, c)
	if err != nil {
		return fmt.Errorf("parse file: %w", err)
	}
	_, err = psql.Update(
		um.Table("fileupload.csv"),
		um.SetCol("rowcount").ToArg(len(parsed)),
		um.Where(psql.Quote("file_id").EQ(psql.Arg(file_id))),
	).Exec(ctx, txn)
	if err != nil {
		return fmt.Errorf("update csv row: %w", err)
	}
	err = processor(ctx, txn, file, c, parsed)
	if err != nil {
		return fmt.Errorf("process parsed file: %w", err)
	}

	file.Update(ctx, txn, &models.FileuploadFileSetter{
		Status: omit.From(enums.FileuploadFilestatustypeParsed),
	})
	log.Info().Int32("file.ID", file.ID).Msg("Set file to parsed")
	txn.Commit(ctx)
	return nil
}
func loadFileAndCSV(ctx context.Context, file_id int32) (*models.FileuploadFile, *models.FileuploadCSV, error) {
	file, err := models.FindFileuploadFile(ctx, db.PGInstance.BobDB, file_id)
	if err != nil {
		return nil, nil, fmt.Errorf("Failed to get file %d from DB: %w", file_id, err)
	}
	c, err := models.FindFileuploadCSV(ctx, db.PGInstance.BobDB, file.ID)
	if err != nil {
		return nil, nil, fmt.Errorf("Failed to get csv file %d from DB: %w", file.ID, err)
	}
	return file, c, nil
}

func addError(ctx context.Context, txn bob.Tx, c *models.FileuploadCSV, row_number int32, column_number int32, msg string) error {
	r, err := models.FileuploadErrorCSVS.Insert(&models.FileuploadErrorCSVSetter{
		Col:       omit.From(column_number),
		CSVFileID: omit.From(c.FileID),
		// ID
		Line:    omit.From(row_number),
		Message: omit.From(msg),
	}).One(ctx, txn)
	if err != nil {
		return fmt.Errorf("Failed to add error: %w", err)
	}
	log.Info().Int32("id", r.ID).Int32("file_id", c.FileID).Str("msg", msg).Int32("row", row_number).Int32("col", column_number).Msg("Created CSV file error")
	return nil
}
func addImportError(file *models.FileuploadFile, err error) {
	log.Debug().Err(err).Int32("file_id", file.ID).Msg("Fake add import error")
}
func parseBool(s string) (bool, error) {
	sl := strings.ToLower(s)
	boolValue, err := strconv.ParseBool(sl)
	if err != nil {
		// Handle some of the stuff that strconv doesn't handle
		switch sl {
		case "yes":
			return true, nil
		case "no":
			return false, nil
		default:
			return false, fmt.Errorf("unrecognized '%s'", sl)
		}

	}
	return boolValue, err
}

func errorMissingHeader(ctx context.Context, txn bob.Tx, c *models.FileuploadCSV, h headerPoolEnum) error {
	msg := fmt.Sprintf("The file is missing the '%s' header", h.String())
	return addError(ctx, txn, c, 0, 0, msg)
}
