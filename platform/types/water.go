package types

type Water struct {
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
