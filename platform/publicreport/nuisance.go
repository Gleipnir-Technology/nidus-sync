package publicreport

import (
	"context"
	"fmt"
	"time"

	"github.com/Gleipnir-Technology/bob"
	"github.com/Gleipnir-Technology/bob/dialect/psql"
	"github.com/Gleipnir-Technology/bob/dialect/psql/sm"
	//"github.com/Gleipnir-Technology/nidus-sync/config"
	"github.com/Gleipnir-Technology/nidus-sync/db"
	"github.com/Gleipnir-Technology/nidus-sync/platform/types"
	//"github.com/Gleipnir-Technology/nidus-sync/db/models"
	//"github.com/google/uuid"
	//"github.com/rs/zerolog/log"
	"github.com/stephenafamo/scan"
)

type Nuisance struct {
	AdditionalInfo      string         `db:"additional_info" json:"additional_info"`
	Address             types.Address  `db:"address" json:"address"`
	AddressRaw          string         `db:"address_raw" json:"address_raw"`
	Created             time.Time      `db:"created" json:"created"`
	Duration            string         `db:"duration" json:"duration"`
	ID                  int32          `db:"id" json:"-"`
	Images              []types.Image  `db:"-" json:"images"`
	IsLocationBackyard  bool           `db:"is_location_backyard" json:"is_location_backyard"`
	IsLocationFrontyard bool           `db:"is_location_frontyard" json:"is_location_frontyard"`
	IsLocationGarden    bool           `db:"is_location_garden" json:"is_location_garden"`
	IsLocationOther     bool           `db:"is_location_other" json:"is_location_other"`
	IsLocationPool      bool           `db:"is_location_pool" json:"is_location_pool"`
	Location            types.Location `db:"location" json:"location"`
	PublicID            string         `db:"public_id" json:"public_id"`
	Reporter            types.Contact  `db:"reporter" json:"reporter"`
	SourceContainer     bool           `db:"source_container" json:"source_container"`
	SourceDescription   string         `db:"source_description" json:"source_description"`
	SourceGutter        bool           `db:"source_gutter" json:"source_gutter"`
	SourceStagnant      bool           `db:"source_stagnant" json:"source_stagnant"`
	TODDay              bool           `db:"tod_day" json:"time_of_day_day"`
	TODEarly            bool           `db:"tod_early" json:"time_of_day_early"`
	TODEvening          bool           `db:"tod_evening" json:"time_of_day_evening"`
	TODNight            bool           `db:"tod_night" json:"time_of_day_night"`
}

func NuisanceReportForOrganization(ctx context.Context, org_id int32) ([]Nuisance, error) {
	reports, err := bob.All(ctx, db.PGInstance.BobDB, psql.Select(
		sm.Columns(
			"additional_info",
			"address_raw AS address_raw",
			"address_country AS \"address.country\"",
			"address_locality AS \"address.locality\"",
			"address_number AS \"address.number\"",
			"address_postal_code AS \"address.postal_code\"",
			"address_region AS \"address.region\"",
			"address_street AS \"address.street\"",
			"created",
			"duration",
			"id",
			"is_location_backyard",
			"is_location_frontyard",
			"is_location_garden",
			"is_location_other",
			"is_location_pool",
			"ST_Y(location::geometry::geometry(point, 4326)) AS \"location.latitude\"",
			"ST_X(location::geometry::geometry(point, 4326)) AS \"location.longitude\"",
			"public_id",
			"reporter_email AS \"reporter.email\"",
			"reporter_name AS \"reporter.name\"",
			"reporter_phone AS \"reporter.phone\"",
			"source_container",
			"source_description",
			"source_gutter",
			"source_stagnant",
			"tod_day",
			"tod_early",
			"tod_evening",
			"tod_night",
		),
		sm.From("publicreport.nuisance"),
		sm.Where(psql.Quote("publicreport", "nuisance", "organization_id").EQ(psql.Arg(org_id))),
		sm.Where(psql.Quote("publicreport", "nuisance", "reviewed").IsNull()),
	), scan.StructMapper[Nuisance]())
	if err != nil {
		return nil, fmt.Errorf("get reports: %w", err)
	}
	report_ids := make([]int32, len(reports))
	for i, report := range reports {
		report_ids[i] = report.ID
	}
	images_by_id, err := loadImagesForReportNuisance(ctx, org_id, report_ids)
	if err != nil {
		return nil, fmt.Errorf("images for report: %w", err)
	}
	for i := range reports {
		images, ok := images_by_id[reports[i].ID]
		if ok {
			reports[i].Images = images
		} else {
			reports[i].Images = []types.Image{}
		}
	}
	return reports, nil
}
func NuisanceReportForOrganizationCount(ctx context.Context, org_id int32) (uint, error) {
	type _Row struct {
		Count uint `db:"count"`
	}
	row, err := bob.One(ctx, db.PGInstance.BobDB, psql.Select(
		sm.Columns(
			"COUNT(*) AS count",
		),
		sm.From("publicreport.nuisance"),
		sm.Where(psql.Quote("publicreport", "nuisance", "organization_id").EQ(psql.Arg(org_id))),
	), scan.StructMapper[_Row]())
	if err != nil {
		return 0, fmt.Errorf("query count: %w", err)
	}
	return row.Count, nil
}
