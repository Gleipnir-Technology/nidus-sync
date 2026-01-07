package htmlpage

import (
	"bytes"
	"embed"
	"errors"
	"fmt"
	"html/template"
	"io"
	"math"
	"net/http"
	"os"
	"path"
	"strconv"
	"strings"
	"time"

	"github.com/Gleipnir-Technology/nidus-sync/config"
	"github.com/aarondl/opt/null"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
)

var TemplatesByFilename = make(map[string]BuiltTemplate, 0)

type BuiltTemplate struct {
	files  []string
	subdir string
	// Nil if we are going to read templates off disk every time we render
	// because we are in development mode.
	template *template.Template
}

func (bt *BuiltTemplate) executeTemplate(w io.Writer, data any) error {
	if bt.template == nil {
		name := path.Base(bt.files[0])
		templ, err := parseFromDisk(bt.files)
		if err != nil {
			return fmt.Errorf("Failed to parse template file: %w", err)
		}
		if templ == nil {
			w.Write([]byte("Failed to read from disk: "))
			return errors.New("Template parsing failed")
		}
		log.Debug().Str("name", templ.Name()).Msg("Parsed template")
		return templ.ExecuteTemplate(w, name, data)
	} else {
		name := path.Base(bt.files[0])
		return bt.template.ExecuteTemplate(w, name, data)
	}
}

func NewBuiltTemplate(embeddedFiles embed.FS, subdir string, files ...string) *BuiltTemplate {
	files_on_disk := true
	for _, f := range files {
		_, err := os.Stat(f)
		if err != nil {
			files_on_disk = false
			if !config.IsProductionEnvironment() {
				log.Warn().Str("file", f).Msg("template file is not on disk")
			}
			break
		}
	}
	var result BuiltTemplate
	if files_on_disk {
		result = BuiltTemplate{
			files:    files,
			subdir:   subdir,
			template: nil,
		}
	} else {
		result = BuiltTemplate{
			files:    files,
			subdir:   subdir,
			template: parseEmbedded(embeddedFiles, subdir, files),
		}
	}
	TemplatesByFilename[files[0]] = result
	return &result
}

func RenderOrError(w http.ResponseWriter, template *BuiltTemplate, context interface{}) {
	buf := &bytes.Buffer{}
	err := template.executeTemplate(buf, context)
	if err != nil {
		log.Error().Err(err).Strs("files", template.files).Msg("Failed to render template")
		respondError(w, "Failed to render template", err, http.StatusInternalServerError)
		return
	}
	buf.WriteTo(w)
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

func makeFuncMap() template.FuncMap {
	funcMap := template.FuncMap{
		"bigNumber":          bigNumber,
		"GISStatement":       gisStatement,
		"latLngDisplay":      latLngDisplay,
		"timeAsRelativeDate": timeAsRelativeDate,
		"timeDelta":          timeDelta,
		"timeElapsed":        timeElapsed,
		"timeInterval":       timeInterval,
		"timeSince":          timeSince,
		"uuidShort":          uuidShort,
	}
	return funcMap
}
func parseEmbedded(embeddedFiles embed.FS, subdir string, files []string) *template.Template {
	funcMap := makeFuncMap()
	// Remap the file names to embedded paths
	embeddedFilePaths := make([]string, 0)
	for _, f := range files {
		embeddedFilePaths = append(embeddedFilePaths, strings.TrimPrefix(f, subdir))
	}
	name := path.Base(embeddedFilePaths[0])
	log.Debug().Str("name", name).Strs("paths", embeddedFilePaths).Msg("Parsing embedded template")
	return template.Must(
		template.New(name).Funcs(funcMap).ParseFS(embeddedFiles, embeddedFilePaths...))
}

func parseFromDisk(files []string) (*template.Template, error) {
	funcMap := makeFuncMap()
	name := path.Base(files[0])
	log.Debug().Str("name", name).Strs("files", files).Msg("parsing from disk")
	templ, err := template.New(name).Funcs(funcMap).ParseFiles(files...)
	if err != nil {
		return nil, fmt.Errorf("Failed to parse %s: %w", files, err)
	}
	return templ, nil
}

func timeAsRelativeDate(d time.Time) string {
	return d.Format("01-02")
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
func timeSince(t *time.Time) string {
	if t == nil {
		return "never"
	}
	now := time.Now()
	diff := now.Sub(*t)

	hours := diff.Hours()
	if hours < 1 {
		minutes := diff.Minutes()
		return fmt.Sprintf("%d minutes ago", int(minutes))
	} else if hours < 24 {
		return fmt.Sprintf("%d hours ago", int(hours))
	} else {
		days := hours / 24
		return fmt.Sprintf("%d days ago", int(days))
	}
}
func uuidShort(uuid uuid.UUID) string {
	s := uuid.String()
	if len(s) < 7 {
		return s // Return as is if too short
	}

	return s[:3] + "..." + s[len(s)-4:]
}
