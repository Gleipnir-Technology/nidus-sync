package types

import (
	"time"
)

type ReviewTask struct {
	Created  time.Time       `db:"created" json:"created"`
	Creator  User            `db:"creator" json:"creator"`
	ID       int32           `db:"id" json:"id"`
	Pool     *ReviewTaskPool `db:"pool" json:"pool"`
	Reviewed *time.Time      `db:"reviewed" json:"reviewed"`
	Reviewer *User           `db:"reviewer" json:"reviewer"`
}
type ReviewTaskPool struct {
	Condition string   `db:"condition" json:"condition"`
	Location  Location `db:"location" json:"location"`
	Site      Site     `db:"site" json:"site"`
}
