package types

import (
	"time"
)

type LogEntry struct {
	Created  time.Time `db:"created" json:"created"`
	ID       int32     `db:"id" json:"-"`
	Message  string    `db:"message" json:"message"`
	ReportID int32     `db:"report_id" json:"-"`
	Type     string    `db:"type_" json:"type"`
	UserID   *int32    `db:"user_id" json:"user_id"`
}
