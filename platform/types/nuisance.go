package types

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
