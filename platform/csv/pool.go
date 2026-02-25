package csv

import (
	"context"
	"encoding/csv"
	"fmt"
	"io"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/Gleipnir-Technology/bob"
	"github.com/Gleipnir-Technology/bob/dialect/psql"
	"github.com/Gleipnir-Technology/bob/dialect/psql/um"
	"github.com/Gleipnir-Technology/nidus-sync/config"
	"github.com/Gleipnir-Technology/nidus-sync/db"
	"github.com/Gleipnir-Technology/nidus-sync/db/enums"
	"github.com/Gleipnir-Technology/nidus-sync/db/models"
	"github.com/Gleipnir-Technology/nidus-sync/h3utils"
	"github.com/Gleipnir-Technology/nidus-sync/platform/geom"
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
	c, err := models.FindFileuploadCSV(ctx, txn, file.ID)
	if err != nil {
		return fmt.Errorf("Failed to get csv file %d from DB: %w", file.ID, err)
	}
	pools, err := parseFile(ctx, txn, file, c)
	if err != nil {
		return fmt.Errorf("parse file: %w", err)
	}
	_, err = psql.Update(
		um.Table("fileupload.csv"),
		um.SetCol("rowcount").ToArg(len(pools)),
		um.Where(psql.Quote("file_id").EQ(psql.Arg(file_id))),
	).Exec(ctx, txn)
	if err != nil {
		return fmt.Errorf("update csv row: %w", err)
	}
	org, err := models.FindOrganization(ctx, db.PGInstance.BobDB, file.OrganizationID)
	if err != nil {
		return fmt.Errorf("get org: %w", err)
	}
	err = bulkGeocode(ctx, txn, file, c, pools, org)
	if err != nil {
		log.Error().Err(err).Msg("Failure during geocoding")
	}
	file.Update(ctx, txn, &models.FileuploadFileSetter{
		Status: omit.From(enums.FileuploadFilestatustypeParsed),
	})
	log.Info().Int32("file.ID", file.ID).Msg("Set file to parsed")
	txn.Commit(ctx)
	return nil
}
func bulkGeocode(ctx context.Context, txn bob.Tx, file *models.FileuploadFile, c *models.FileuploadCSV, pools []*models.FileuploadPool, org *models.Organization) error {
	if len(pools) == 0 {
		return nil
	}
	log.Info().Int("len pools", len(pools)).Msg("bulk geocoding")
	client := stadia.NewStadiaMaps(config.StadiaMapsAPIKey)
	jobs := make(chan *jobGeocode, len(pools))
	errors := make(chan error, len(pools))

	var wg sync.WaitGroup
	for i := 0; i < 20; i++ {
		wg.Add(1)
		go worker(ctx, txn, client, jobs, errors, &wg)
	}

	for i, pool := range pools {
		jobs <- &jobGeocode{
			csv:       c,
			rownumber: int32(i),
			org:       org,
			pool:      pool,
		}
	}
	close(jobs)

	go func() {
		wg.Wait()
		close(errors)
	}()

	error_count := 0
	for err := range errors {
		log.Error().Err(err).Msg("failed to geocode")
		error_count++
	}
	if error_count > 0 {
		txn.Rollback(ctx)
		return fmt.Errorf("%d errors encountered in bulk geocode", error_count)
	}
	update_query := `
	UPDATE fileupload.pool p
		SET is_in_district = (
    			EXISTS (
        			SELECT 1
        			FROM import.district d
				JOIN organization o ON d.gid = o.import_district_gid
				WHERE o.id = p.organization_id
					AND ST_Contains(d.geom_4326, p.geom)
			)
		)
	WHERE p.geom IS NOT NULL;`
	_, err := txn.ExecContext(ctx, update_query)
	if err != nil {
		return fmt.Errorf("failed to update is_in_district: %w", err)
	}
	return nil
}

type jobGeocode struct {
	csv       *models.FileuploadCSV
	rownumber int32
	org       *models.Organization
	pool      *models.FileuploadPool
}

func geocode(ctx context.Context, txn bob.Tx, client *stadia.StadiaMaps, job *jobGeocode) error {
	pool := job.pool
	sublog := log.With().
		Str("pool.address_postal", pool.AddressPostalCode).
		Str("pool.address_street", pool.AddressStreet).
		Str("pool.postal", pool.AddressPostalCode).
		Logger()
	req := stadia.StructuredGeocodeRequest{
		Address:    &pool.AddressStreet,
		Locality:   &pool.AddressCity,
		PostalCode: &pool.AddressPostalCode,
	}
	maybeAddServiceArea(&req, job.org)
	resp, err := client.StructuredGeocode(ctx, req)
	if err != nil {
		return fmt.Errorf("client structured geocode failure on %s, %s, %s: %w", pool.AddressStreet, pool.AddressCity, pool.AddressPostalCode, err)
	}
	if len(resp.Features) > 1 {
		sublog.Warn().Int("len", len(resp.Features)).Msg("More than one feature")
		addError(ctx, txn, job.csv, job.rownumber, 0, "The address provided matched more than one location")
	}
	feature := resp.Features[0]
	if feature.Geometry.Type != "Point" {
		return fmt.Errorf("wrong type %s from %s %s", feature.Geometry.Type, pool.AddressStreet, pool.AddressPostalCode)
	}
	longitude := feature.Geometry.Coordinates[0]
	latitude := feature.Geometry.Coordinates[1]
	cell, err := h3utils.GetCell(longitude, latitude, 15)
	if err != nil {
		return fmt.Errorf("failed to convert lat %f lng %f to h3 cell", longitude, latitude)
	}
	geom_query := geom.PostgisPointQuery(longitude, latitude)
	_, err = psql.Update(
		um.Table("fileupload.pool"),
		um.SetCol("h3cell").ToArg(cell),
		um.SetCol("geom").To(geom_query),
		um.Where(psql.Quote("id").EQ(psql.Arg(pool.ID))),
	).Exec(ctx, txn)
	if err != nil {
		return fmt.Errorf("failed to update pool: %w", err)
	}
	return nil
}
func parseFile(ctx context.Context, txn bob.Tx, file *models.FileuploadFile, c *models.FileuploadCSV) ([]*models.FileuploadPool, error) {
	pools := make([]*models.FileuploadPool, 0)
	r, err := userfile.NewFileReader(userfile.CollectionCSV, file.FileUUID)
	if err != nil {
		return pools, fmt.Errorf("Failed to get filereader for %d: %w", file.ID, err)
	}
	reader := csv.NewReader(r)
	h, err := reader.Read()
	if err != nil {
		return pools, fmt.Errorf("Failed to read header of CSV for file %d: %w", file.ID, err)
	}
	header_types, header_names := parseHeaders(h)
	missing_headers := missingRequiredHeaders(header_types)
	for _, mh := range missing_headers {
		errorMissingHeader(ctx, txn, c, mh)
		file.Update(ctx, txn, &models.FileuploadFileSetter{
			Status: omit.From(enums.FileuploadFilestatustypeError),
		})
		return pools, nil
	}
	// Start at 2 because the header is line 1, not line 0
	line_number := int32(2)
	for {
		row, err := reader.Read()
		if err != nil {
			if err == io.EOF {
				return pools, nil
			}
			return pools, fmt.Errorf("Failed to read all CSV records for file %d: %w", file.ID, err)
		}
		tags := make(map[string]string, 0)
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
			IsNew:                  omit.From(true),
			LineNumber:             omit.From(line_number),
			Notes:                  omit.From(""),
			OrganizationID:         omit.From(file.OrganizationID),
			PropertyOwnerName:      omit.From(""),
			PropertyOwnerPhoneE164: omitnull.FromPtr[string](nil),
			ResidentOwned:          omitnull.FromPtr[bool](nil),
			ResidentPhoneE164:      omitnull.FromPtr[string](nil),
			//Tags:       		convertToPGData(tags),
		}
		for i, col := range row {
			hdr_t := header_types[i]
			if col == "" {
				continue
			}
			switch hdr_t {
			case headerAddressCity:
				setter.AddressCity = omit.From(col)
			case headerAddressPostalCode:
				setter.AddressPostalCode = omit.From(col)
			case headerAddressStreet:
				setter.AddressStreet = omit.From(col)
			case headerCondition:
				var condition enums.FileuploadPoolconditiontype
				col_l := strings.ToLower(col)
				col_translated := col_l
				switch col_l {
				case "empty":
					col_translated = "dry"
				}
				err := condition.Scan(col_translated)
				if err != nil {
					addError(ctx, txn, c, int32(line_number), int32(i), fmt.Sprintf("'%s' is not a pool condition that we recognize. It should be one of %s", col, poolConditionValidValues()))
					setter.Condition = omit.From(enums.FileuploadPoolconditiontypeUnknown)
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
					addError(ctx, txn, c, int32(line_number), int32(i), fmt.Sprintf("'%s' is not a phone number that we recognize. Ideally it should be of the form '+12223334444'", col))
					continue
				}
				text.EnsureInDB(ctx, txn, *phone)
				setter.PropertyOwnerPhoneE164 = omitnull.From(text.PhoneString(*phone))
			case headerResidentOwned:
				boolValue, err := parseBool(col)
				if err != nil {
					addError(ctx, txn, c, int32(line_number), int32(i), fmt.Sprintf("'%s' is not something that we recognize as a true/false value. Please use either 'true' or 'false'", col))
					continue
				}
				setter.ResidentOwned = omitnull.From(boolValue)
			case headerResidentPhone:
				phone, err := text.ParsePhoneNumber(col)
				if err != nil {
					addError(ctx, txn, c, int32(line_number), int32(i), fmt.Sprintf("'%s' is not a phone number that we recognize. Ideally it should be of the form '+12223334444'", col))
					continue
				}
				text.EnsureInDB(ctx, txn, *phone)
				setter.ResidentPhoneE164 = omitnull.From(text.PhoneString(*phone))
			case headerTag:
				tags[header_names[i]] = col
			}

		}
		setter.Tags = omit.From(db.ConvertToPGData(tags))
		pool, err := models.FileuploadPools.Insert(&setter).One(ctx, txn)
		if err != nil {
			return pools, fmt.Errorf("Failed to create pool: %w", err)
		}
		pools = append(pools, pool)
		line_number = line_number + 1
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
func maybeAddServiceArea(req *stadia.StructuredGeocodeRequest, org *models.Organization) {
	if org.ServiceAreaXmax.IsNull() ||
		org.ServiceAreaYmax.IsNull() ||
		org.ServiceAreaXmin.IsNull() ||
		org.ServiceAreaYmin.IsNull() {
		return
	}
	xmax := org.ServiceAreaXmax.MustGet()
	ymax := org.ServiceAreaYmax.MustGet()
	xmin := org.ServiceAreaXmin.MustGet()
	ymin := org.ServiceAreaYmin.MustGet()
	req.BoundaryRectMaxLon = &xmax
	req.BoundaryRectMaxLat = &ymax
	req.BoundaryRectMinLon = &xmin
	req.BoundaryRectMinLat = &ymin
}
func parseHeaders(row []string) ([]headerPoolEnum, []string) {
	result_enums := make([]headerPoolEnum, 0)
	result_names := make([]string, 0)
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
		result_enums = append(result_enums, type_)
		result_names = append(result_names, hl)
	}

	return result_enums, result_names
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
func worker(ctx context.Context, txn bob.Tx, client *stadia.StadiaMaps, jobs <-chan *jobGeocode, errors chan<- error, wg *sync.WaitGroup) {
	defer wg.Done()

	for job := range jobs {
		err := geocode(ctx, txn, client, job)

		if err != nil {
			errors <- err
		}
	}
}
