package main

import (
	"time"

	"github.com/aarondl/opt/null"
)

func fsTimestampToTime(t null.Val[int64]) *time.Time {
	if t.IsNull() {
		return nil
	}
	result := time.UnixMilli(t.MustGet())
	return &result
}
