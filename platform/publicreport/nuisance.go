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
	AdditionalInfo      string        `db:"additional_info"`
	Address             types.Address `db:"address"`
	AddressAsGiven      string        `db:"address_as_given"`
	Created             time.Time     `db:"created"`
	Duration            string        `db:"duration"`
	ID                  int32         `db:"id"`
	Images              []types.Image
	IsLocationBackyard  bool           `db:"is_location_backyard"`
	IsLocationFrontyard bool           `db:"is_location_frontyard"`
	IsLocationGarden    bool           `db:"is_location_garden"`
	IsLocationOther     bool           `db:"is_location_other"`
	IsLocationPool      bool           `db:"is_location_pool"`
	Location            types.Location `db:"location"`
	PublicID            string         `db:"public_id"`
	Reporter            Reporter       `db:"reporter"`
	SourceContainer     bool           `db:"source_container"`
	SourceDescription   string         `db:"source_description"`
	SourceGutter        bool           `db:"source_gutter"`
	SourceStagnant      bool           `db:"source_stagnant"`
	TODDay              bool           `db:"tod_day"`
	TODEarly            bool           `db:"tod_early"`
	TODEvening          bool           `db:"tod_evening"`
	TODNight            bool           `db:"tod_night"`
}
type Reporter struct {
	Email *string `db:"reporter_email"`
	Name  *string `db:"reporter_name"`
	Phone *string `db:"reporter_phone"`
}

func NuisanceReportForOrganization(ctx context.Context, org_id int32) ([]Nuisance, error) {
	reports, err := bob.All(ctx, db.PGInstance.BobDB, psql.Select(
		sm.Columns(
			"additional_info",
			"address AS address_as_given",
			"address_country AS \"address.country\"",
			"address_number AS \"address.number\"",
			"address_place AS \"address.place\"",
			"address_postcode AS \"address.postcode\"",
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
	), scan.StructMapper[Nuisance]())
	if err != nil {
		return nil, fmt.Errorf("get reports: %w", err)
	}
	report_ids := make([]int32, len(reports))
	for _, report := range reports {
		report_ids = append(report_ids, report.ID)
	}
	images_by_id, err := loadImagesForReportNuisance(ctx, org_id, report_ids)
	if err != nil {
		return nil, fmt.Errorf("images for report: %w", err)
	}
	for _, report := range reports {
		report.Images = images_by_id[report.ID]
	}
	return reports, nil
}
