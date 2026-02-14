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
	"github.com/Gleipnir-Technology/nidus-sync/db"
	"github.com/Gleipnir-Technology/nidus-sync/db/enums"
	"github.com/Gleipnir-Technology/nidus-sync/db/models"
	"github.com/Gleipnir-Technology/nidus-sync/platform/text"
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
	c, err := models.FindFileuploadCSV(ctx, db.PGInstance.BobDB, file_id)
	if err != nil {
		return fmt.Errorf("Failed to get file %d from DB: %w", file_id, err)
	}
	r, err := userfile.NewFileReader(userfile.CollectionCSV, file.FileUUID)
	if err != nil {
		return fmt.Errorf("Failed to get filereader for %d: %w", file_id, err)
	}
	reader := csv.NewReader(r)
	h, err := reader.Read()
	if err != nil {
		return fmt.Errorf("Failed to read header of CSV for file %d: %w", file_id, err)
	}
	txn, err := db.PGInstance.BobDB.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("Failed to start transaction: %w", err)
	}
	defer txn.Rollback(ctx)
	headers := parseHeaders(h)
	missing_headers := missingRequiredHeaders(headers)
	for _, mh := range missing_headers {
		errorMissingHeader(ctx, txn, c, mh)
		file.Update(ctx, txn, &models.FileuploadFileSetter{
			Status: omit.From(enums.FileuploadFilestatustypeError),
		})
		txn.Commit(ctx)
		return nil
	}
	row_number := 0
	for {
		row, err := reader.Read()
		if err != nil {
			if err == io.EOF {
				file.Update(ctx, txn, &models.FileuploadFileSetter{
					Status: omit.From(enums.FileuploadFilestatustypeParsed),
				})
				log.Info().Int32("file_id", file_id).Msg("Set file to parsed")
				txn.Commit(ctx)
				return nil
			}
			return fmt.Errorf("Failed to read all CSV records for file %d: %w", file_id, err)
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
		_, err = models.FileuploadPools.Insert(&setter).Exec(ctx, txn)
		if err != nil {
			return fmt.Errorf("Failed to create pool: %w", err)
		}
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
