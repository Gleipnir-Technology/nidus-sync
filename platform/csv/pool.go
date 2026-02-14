package csv

import (
	"context"
	"encoding/csv"
	"fmt"
	"io"
	"strconv"
	"strings"
	"time"

	"github.com/Gleipnir-Technology/bob"
	"github.com/Gleipnir-Technology/nidus-sync/config"
	"github.com/Gleipnir-Technology/nidus-sync/db"
	"github.com/Gleipnir-Technology/nidus-sync/db/enums"
	"github.com/Gleipnir-Technology/nidus-sync/db/models"
	"github.com/Gleipnir-Technology/nidus-sync/platform/text"
	"github.com/Gleipnir-Technology/nidus-sync/stadia"
	"github.com/Gleipnir-Technology/nidus-sync/userfile"
	"github.com/aarondl/opt/omit"
	"github.com/aarondl/opt/omitnull"
	"github.com/rs/zerolog/log"
)

type headerPoolEnum int

const (
	headerAddressCity = iota
	headerAddressPostalCode
	headerAddressStreet
	headerCondition
	headerNotes
	headerPropertyOwnerName
	headerPropertyOwnerPhone
	headerResidentOwned
	headerResidentPhone
	headerTag
)

func (e headerPoolEnum) String() string {
	switch e {
	case headerAddressCity:
		return "City"
	case headerAddressPostalCode:
		return "Postal Code"
	case headerAddressStreet:
		return "Street Address"
	case headerCondition:
		return "Condition"
	case headerNotes:
		return "Notes"
	case headerPropertyOwnerName:
		return "Property Owner Name"
	case headerPropertyOwnerPhone:
		return "Property Owner Phone"
	case headerResidentOwned:
		return "Resident Owned"
	case headerResidentPhone:
		return "Resident Phone"
	default:
		return "bad programmer"
	}
}
func ProcessJob(ctx context.Context, file_id int32) error {
	file, err := models.FindFileuploadFile(ctx, db.PGInstance.BobDB, file_id)
	if err != nil {
		return fmt.Errorf("Failed to get file %d from DB: %w", file_id, err)
	}
	txn, err := db.PGInstance.BobDB.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("Failed to start transaction: %w", err)
	}
	defer txn.Rollback(ctx)
	pools, err := parseFile(ctx, txn, *file)
	if err != nil {
		return fmt.Errorf("parse file: %w", err)
	}
	err = bulkGeocode(ctx, txn, *file, pools)
	if err != nil {
		return fmt.Errorf("bulk geocode: %w", err)
	}
	return nil
}
func bulkGeocode(ctx context.Context, txn bob.Tx, file models.FileuploadFile, pools []*models.FileuploadPool) error {
	if len(pools) == 0 {
		return nil
	}
	client := stadia.NewStadiaMaps(config.StadiaMapsAPIKey)
	requests := make([]stadia.BulkGeocodeQuery, 0)
	for _, pool := range pools {
		requests = append(requests, stadia.StructuredGeocodeRequest{
			Address:    &pool.AddressStreet,
			PostalCode: &pool.AddressPostalCode,
		})
	}
	log.Info().Int("len pools", len(pools)).Int("len requests", len(requests)).Msg("bulk querying")
	responses, err := client.BulkGeocode(requests)
	if err != nil {
		return fmt.Errorf("client bulk geocode: %w", err)
	}
	log.Info().Int("len", len(responses)).Msg("bulk query response")
	//setters := make([]*models.FileuploadPoolSetter, 0)
	for i, resp := range responses {
		pool := pools[i]
		if resp.Status != 200 {
			log.Warn().Int("row", i).Int("status", resp.Status).Str("pool.address", pool.AddressStreet).Str("pool.postal", pool.AddressPostalCode).Str("msg", resp.Message)
		}
		if resp.Response != nil {
			for _, feature := range resp.Response.Features {
				log.Info().Int("row", i).Int("status", resp.Status).Str("pool.address", pool.AddressStreet).Str("pool.postal", pool.AddressPostalCode).Str("feature.type", feature.Type).Str("formatted", feature.Properties.FormattedAddressLine).Float64("coordinates[0]", feature.Geometry.Coordinates[0]).Float64("coordinates[1]", feature.Geometry.Coordinates[1]).Msg("geocode")
			}
		}
	}
	return nil
}
func parseFile(ctx context.Context, txn bob.Tx, file models.FileuploadFile) ([]*models.FileuploadPool, error) {
	pools := make([]*models.FileuploadPool, 0)
	c, err := models.FindFileuploadCSV(ctx, db.PGInstance.BobDB, file.ID)
	if err != nil {
		return pools, fmt.Errorf("Failed to get file %d from DB: %w", file.ID, err)
	}
	r, err := userfile.NewFileReader(userfile.CollectionCSV, file.FileUUID)
	if err != nil {
		return pools, fmt.Errorf("Failed to get filereader for %d: %w", file.ID, err)
	}
	reader := csv.NewReader(r)
	h, err := reader.Read()
	if err != nil {
		return pools, fmt.Errorf("Failed to read header of CSV for file %d: %w", file.ID, err)
	}
	headers := parseHeaders(h)
	missing_headers := missingRequiredHeaders(headers)
	for _, mh := range missing_headers {
		errorMissingHeader(ctx, txn, c, mh)
		file.Update(ctx, txn, &models.FileuploadFileSetter{
			Status: omit.From(enums.FileuploadFilestatustypeError),
		})
		txn.Commit(ctx)
		return pools, nil
	}
	row_number := 0
	for {
		row, err := reader.Read()
		if err != nil {
			if err == io.EOF {
				file.Update(ctx, txn, &models.FileuploadFileSetter{
					Status: omit.From(enums.FileuploadFilestatustypeParsed),
				})
				log.Info().Int32("file.ID", file.ID).Msg("Set file to parsed")
				txn.Commit(ctx)
				return pools, nil
			}
			return pools, fmt.Errorf("Failed to read all CSV records for file %d: %w", file.ID, err)
		}
		setter := models.FileuploadPoolSetter{
			// required fields
			//AddressCity: omit.From(),
			//AddressPostalCode: omit.From(),
			//AddressStreet: omit.From(),
			Committed: omit.From(false),
			Condition: omit.From(enums.FileuploadPoolconditiontypeUnknown),
			Created:   omit.From(time.Now()),
			CreatorID: omit.From(file.CreatorID),
			CSVFile:   omit.From(file.ID),
			Deleted:   omitnull.FromPtr[time.Time](nil),
			Geom:      omitnull.FromPtr[string](nil),
			H3cell:    omitnull.FromPtr[string](nil),
			// ID - generated
			IsInDistrict:           omit.From(false),
			IsNew:                  omit.From(false),
			Notes:                  omit.From(""),
			OrganizationID:         omit.From(file.OrganizationID),
			PropertyOwnerName:      omit.From(""),
			PropertyOwnerPhoneE164: omitnull.FromPtr[string](nil),
			ResidentOwned:          omitnull.FromPtr[bool](nil),
			ResidentPhoneE164:      omitnull.FromPtr[string](nil),
			Version:                omit.From(int32(0)),
		}
		for i, col := range row {
			hdr := headers[i]
			if col == "" {
				continue
			}
			switch hdr {
			case headerAddressCity:
				setter.AddressCity = omit.From(col)
			case headerAddressPostalCode:
				setter.AddressPostalCode = omit.From(col)
			case headerAddressStreet:
				setter.AddressStreet = omit.From(col)
			case headerCondition:
				var condition enums.FileuploadPoolconditiontype
				err := condition.Scan(strings.ToLower(col))
				if err != nil {
					addError(ctx, txn, c, int32(row_number), int32(i), fmt.Sprintf("'%s' is not a pool condition that we recognize. It should be one of %s", col, poolConditionValidValues()))
					continue
				}
				setter.Condition = omit.From(condition)
			case headerNotes:
				setter.Notes = omit.From(col)
			case headerPropertyOwnerName:
				setter.PropertyOwnerName = omit.From(col)
			case headerPropertyOwnerPhone:
				phone, err := text.ParsePhoneNumber(col)
				if err != nil {
					addError(ctx, txn, c, int32(row_number), int32(i), fmt.Sprintf("'%s' is not a phone number that we recognize. Ideally it should be of the form '+12223334444'", col))
					continue
				}
				text.EnsureInDB(ctx, txn, *phone)
				setter.PropertyOwnerPhoneE164 = omitnull.From(text.PhoneString(*phone))
			case headerResidentOwned:
				boolValue, err := parseBool(col)
				if err != nil {
					addError(ctx, txn, c, int32(row_number), int32(i), fmt.Sprintf("'%s' is not something that we recognize as a true/false value. Please use either 'true' or 'false'", col))
					continue
				}
				setter.ResidentOwned = omitnull.From(boolValue)
			case headerResidentPhone:
				phone, err := text.ParsePhoneNumber(col)
				if err != nil {
					addError(ctx, txn, c, int32(row_number), int32(i), fmt.Sprintf("'%s' is not a phone number that we recognize. Ideally it should be of the form '+12223334444'", col))
					continue
				}
				text.EnsureInDB(ctx, txn, *phone)
				setter.ResidentPhoneE164 = omitnull.From(text.PhoneString(*phone))
			}
		}
		pool, err := models.FileuploadPools.Insert(&setter).One(ctx, txn)
		if err != nil {
			return pools, fmt.Errorf("Failed to create pool: %w", err)
		}
		pools = append(pools, pool)
		row_number = row_number + 1
	}
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
func parseHeaders(row []string) []headerPoolEnum {
	results := make([]headerPoolEnum, 0)
	for _, h := range row {
		ht := strings.TrimSpace(h)
		hl := strings.ToLower(ht)
		log.Debug().Str("header", hl).Msg("Saw CSV header")
		var type_ headerPoolEnum
		switch hl {
		case "city":
			type_ = headerAddressCity
		case "zip":
		case "postal code":
			type_ = headerAddressPostalCode
		case "street address":
			type_ = headerAddressStreet
		case "condition":
		case "pool condition":
			type_ = headerCondition
		case "notes":
			type_ = headerNotes
		case "property owner":
		case "property owner name":
			type_ = headerPropertyOwnerName
		case "property owner phone":
			type_ = headerPropertyOwnerPhone
		case "resident owned":
			type_ = headerResidentOwned
		case "resident phone":
		case "resident phone number":
			type_ = headerResidentPhone
		default:
			type_ = headerTag
		}
		results = append(results, type_)
	}

	return results
}
func missingRequiredHeaders(headers []headerPoolEnum) []headerPoolEnum {
	results := make([]headerPoolEnum, 0)
	for _, rh := range []headerPoolEnum{headerAddressCity, headerAddressPostalCode, headerAddressStreet} {
		present := false
		for _, h := range headers {
			if h == rh {
				present = true
				break
			}
		}
		if !present {
			results = append(results, rh)
		}
	}
	return results
}
func poolConditionValidValues() string {
	var b strings.Builder
	for i, cond := range enums.AllFileuploadPoolconditiontype() {
		if i == 0 {
			fmt.Fprintf(&b, "'%s'", cond)
		} else {
			fmt.Fprintf(&b, ", '%s'", cond)
		}
	}
	return b.String()
}
