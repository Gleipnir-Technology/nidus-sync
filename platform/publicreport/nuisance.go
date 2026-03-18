package publicreport

import (
	"context"
	"fmt"

	"github.com/Gleipnir-Technology/bob"
	"github.com/Gleipnir-Technology/bob/dialect/psql"
	"github.com/Gleipnir-Technology/bob/dialect/psql/sm"
	//"github.com/Gleipnir-Technology/nidus-sync/config"
	"github.com/Gleipnir-Technology/nidus-sync/db"
	//"github.com/google/uuid"
	//"github.com/rs/zerolog/log"
	"github.com/stephenafamo/scan"
)

type Nuisance struct {
	AdditionalInfo      string `db:"additional_info" json:"additional_info"`
	Duration            string `db:"duration" json:"duration"`
	IsLocationBackyard  bool   `db:"is_location_backyard" json:"is_location_backyard"`
	IsLocationFrontyard bool   `db:"is_location_frontyard" json:"is_location_frontyard"`
	IsLocationGarden    bool   `db:"is_location_garden" json:"is_location_garden"`
	IsLocationOther     bool   `db:"is_location_other" json:"is_location_other"`
	IsLocationPool      bool   `db:"is_location_pool" json:"is_location_pool"`
	ReportID            int32  `db:"report_id" json:"-"`
	SourceContainer     bool   `db:"source_container" json:"source_container"`
	SourceDescription   string `db:"source_description" json:"source_description"`
	SourceGutter        bool   `db:"source_gutter" json:"source_gutter"`
	SourceStagnant      bool   `db:"source_stagnant" json:"source_stagnant"`
	TODDay              bool   `db:"tod_day" json:"time_of_day_day"`
	TODEarly            bool   `db:"tod_early" json:"time_of_day_early"`
	TODEvening          bool   `db:"tod_evening" json:"time_of_day_evening"`
	TODNight            bool   `db:"tod_night" json:"time_of_day_night"`
}

func nuisancesByReportID(ctx context.Context, report_ids []int32) (map[int32]*Nuisance, error) {
	rows, err := bob.All(ctx, db.PGInstance.BobDB, psql.Select(
		sm.Columns(
			"additional_info",
			"duration",
			"is_location_backyard",
			"is_location_frontyard",
			"is_location_garden",
			"is_location_other",
			"is_location_pool",
			"report_id",
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
		sm.Where(psql.Quote("report_id").EQ(
			psql.Any(report_ids),
		)),
	), scan.StructMapper[Nuisance]())
	if err != nil {
		return nil, fmt.Errorf("query nuisance: %w", err)
	}
	results := make(map[int32]*Nuisance, len(rows))
	for _, row := range rows {
		results[row.ReportID] = &Nuisance{
			AdditionalInfo:      row.AdditionalInfo,
			Duration:            row.Duration,
			IsLocationBackyard:  row.IsLocationBackyard,
			IsLocationFrontyard: row.IsLocationFrontyard,
			IsLocationGarden:    row.IsLocationGarden,
			IsLocationOther:     row.IsLocationOther,
			IsLocationPool:      row.IsLocationPool,
			SourceContainer:     row.SourceContainer,
			SourceDescription:   row.SourceDescription,
			SourceGutter:        row.SourceGutter,
			SourceStagnant:      row.SourceStagnant,
			TODDay:              row.TODDay,
			TODEarly:            row.TODEarly,
			TODEvening:          row.TODEvening,
			TODNight:            row.TODNight,
		}
	}
	return results, nil
}
