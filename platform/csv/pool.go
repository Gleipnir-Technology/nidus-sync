package csv

import (
	"context"
	"encoding/csv"
	"fmt"
	"io"
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
	"github.com/Gleipnir-Technology/nidus-sync/platform/file"
	"github.com/Gleipnir-Technology/nidus-sync/platform/geocode"
	"github.com/Gleipnir-Technology/nidus-sync/platform/geom"
	"github.com/Gleipnir-Technology/nidus-sync/platform/text"
	"github.com/Gleipnir-Technology/nidus-sync/platform/types"
	"github.com/Gleipnir-Technology/nidus-sync/stadia"
	"github.com/aarondl/opt/omit"
	"github.com/aarondl/opt/omitnull"
	"github.com/rs/zerolog/log"
)

type headerPoolEnum int

const (
	headerPoolUnknown headerPoolEnum = iota
	headerPoolAddressLocality
	headerPoolAddressPostalCode
	headerPoolAddressRegion
	headerPoolAddressStreet
	headerPoolCondition
	headerPoolNotes
	headerPoolPropertyOwnerName
	headerPoolPropertyOwnerPhone
	headerPoolResidentOwned
	headerPoolResidentPhone
	headerPoolTag
)

func (e headerPoolEnum) String() string {
	switch e {
	case headerPoolAddressPostalCode:
		return "Postal Code"
	case headerPoolAddressLocality:
		return "City"
	case headerPoolAddressStreet:
		return "Street Address"
	case headerPoolCondition:
		return "Condition"
	case headerPoolNotes:
		return "Notes"
	case headerPoolPropertyOwnerName:
		return "Property Owner Name"
	case headerPoolPropertyOwnerPhone:
		return "Property Owner Phone"
	case headerPoolResidentOwned:
		return "Resident Owned"
	case headerPoolResidentPhone:
		return "Resident Phone"
	case headerPoolAddressRegion:
		return "State"
	default:
		return "bad programmer"
	}
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
        			FROM organization o, fileupload.file f
				WHERE
					p.csv_file = f.id AND
					f.organization_id = o.id AND (
						ST_Contains(o.service_area_geometry, p.geom) OR
						o.is_catchall
					)
						
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

func geocodePool(ctx context.Context, txn bob.Tx, client *stadia.StadiaMaps, job *jobGeocode) error {
	pool := job.pool
	a := types.Address{
		Number:     pool.AddressNumber,
		Locality:   pool.AddressLocality,
		PostalCode: pool.AddressPostalCode,
		Region:     pool.AddressRegion,
		Street:     pool.AddressStreet,
	}
	geo, err := geocode.GeocodeStructured(ctx, job.org, a)
	if err != nil {
		addError(ctx, txn, job.csv, job.rownumber, 0, err.Error())
		return nil
	}
	if geo.Address.Location == nil {
		addError(ctx, txn, job.csv, job.rownumber, 0, fmt.Sprintf("nil location from geocoding"))
		return nil
	}
	geom_query := geom.PostgisPointQuery(*geo.Address.Location)
	_, err = psql.Update(
		um.Table("fileupload.pool"),
		um.SetCol("h3cell").ToArg(geo.Cell),
		um.SetCol("geom").To(geom_query),
		um.SetCol("address_id").To(*geo.Address.ID),
		um.Where(psql.Quote("id").EQ(psql.Arg(pool.ID))),
	).Exec(ctx, txn)
	if err != nil {
		return fmt.Errorf("failed to update pool: %w", err)
	}
	return nil
}
func parseCSVPoollist(ctx context.Context, txn bob.Tx, f *models.FileuploadFile, c *models.FileuploadCSV) ([]*models.FileuploadPool, error) {
	pools := make([]*models.FileuploadPool, 0)
	r, err := file.NewFileReader(file.CollectionCSV, f.FileUUID)
	if err != nil {
		return pools, fmt.Errorf("Failed to get filereader for %d: %w", f.ID, err)
	}
	reader := csv.NewReader(r)
	h, err := reader.Read()
	if err != nil {
		return pools, fmt.Errorf("Failed to read header of CSV for file %d: %w", f.ID, err)
	}
	header_types, header_names := parseHeaders(h)
	missing_headers := missingRequiredHeaders(header_types)
	for _, mh := range missing_headers {
		errorMissingHeader(ctx, txn, c, mh)
		f.Update(ctx, txn, &models.FileuploadFileSetter{
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
			return pools, fmt.Errorf("Failed to read all CSV records for file %d: %w", f.ID, err)
		}
		tags := make(map[string]string, 0)
		setter := models.FileuploadPoolSetter{
			// required fields
			//AddressNumber: omit.From(),
			//AddressLocality: omit.From(),
			//AddressPostalCode: omit.From(),
			//AddressRegion: omit.From(),
			//AddressStreet: omit.From(),
			Committed: omit.From(false),
			Condition: omit.From(enums.PoolconditiontypeUnknown),
			Created:   omit.From(time.Now()),
			CreatorID: omit.From(f.CreatorID),
			CSVFile:   omit.From(f.ID),
			Deleted:   omitnull.FromPtr[time.Time](nil),
			Geom:      omitnull.FromPtr[string](nil),
			H3cell:    omitnull.FromPtr[string](nil),
			// ID - generated
			IsInDistrict:           omit.From(false),
			IsNew:                  omit.From(true),
			LineNumber:             omit.From(line_number),
			Notes:                  omit.From(""),
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
			case headerPoolAddressLocality:
				setter.AddressLocality = omit.From(col)
			case headerPoolAddressPostalCode:
				setter.AddressPostalCode = omit.From(col)
			case headerPoolAddressRegion:
				setter.AddressRegion = omit.From(col)
			case headerPoolAddressStreet:
				// This type of spreadsheet normally has '123 Main Str'
				parts := strings.SplitN(col, " ", 2)
				if len(parts) != 2 {
					addError(ctx, txn, c, int32(line_number), int32(i), fmt.Sprintf("'%s' is not a house number and street. It needs to be in the form '123 main'", col))
					continue
				}
				setter.AddressNumber = omit.From(parts[0])
				setter.AddressStreet = omit.From(parts[1])
			case headerPoolCondition:
				var condition enums.Poolconditiontype
				col_l := strings.ToLower(col)
				col_translated := col_l
				switch col_l {
				case "empty":
					col_translated = "dry"
				}
				err := condition.Scan(col_translated)
				if err != nil {
					addError(ctx, txn, c, int32(line_number), int32(i), fmt.Sprintf("'%s' is not a pool condition that we recognize. It should be one of %s", col, poolConditionValidValues()))
					setter.Condition = omit.From(enums.PoolconditiontypeUnknown)
					continue
				}
				setter.Condition = omit.From(condition)
			case headerPoolNotes:
				setter.Notes = omit.From(col)
			case headerPoolPropertyOwnerName:
				setter.PropertyOwnerName = omit.From(col)
			case headerPoolPropertyOwnerPhone:
				phone, err := text.ParsePhoneNumber(col)
				if err != nil {
					addError(ctx, txn, c, int32(line_number), int32(i), fmt.Sprintf("'%s' is not a phone number that we recognize. Ideally it should be of the form '+12223334444'", col))
					continue
				}
				text.EnsureInDB(ctx, txn, *phone)
				setter.PropertyOwnerPhoneE164 = omitnull.From(phone.PhoneString())
			case headerPoolResidentOwned:
				boolValue, err := parseBool(col)
				if err != nil {
					addError(ctx, txn, c, int32(line_number), int32(i), fmt.Sprintf("'%s' is not something that we recognize as a true/false value. Please use either 'true' or 'false'", col))
					continue
				}
				setter.ResidentOwned = omitnull.From(boolValue)
			case headerPoolResidentPhone:
				phone, err := text.ParsePhoneNumber(col)
				if err != nil {
					addError(ctx, txn, c, int32(line_number), int32(i), fmt.Sprintf("'%s' is not a phone number that we recognize. Ideally it should be of the form '+12223334444'", col))
					continue
				}
				text.EnsureInDB(ctx, txn, *phone)
				setter.ResidentPhoneE164 = omitnull.From(phone.PhoneString())
			case headerPoolTag:
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
func processCSVPoollist(ctx context.Context, txn bob.Tx, f *models.FileuploadFile, c *models.FileuploadCSV, parsed []*models.FileuploadPool) error {
	org, err := models.FindOrganization(ctx, db.PGInstance.BobDB, f.OrganizationID)
	if err != nil {
		return fmt.Errorf("get org: %w", err)
	}
	err = bulkGeocode(ctx, txn, f, c, parsed, org)
	if err != nil {
		log.Error().Err(err).Msg("Failure during geocoding")
	}
	return nil
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
			type_ = headerPoolAddressLocality
		case "zip":
		case "postal code":
			type_ = headerPoolAddressPostalCode
		case "state":
			type_ = headerPoolAddressRegion
		case "street address":
			type_ = headerPoolAddressStreet
		case "condition":
		case "pool condition":
			type_ = headerPoolCondition
		case "notes":
			type_ = headerPoolNotes
		case "property owner":
		case "property owner name":
			type_ = headerPoolPropertyOwnerName
		case "property owner phone":
			type_ = headerPoolPropertyOwnerPhone
		case "resident owned":
			type_ = headerPoolResidentOwned
		case "resident phone":
		case "resident phone number":
			type_ = headerPoolResidentPhone
		default:
			type_ = headerPoolTag
		}
		result_enums = append(result_enums, type_)
		result_names = append(result_names, hl)
	}

	return result_enums, result_names
}
func missingRequiredHeaders(headers []headerPoolEnum) []headerPoolEnum {
	results := make([]headerPoolEnum, 0)
	for _, rh := range []headerPoolEnum{headerPoolAddressLocality, headerPoolAddressRegion, headerPoolAddressPostalCode, headerPoolAddressStreet} {
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
	for i, cond := range enums.AllPoolconditiontype() {
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
		err := geocodePool(ctx, txn, client, job)

		if err != nil {
			errors <- err
		}
	}
}
