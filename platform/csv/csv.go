package csv

import (
	"context"
	//"encoding/csv"
	"fmt"
	//"io"
	"strconv"
	"strings"
	//"sync"
	"time"

	"github.com/Gleipnir-Technology/bob"
	"github.com/Gleipnir-Technology/bob/dialect/psql"
	"github.com/Gleipnir-Technology/bob/dialect/psql/um"
	//"github.com/Gleipnir-Technology/nidus-sync/config"
	"github.com/Gleipnir-Technology/nidus-sync/db"
	"github.com/Gleipnir-Technology/nidus-sync/db/enums"
	"github.com/Gleipnir-Technology/nidus-sync/db/models"
	"github.com/Gleipnir-Technology/nidus-sync/platform/event"
	"github.com/Gleipnir-Technology/nidus-sync/platform/geocode"
	"github.com/Gleipnir-Technology/nidus-sync/platform/types"
	//"github.com/Gleipnir-Technology/nidus-sync/h3utils"
	//"github.com/Gleipnir-Technology/nidus-sync/platform/geom"
	//"github.com/Gleipnir-Technology/nidus-sync/platform/text"
	//"github.com/Gleipnir-Technology/nidus-sync/stadia"
	//"github.com/Gleipnir-Technology/nidus-sync/userfile"
	"github.com/aarondl/opt/omit"
	"github.com/aarondl/opt/omitnull"
	"github.com/rs/zerolog/log"
)

type csvParserFunc[T any] = func(context.Context, bob.Tx, *models.FileuploadFile, *models.FileuploadCSV) ([]T, error)
type csvProcessorFunc[T any] = func(context.Context, bob.Tx, *models.FileuploadFile, *models.FileuploadCSV, []T) error

func JobCommit(ctx context.Context, txn bob.Executor, file_id int32) error {
	log.Debug().Int32("file_id", file_id).Msg("begin job commit")
	file, err := models.FindFileuploadFile(ctx, txn, file_id)
	if err != nil {
		return fmt.Errorf("Failed to get csv file %d from DB: %w", file_id, err)
	}
	org, err := models.FindOrganization(ctx, txn, file.OrganizationID)
	if err != nil {
		return fmt.Errorf("Failed to get org %d from DB: %w", file.OrganizationID, err)
	}

	rows, err := models.FileuploadPools.Query(
		models.SelectWhere.FileuploadPools.CSVFile.EQ(file_id),
	).All(ctx, txn)
	if err != nil {
		return fmt.Errorf("Failed to get all rows of file %d: %w", file_id, err)
	}
	address_ids := make([]int32, 0)
	for _, r := range rows {
		if r.AddressID.IsValue() {
			address_ids = append(address_ids, r.AddressID.MustGet())
		}
	}
	addresses, err := types.AddressList(ctx, address_ids)
	if err != nil {
		return fmt.Errorf("get address list: %w", err)
	}
	for _, row := range rows {
		var a *types.Address
		var parcel *models.Parcel
		if row.AddressID.IsValue() {
			var ok bool
			a, ok = addresses[row.AddressID.MustGet()]
			if !ok {
				log.Error().Int32("id", row.AddressID.MustGet()).Msg("address is missing")
				continue
			}
			parcel, err = geocode.GetParcel(ctx, txn, *a)
			if err != nil {
				return fmt.Errorf("get parcel: %w", err)
			}
		} else {
			a = &types.Address{
				Country:    "usa",
				Locality:   row.AddressLocality,
				Number:     row.AddressNumber,
				PostalCode: row.AddressPostalCode,
				Region:     row.AddressRegion,
				Street:     row.AddressStreet,
				Unit:       "",
			}
		}
		var site *models.Site
		site, err = models.Sites.Query(
			models.SelectWhere.Sites.AddressID.EQ(*a.ID),
		).One(ctx, txn)
		if err != nil {
			if err.Error() != "sql: no rows in result set" {
				return fmt.Errorf("query site: %w", err)
			}
			var parcel_id *int32
			if parcel != nil {
				parcel_id = &(*parcel).ID
			}
			setter := models.SiteSetter{
				AddressID: omit.From(*a.ID),
				Created:   omit.From(time.Now()),
				CreatorID: omit.FromPtr(file.Committer.Ptr()),
				FileID:    omitnull.From(file_id),
				//ID             omit.Val[int32]          `db:"id,pk" `
				Notes:          omit.From(row.Notes),
				OrganizationID: omit.From(org.ID),
				OwnerName:      omit.From(row.PropertyOwnerName),
				OwnerPhoneE164: omitnull.FromPtr(row.PropertyOwnerPhoneE164.Ptr()),
				ParcelID:       omitnull.FromPtr(parcel_id),
				ResidentOwned:  omitnull.FromPtr(row.ResidentOwned.Ptr()),
				Tags:           omit.From(row.Tags),
				Version:        omit.From(int32(1)),
			}
			site, err = models.Sites.Insert(&setter).One(ctx, txn)
			if err != nil {
				return fmt.Errorf("insert site: %w", err)
			}
		}
		var feature *models.Feature
		feature, err = models.Features.Query(
			models.SelectWhere.Features.OrganizationID.EQ(org.ID),
			models.SelectWhere.Features.SiteID.EQ(site.ID),
		).One(ctx, txn)
		if err != nil {
			if err.Error() != "sql: no rows in result set" {
				return fmt.Errorf("query site: %w", err)
			}
			feature, err = models.Features.Insert(&models.FeatureSetter{
				Created:   omit.From(time.Now()),
				CreatorID: omit.From(file.Committer.MustGet()),
				//ID: row.Address,
				OrganizationID: omit.From(org.ID),
				SiteID:         omit.From(site.ID),
			}).One(ctx, txn)
			if err != nil {
				return fmt.Errorf("insert feature: %w", err)
			}
			_, err := models.FeaturePools.Insert(&models.FeaturePoolSetter{
				Condition: omit.From(row.Condition),
				FeatureID: omit.From(feature.ID),
			}).One(ctx, txn)
			if err != nil {
				return fmt.Errorf("insert feature_pool: %w", err)
			}
		}
		review_task, err := models.ReviewTasks.Insert(&models.ReviewTaskSetter{
			Created:   omit.From(time.Now()),
			CreatorID: omitnull.From(file.Committer.MustGet()),
			//ID: row.Address,
			OrganizationID: omit.From(org.ID),
			Reviewed:       omitnull.FromPtr[time.Time](nil),
			ReviewerID:     omitnull.FromPtr[int32](nil),
		}).One(ctx, txn)
		if err != nil {
			return fmt.Errorf("insert review task: %w", err)
		}
		_, err = models.ReviewTaskPools.Insert(&models.ReviewTaskPoolSetter{
			FeaturePoolID: omit.From(feature.ID),
			Location:      omitnull.FromPtr[string](nil),
			Geometry:      omitnull.FromPtr[string](nil),
			ReviewTaskID:  omit.From(review_task.ID),
		}).One(ctx, txn)

		if err != nil {
			return fmt.Errorf("insert review task pool: %w", err)
		}
		/*
			Not sure why SignalPools doesn't have an Insert method
			_, err = models.SignalPools.Insert(&models.SignalPoolSetter{
				PoolID: omit.From(pool.ID),
				SignalID: omit.From(signal.ID),
			}).One(ctx, txn)
		*/
	}
	err = file.Update(ctx, txn, &models.FileuploadFileSetter{
		Status: omit.From(enums.FileuploadFilestatustypeCommitted),
	})
	if err != nil {
		return fmt.Errorf("update file status to committed: %w", err)
	}
	event.Updated(event.TypeFileCSV, file.OrganizationID, strconv.Itoa(int(file.ID)))
	return nil
}
func JobImport(ctx context.Context, txn bob.Executor, file_id int32) error {
	file, err := models.FileuploadFiles.Query(
		models.SelectWhere.FileuploadFiles.ID.EQ(file_id),
	).One(ctx, txn)
	if err != nil {
		return fmt.Errorf("find file: %w", err)
	}
	csv, err := models.FileuploadCSVS.Query(
		models.SelectWhere.FileuploadCSVS.FileID.EQ(file_id),
	).One(ctx, txn)
	if err != nil {
		return fmt.Errorf("find csv: %w", err)
	}

	switch csv.Type {
	case enums.FileuploadCsvtypePoollist:
		err = importCSV(ctx, file_id, parseCSVPoollist, processCSVPoollist)
	case enums.FileuploadCsvtypeFlyover:
		err = importCSV(ctx, file_id, parseCSVFlyover, processCSVFlyover)
	}
	if err != nil {
		log.Debug().Err(err).Msg("failed to import CSV")
		_, err := psql.Update(
			um.Table("fileupload.file"),
			um.SetCol("status").ToArg("error"),
			um.SetCol("error").ToArg(err.Error()),
			um.Where(psql.Quote("id").EQ(psql.Arg(file_id))),
		).Exec(ctx, db.PGInstance.BobDB)
		if err != nil {
			log.Error().Err(err).Msg("Failed to set upload to error status")
		}
	}
	event.Updated(event.TypeFileCSV, file.OrganizationID, strconv.Itoa(int(file.ID)))
	return nil
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
