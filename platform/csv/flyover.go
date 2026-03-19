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
	"github.com/Gleipnir-Technology/bob/dialect/psql"
	"github.com/Gleipnir-Technology/bob/dialect/psql/um"
	"github.com/Gleipnir-Technology/nidus-sync/db"
	"github.com/Gleipnir-Technology/nidus-sync/db/enums"
	"github.com/Gleipnir-Technology/nidus-sync/db/models"
	"github.com/Gleipnir-Technology/nidus-sync/h3utils"
	"github.com/Gleipnir-Technology/nidus-sync/platform/file"
	"github.com/Gleipnir-Technology/nidus-sync/platform/geom"
	"github.com/Gleipnir-Technology/nidus-sync/platform/types"
	"github.com/aarondl/opt/omit"
	"github.com/aarondl/opt/omitnull"
	"github.com/rs/zerolog/log"
)

type Enum interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 | ~string
}

type headerFlyoverEnum int

const (
	headerFlyoverAddressLocality headerFlyoverEnum = iota
	headerFlyoverAddressNumber
	headerFlyoverAddressPostalCode
	headerFlyoverAddressRegion
	headerFlyoverAddressStreet
	headerFlyoverComment
	headerFlyoverLatitude
	headerFlyoverLongitude
	headerFlyoverNone
)

func (e headerFlyoverEnum) String() string {
	switch e {
	case headerFlyoverAddressLocality:
		return "City"
	case headerFlyoverAddressNumber:
		return "HouseNo"
	case headerFlyoverAddressPostalCode:
		return "ZIP"
	case headerFlyoverAddressRegion:
		return "State"
	case headerFlyoverAddressStreet:
		return "Street"
	case headerFlyoverComment:
		return "Comment"
	case headerFlyoverLatitude:
		return "TargetLat"
	case headerFlyoverLongitude:
		return "TargetLon"
	default:
		return "bad programmer"
	}
}

var parseCSVFlyover = makeParseCSV(
	makeParseHeaders(map[string]headerFlyoverEnum{
		"comment":   headerFlyoverComment,
		"houseno":   headerFlyoverAddressNumber,
		"state":     headerFlyoverAddressRegion,
		"street":    headerFlyoverAddressStreet,
		"city":      headerFlyoverAddressLocality,
		"targetlat": headerFlyoverLatitude,
		"targetlon": headerFlyoverLongitude,
		"zip":       headerFlyoverAddressPostalCode,
		"*":         headerFlyoverNone,
	}),
	insertFlyover,
)

type insertModelFunc[ModelType any, HeaderType Enum] = func(context.Context, bob.Tx, *models.FileuploadFile, *models.FileuploadCSV, int32, []HeaderType, []string, []string) (ModelType, error)
type parseCSVFunc[ModelType any] = func(ctx context.Context, txn bob.Tx, f *models.FileuploadFile, c *models.FileuploadCSV) ([]ModelType, error)

func makeParseCSV[ModelType any, HeaderType Enum](parseHeader parseHeaderFunc[HeaderType], insertModel insertModelFunc[ModelType, HeaderType]) parseCSVFunc[ModelType] {
	return func(ctx context.Context, txn bob.Tx, f *models.FileuploadFile, c *models.FileuploadCSV) ([]ModelType, error) {
		rows := make([]ModelType, 0)
		r, err := file.NewFileReader(file.CollectionCSV, f.FileUUID)
		if err != nil {
			return rows, fmt.Errorf("Failed to get filereader for %d: %w", f.ID, err)
		}
		reader := csv.NewReader(r)
		h, err := reader.Read()
		if err != nil {
			return rows, fmt.Errorf("Failed to read header of CSV for file %d: %w", f.ID, err)
		}
		header_types, header_names := parseHeader(h)
		/*
			TODO: Add support for missing headersi
			missing_headers := missingRequiredHeaders(header_types)
			for _, mh := range missing_headers {
				errorMissingHeader(ctx, txn, c, mh)
				file.Update(ctx, txn, &models.FileuploadFileSetter{
					Status: omit.From(enums.FileuploadFilestatustypeError),
				})
				return pools, nil
			}
		*/
		// Start at 2 because the header is line 1, not line 0
		line_number := int32(2)
		for {
			row, err := reader.Read()
			if err != nil {
				if err == io.EOF {
					return rows, nil
				}
				return rows, fmt.Errorf("Failed to read all CSV records for file %d: %w", f.ID, err)
			}
			m, err := insertModel(ctx, txn, f, c, line_number, header_types, header_names, row)
			if err != nil {
				return rows, fmt.Errorf("insert models: %w", err)
			}
			rows = append(rows, m)
			line_number = line_number + 1
		}
	}
}
func insertFlyover(ctx context.Context, txn bob.Tx, file *models.FileuploadFile, c *models.FileuploadCSV, line_number int32, header_types []headerFlyoverEnum, header_names []string, row []string) (*models.FileuploadPool, error) {
	/*
		setter := models.FileuploadFlyoverAerialServiceSetter{
			Committed: omit.From(false),
			Condition: omit.From(enums.FileuploadPoolconditiontypeUnknown),
			Created:   omit.From(time.Now()),
			CreatorID: omit.From(file.CreatorID),
			CSVFile:   omit.From(file.ID),
			Deleted:   omitnull.FromPtr[time.Time](nil),
			Geom:      omitnull.FromPtr[string](nil),
			H3cell:    omitnull.FromPtr[string](nil),
			// ID - generated
			OrganizationID: omit.From(file.OrganizationID),
		}
	*/
	setter := models.FileuploadPoolSetter{
		// required fields
		//AddressLocality: omit.From(),
		//AddressPostalCode: omit.From(),
		//AddressRegion: omit.From(),
		//AddressStreet: omit.From(),
		Committed: omit.From(false),
		Condition: omit.From(enums.PoolconditiontypeUnknown),
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
		PropertyOwnerName:      omit.From(""),
		PropertyOwnerPhoneE164: omitnull.FromPtr[string](nil),
		ResidentOwned:          omitnull.FromPtr[bool](nil),
		ResidentPhoneE164:      omitnull.FromPtr[string](nil),
		//Tags:       		convertToPGData(tags),
	}
	var lat, lng float64
	var err error
	for i, value := range row {
		if value == "" {
			continue
		}
		header_type := header_types[i]
		switch header_type {
		case headerFlyoverAddressLocality:
			setter.AddressLocality = omit.From(value)
		case headerFlyoverAddressNumber:
			setter.AddressNumber = omit.From(value)
		case headerFlyoverAddressPostalCode:
			setter.AddressPostalCode = omit.From(value)
		case headerFlyoverAddressRegion:
			setter.AddressRegion = omit.From(value)
		case headerFlyoverAddressStreet:
			setter.AddressStreet = omit.From(value)
		case headerFlyoverComment:
			condition, err := parsePoolCondition(value)
			if err == nil {
				setter.Condition = omit.From(condition)
			} else {
				addError(ctx, txn, c, int32(line_number), int32(i), fmt.Sprintf("'%s' is not a pool condition that we recognize. It should be one of %s", value, poolConditionValidValues()))
				continue
			}
		case headerFlyoverLatitude:
			lat, err = strconv.ParseFloat(value, 10)
			if err != nil {
				addError(ctx, txn, c, int32(line_number), int32(i), fmt.Sprintf("'%s' is not decimal value", value))
				continue
			}
		case headerFlyoverLongitude:
			lng, err = strconv.ParseFloat(value, 10)
			if err != nil {
				addError(ctx, txn, c, int32(line_number), int32(i), fmt.Sprintf("'%s' is not decimal value", value))
				continue
			}
		}
	}
	setter.Tags = omit.From(db.ConvertToPGData(map[string]string{}))
	flyover, err := models.FileuploadPools.Insert(&setter).One(ctx, txn)
	if err != nil {
		return nil, fmt.Errorf("Failed to create flyover: %w", err)
	}
	cell, err := h3utils.GetCell(lng, lat, 15)
	if err != nil {
		return nil, fmt.Errorf("failed to convert lat %f lng %f to h3 cell", lng, lat)
	}
	geom_query := geom.PostgisPointQuery(types.Location{
		Latitude:  lat,
		Longitude: lng,
	})
	_, err = psql.Update(
		um.TableAs("fileupload.pool", "pool"),
		um.SetCol("h3cell").ToArg(cell),
		um.SetCol("geom").To(geom_query),
		um.SetCol("is_in_district").To(
			psql.F("COALESCE",
				psql.F("ST_Contains", "org.service_area_geometry", geom_query),
				psql.Quote("org", "is_catchall"),
			),
		),
		um.From("fileupload.csv").As("csv"),
		um.InnerJoin("fileupload.file").As("file").OnEQ(psql.Quote("csv", "file_id"), psql.Quote("file", "id")),
		um.InnerJoin("organization").As("org").OnEQ(psql.Quote("file", "organization_id"), psql.Quote("org", "id")),
		um.Where(psql.Quote("pool", "id").EQ(psql.Arg(flyover.ID))),
	).Exec(ctx, txn)
	if err != nil {
		return nil, fmt.Errorf("failed to update flyover geometry: %w", err)
	}
	return flyover, nil
}
func insertPoollistRow(ctx context.Context, txn bob.Tx, file *models.FileuploadFile, c *models.FileuploadCSV, line_number int32, header_types []headerFlyoverEnum, header_names []string, row []string) (*models.FileuploadPool, error) {
	tags := make(map[string]string, 0)
	// Start with a setter with default values, comment out the required fields to ensure they're set
	setter := models.FileuploadPoolSetter{
		// AddressCity: omit.From(),
		// AddressPostalCode: omit.From(),
		// AddressStreet: omit.From(),
		Committed: omit.From(false),
		Condition: omit.From(enums.PoolconditiontypeUnknown),
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
		PropertyOwnerName:      omit.From(""),
		PropertyOwnerPhoneE164: omitnull.FromPtr[string](nil),
		ResidentOwned:          omitnull.FromPtr[bool](nil),
		ResidentPhoneE164:      omitnull.FromPtr[string](nil),
		// Can't set this via a Setter
		// Tags:       		convertToPGData(tags),
	}
	for i, value := range row {
		if value == "" {
			continue
		}
		header_type := header_types[i]
		switch header_type {
		case headerFlyoverAddressLocality:
			setter.AddressLocality = omit.From(value)
		case headerFlyoverAddressPostalCode:
			setter.AddressPostalCode = omit.From(value)
		case headerFlyoverAddressStreet:
			setter.AddressStreet = omit.From(value)
		case headerFlyoverComment:
			condition, err := parsePoolCondition(value)
			if err == nil {
				setter.Condition = omit.From(condition)
			} else {
				addError(ctx, txn, c, int32(line_number), int32(i), fmt.Sprintf("'%s' is not a pool condition that we recognize. It should be one of %s", value, poolConditionValidValues()))
				continue
			}
		}

	}
	setter.Tags = omit.From(db.ConvertToPGData(tags))
	return models.FileuploadPools.Insert(&setter).One(ctx, txn)
}

type parseHeaderFunc[EnumType any] = func(row []string) ([]EnumType, []string)

func makeParseHeaders[EnumType any](headerToType map[string]EnumType) parseHeaderFunc[EnumType] {
	return func(row []string) ([]EnumType, []string) {
		result_enums := make([]EnumType, len(row))
		result_names := make([]string, len(row))
		for i, h := range row {
			ht := strings.TrimSpace(h)
			hl := strings.ToLower(ht)
			log.Debug().Str("header", hl).Msg("Saw CSV header")
			var type_ EnumType
			type_, ok := headerToType[hl]
			if !ok {
				// See if there is a '*' entry which should match anything
				all_type, ok2 := headerToType["*"]
				if !ok2 {
					log.Error().Str("name", hl).Msg("No header type matches column. You should add a '*' to the makeParseHeaders call")
					continue
				} else {
					type_ = all_type
				}
			}
			result_enums[i] = type_
			result_names[i] = hl
		}

		return result_enums, result_names
	}
}

func processCSVFlyover(ctx context.Context, txn bob.Tx, file *models.FileuploadFile, c *models.FileuploadCSV, rows []*models.FileuploadPool) error {
	return nil
}

var poolConditionAliases = map[string]string{
	"covered":       "unknown",
	"dark bottom":   "unknown",
	"no data":       "unknown",
	"empty":         "dry",
	"green":         "green",
	"murky pool":    "murky",
	"putting green": "false pool",
	"questionable":  "unknown",
}

func parsePoolCondition(c string) (enums.Poolconditiontype, error) {
	var condition enums.Poolconditiontype
	col_l := strings.ToLower(c)
	col_translated, ok := poolConditionAliases[col_l]
	if ok {
		col_l = col_translated
	}
	err := condition.Scan(col_l)
	return condition, err
}
