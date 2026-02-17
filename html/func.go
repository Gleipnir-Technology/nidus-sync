package html

import (
	//"bytes"
	"fmt"
	"html/template"
	//"io/fs"
	"math"
	//"net/http"
	//"os"
	"strconv"
	"strings"
	"time"

	"github.com/aarondl/opt/null"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
)

func addFuncMap(t *template.Template) {
	funcMap := template.FuncMap{
		"bigNumber":          bigNumber,
		"duration":           duration,
		"hasPassed":          hasPassed,
		"html":               unescapeHTML,
		"json":               unescapeJS,
		"GISStatement":       gisStatement,
		"latLngDisplay":      latLngDisplay,
		"publicReportID":     publicReportID,
		"timeAsRelativeDate": timeAsRelativeDate,
		"timeDelta":          timeDelta,
		"timeElapsed":        timeElapsed,
		"timeInterval":       timeInterval,
		"timeRelative":       timeRelative,
		"timeRelativePtr":    timeRelativePtr,
		"uuidShort":          uuidShort,
	}
	t.Funcs(funcMap)
}

func bigNumber(n int) string {
	// Convert the number to a string
	numStr := strconv.FormatInt(int64(n), 10)

	// Add commas every three digits from the right
	var result strings.Builder
	for i, char := range numStr {
		if i > 0 && (len(numStr)-i)%3 == 0 {
			result.WriteByte(',')
		}
		result.WriteRune(char)
	}

	return result.String()
}

func duration(d time.Duration) string {
	seconds := int(d.Seconds())

	if seconds < 60 {
		if seconds == 1 {
			return "1 second ago"
		}
		return fmt.Sprintf("%d seconds ago", seconds)
	}

	minutes := int(d.Minutes())
	if minutes < 60 {
		if minutes == 1 {
			return "1 minute ago"
		}
		return fmt.Sprintf("%d minutes ago", minutes)
	}

	hours := int(d.Hours())
	if hours < 24 {
		if hours == 1 {
			return "1 hour ago"
		}
		return fmt.Sprintf("%d hours ago", hours)
	}

	days := hours / 24
	if days < 30 {
		if days == 1 {
			return "1 day ago"
		}
		return fmt.Sprintf("%d days ago", days)
	}

	months := days / 30
	if months < 12 {
		if months == 1 {
			return "1 month ago"
		}
		return fmt.Sprintf("%d months ago", months)
	}

	years := days / 365
	if years == 1 {
		return "1 year ago"
	}
	return fmt.Sprintf("%d years ago", years)
}
func publicReportID(s string) string {
	if len(s) != 12 {
		return s
	}
	return s[0:4] + "-" + s[4:8] + "-" + s[8:12]
}
func timeAsRelativeDate(d time.Time) string {
	return d.Format("01-02")
}

func hasPassed(t time.Time) bool {
	return t.Before(time.Now())
}

// FormatTimeDuration returns a human-readable string representing a time.Duration
// as "X units early" or "X units late"
func timeDelta(d time.Duration) string {
	suffix := "late"
	if d < 0 {
		suffix = "early"
		d = -d // Make duration positive for calculations
	}

	const (
		day  = 24 * time.Hour
		week = 7 * day
	)

	log.Info().Int64("delta", int64(d)).Str("suffix", suffix).Msg("Time delta")
	switch {
	case d >= week:
		weeks := d / week
		if weeks == 1 {
			return "1 week " + suffix
		}
		return fmt.Sprintf("%d weeks %s", weeks, suffix)

	case d >= day:
		days := d / day
		if days == 1 {
			return "1 day " + suffix
		}
		return fmt.Sprintf("%d days %s", days, suffix)

	case d >= time.Hour:
		hours := d / time.Hour
		if hours == 1 {
			return "1 hour " + suffix
		}
		return fmt.Sprintf("%d hours %s", hours, suffix)

	case d >= time.Minute:
		minutes := d / time.Minute
		if minutes == 1 {
			return "1 minute " + suffix
		}
		return fmt.Sprintf("%d minutes %s", minutes, suffix)

	default:
		seconds := d / time.Second
		if seconds == 1 {
			return "1 second " + suffix
		}
		return fmt.Sprintf("%d seconds %s", seconds, suffix)
	}
}

func timeElapsed(seconds null.Val[float32]) string {
	if !seconds.IsValue() {
		return "none"
	}
	s := int(seconds.MustGet())
	hours := s / 3600
	remainder := s - (hours * 3600)
	minutes := remainder / 60
	remainder = remainder - (minutes * 60)
	if hours > 0 {
		return fmt.Sprintf("%02d:%02d:%02d", hours, minutes, remainder)
	} else if minutes > 0 {
		return fmt.Sprintf("%02d:%02d", minutes, remainder)
	} else {
		return fmt.Sprintf("%d seconds", remainder)
	}
}

func timeInterval(d time.Duration) string {
	seconds := d.Seconds()

	// Less than 120 seconds -> show in seconds
	if seconds < 120 {
		return fmt.Sprintf("every %d seconds", int(math.Round(seconds)))
	}

	minutes := d.Minutes()
	// Less than 120 minutes -> show in minutes
	if minutes < 120 {
		return fmt.Sprintf("every %d minutes", int(math.Round(minutes)))
	}

	hours := d.Hours()
	// Less than 48 hours -> show in hours
	if hours < 48 {
		return fmt.Sprintf("every %d hours", int(math.Round(hours)))
	}

	days := hours / 24
	// Less than 14 days -> show in days
	if days < 14 {
		return fmt.Sprintf("every %d days", int(math.Round(days)))
	}

	weeks := days / 7
	// Less than 8 weeks -> show in weeks
	if weeks < 8 {
		return fmt.Sprintf("every %d weeks", int(math.Round(weeks)))
	}

	months := days / 30
	// Less than 24 months -> show in months
	if months < 24 {
		return fmt.Sprintf("every %d months", int(math.Round(months)))
	}

	years := days / 365
	return fmt.Sprintf("every %d years", int(math.Round(years)))
}
func timeRelative(t time.Time) string {
	now := time.Now()
	diff := now.Sub(t)

	hours := diff.Hours()
	if hours > 0 {
		if hours < 1 {
			minutes := diff.Minutes()
			return fmt.Sprintf("%d minutes ago", int(minutes))
		} else if hours < 24 {
			return fmt.Sprintf("%d hours ago", int(hours))
		} else {
			days := hours / 24
			return fmt.Sprintf("%d days ago", int(days))
		}
	} else {
		if hours < -24 {
			days := hours / 24
			return fmt.Sprintf("in %d days", -1*int(days))
		} else if hours < -1 {
			return fmt.Sprintf("in %d hours", -1*int(hours))
		} else {
			minutes := diff.Minutes()
			if minutes > -1 {
				seconds := diff.Seconds()
				return fmt.Sprintf("in %d seconds", -1*int(seconds))
			}
			return fmt.Sprintf("in %d minutes", -1*int(minutes))
		}
	}
}
func timeRelativePtr(t *time.Time) string {
	if t == nil {
		return "never"
	}
	return timeRelative(*t)
}
func unescapeHTML(s string) template.HTML {
	return template.HTML(s)
}
func unescapeJS(s string) template.JS {
	return template.JS(s)
}
func uuidShort(uuid uuid.UUID) string {
	s := uuid.String()
	if len(s) < 7 {
		return s // Return as is if too short
	}

	return s[:3] + "..." + s[len(s)-4:]
}
