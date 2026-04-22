package types

import (
	"time"
)

type PublicReport struct {
	Address    Address                           `db:"address" json:"address"`
	Concerns   []*ConcernComplianceReportRequest `db:"-" json:"concerns"`
	Created    time.Time                         `db:"created" json:"created"`
	ID         int32                             `db:"id" json:"-"`
	Images     []Image                           `db:"images" json:"images"`
	Location   *Location                         `db:"location" json:"location"`
	Log        []LogEntry                        `db:"-" json:"log"`
	DistrictID *int32                            `db:"organization_id" json:"-"`
	District   *string                           `db:"-" json:"district"`
	PublicID   string                            `db:"public_id" json:"public_id"`
	Reporter   Contact                           `db:"reporter" json:"reporter"`
	Status     string                            `db:"status" json:"status"`
	Type       string                            `db:"report_type" json:"type"`
	URI        string                            `db:"-" json:"uri"`
}
type PublicReportCompliance struct {
	PublicReport

	AccessInstructions string `db:"access_instructions" json:"access_instructions"`
	AvailabilityNotes  string `db:"availability_notes" json:"availability_notes"`
	Comments           string `db:"comments" json:"comments"`
	GateCode           string `db:"gate_code" json:"gate_code"`
	HasDog             *bool  `db:"has_dog" json:"has_dog"`
	PermissionType     string `db:"permission_type" json:"permission_type"`
	ReportID           int32  `db:"report_id" json:"-"`
	ReportPhoneCanText *bool  `db:"report_phone_can_text" json:"can_text"`
	WantsScheduled     *bool  `db:"wants_scheduled" json:"wants_scheduled"`
}
type PublicReportNuisance struct {
	PublicReport

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
type PublicReportWater struct {
	PublicReport
	AccessComments         string  `db:"access_comments" json:"access_comments"`
	AccessGate             bool    `db:"access_gate" json:"access_gate"`
	AccessFence            bool    `db:"access_fence" json:"access_fence"`
	AccessLocked           bool    `db:"access_locked" json:"access_locked"`
	AccessDog              bool    `db:"access_dog" json:"access_dog"`
	AccessOther            bool    `db:"access_other" json:"access_other"`
	Comments               string  `db:"comments" json:"comments"`
	HasAdult               bool    `db:"has_adult" json:"has_adult"`
	HasBackyardPermission  bool    `db:"has_backyard_permission" json:"has_backyard_permission"`
	HasLarvae              bool    `db:"has_larvae" json:"has_larvae"`
	HasPupae               bool    `db:"has_pupae" json:"has_pupae"`
	IsReporterConfidential bool    `db:"is_reporter_confidential" json:"is_reporter_confidential"`
	IsReporterOwner        bool    `db:"is_reporter_owner" json:"is_reporter_owner"`
	Owner                  Contact `db:"owner" json:"owner"`
	ReportID               int32   `db:"report_id" json:"-"`
}
