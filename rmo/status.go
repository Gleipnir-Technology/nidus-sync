package rmo

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/Gleipnir-Technology/bob"
	"github.com/Gleipnir-Technology/bob/dialect/psql"
	"github.com/Gleipnir-Technology/bob/dialect/psql/sm"
	"github.com/Gleipnir-Technology/nidus-sync/db"
	"github.com/Gleipnir-Technology/nidus-sync/db/models"
	"github.com/Gleipnir-Technology/nidus-sync/db/sql"
	"github.com/Gleipnir-Technology/nidus-sync/html"
	"github.com/go-chi/chi/v5"
	//"github.com/rs/zerolog/log"
	"github.com/stephenafamo/scan"
	/*
		"github.com/Gleipnir-Technology/nidus-sync/db"
		"github.com/Gleipnir-Technology/nidus-sync/h3utils"
		"github.com/aarondl/opt/omit"
		"github.com/aarondl/opt/omitnull"
	*/)

type ContentStatus struct {
	District *ContentDistrict
	Error    string
	ReportID string
	URL      ContentURL
}
type ContentStatusByID struct {
	District *ContentDistrict
	Report   Report
	Timeline []TimelineEntry
	URL      ContentURL
}
type DetailEntry struct {
	Name  string
	Value string
}
type Report struct {
	Address    string
	Comments   string
	Created    time.Time
	Details    []DetailEntry
	ID         string
	ImageCount int
	Location   string // GeoJSON
	Status     string
	Type       string
}
type TimelineEntry struct {
	At     time.Time
	Detail string
	Title  string
}

func formatReportID(s string) string {
	// truncate down if too long
	if len(s) > 12 {
		s = s[:12]
	}

	// If less than 4 characters, return as is
	if len(s) < 4 {
		return s
	}

	// If at least 8 characters, add hyphens at positions 4 and 8
	if len(s) >= 8 {
		return s[0:4] + "-" + s[4:8] + "-" + s[8:]
	}

	// If at least 4 characters but less than 8, add hyphen only at position 4
	return s[0:4] + "-" + s[4:]
}

func getStatus(w http.ResponseWriter, r *http.Request) {
	report_id_str := r.URL.Query().Get("report")
	content := ContentStatus{
		Error:    "",
		ReportID: "",
		URL:      makeContentURL(nil),
	}
	if report_id_str == "" {
		html.RenderOrError(w, "rmo/status.html", content)
		return
	}
	report_id := sanitizeReportID(report_id_str)
	report_id_str = formatReportID(report_id)
	//some_report, e := report.FindSomeReport(r.Context(), report_id)
	content.Error = "Sorry, we can't find that report"
	html.RenderOrError(w, "rmo/status.html", content)
}
func contentFromNuisance(ctx context.Context, report_id string) (result ContentStatusByID, err error) {
	nuisance, err := models.PublicreportNuisances.Query(
		models.SelectWhere.PublicreportNuisances.PublicID.EQ(report_id),
	).One(ctx, db.PGInstance.BobDB)
	if err != nil {
		return result, fmt.Errorf("Failed to query nuisance %s: %w", report_id, err)
	}

	images, err := sql.PublicreportImageWithJSONByNuisanceID(nuisance.ID).All(ctx, db.PGInstance.BobDB)
	if err != nil {
		return result, fmt.Errorf("Failed to get images %s: %w", report_id, err)
	}

	if !nuisance.OrganizationID.IsNull() {
		org_id := nuisance.OrganizationID.MustGet()
		org, err := models.FindOrganization(ctx, db.PGInstance.BobDB, org_id)
		if err != nil {
			return result, fmt.Errorf("Failed to get district %d information: %w", org_id, err)
		}
		result.District = newContentDistrict(org)
	}
	result.Report.ID = report_id
	result.Report.Address = nuisance.AddressRaw
	result.Report.Created = nuisance.Created
	result.Report.ImageCount = len(images)
	result.Report.Status = strings.Title(nuisance.Status.String())
	result.Report.Type = "Mosquito Nuisance"
	result.Report.Details = []DetailEntry{
		DetailEntry{
			Name:  "Active early morning (5a-8a)?",
			Value: strconv.FormatBool(nuisance.TodEarly),
		},
		DetailEntry{
			Name:  "Active daytime (8a-5p)?",
			Value: strconv.FormatBool(nuisance.TodDay),
		},
		DetailEntry{
			Name:  "Active evening (5p-9p)?",
			Value: strconv.FormatBool(nuisance.TodEvening),
		},
		DetailEntry{
			Name:  "Active night (9p-5a)?",
			Value: strconv.FormatBool(nuisance.TodNight),
		},
		DetailEntry{
			Name:  "Duration",
			Value: nuisance.Duration.String(),
		},
		DetailEntry{
			Name:  "Active in backyard?",
			Value: strconv.FormatBool(nuisance.IsLocationBackyard),
		},
		DetailEntry{
			Name:  "Active in frontyard?",
			Value: strconv.FormatBool(nuisance.IsLocationFrontyard),
		},
		DetailEntry{
			Name:  "Active in garden?",
			Value: strconv.FormatBool(nuisance.IsLocationGarden),
		},
		DetailEntry{
			Name:  "Active in other location?",
			Value: strconv.FormatBool(nuisance.IsLocationOther),
		},
		DetailEntry{
			Name:  "Active in pool area?",
			Value: strconv.FormatBool(nuisance.IsLocationPool),
		},
		DetailEntry{
			Name:  "Stagnant Water",
			Value: strconv.FormatBool(nuisance.SourceStagnant),
		},
		DetailEntry{
			Name:  "Container",
			Value: strconv.FormatBool(nuisance.SourceContainer),
		},
		DetailEntry{
			Name:  "Sprinklers & Gutters",
			Value: strconv.FormatBool(nuisance.SourceGutter),
		},
	}
	result.Timeline = []TimelineEntry{
		TimelineEntry{
			At:     nuisance.Created,
			Detail: "Initial report was submitted",
			Title:  "Created",
		},
	}

	type LocationGeoJSON struct {
		Location string
	}
	location, err := bob.One(ctx, db.PGInstance.BobDB, psql.Select(
		sm.Columns(
			psql.F("ST_AsGeoJSON", "location"),
		),
		sm.From("publicreport.nuisance"),
		sm.Where(psql.Quote("public_id").EQ(psql.Arg(report_id))),
	), scan.SingleColumnMapper[string])
	if err != nil {
		return result, fmt.Errorf("Failed to query nuisance %s: %w", report_id, err)
	}
	result.Report.Location = location

	return result, err
}
func contentFromWater(ctx context.Context, report_id string) (result ContentStatusByID, err error) {
	water, err := models.PublicreportWaters.Query(
		models.SelectWhere.PublicreportWaters.PublicID.EQ(report_id),
	).One(ctx, db.PGInstance.BobDB)
	if err != nil {
		return result, fmt.Errorf("Failed to query water %s: %w", report_id, err)
	}

	images, err := sql.PublicreportImageWithJSONByWaterID(water.ID).All(ctx, db.PGInstance.BobDB)
	if err != nil {
		return result, fmt.Errorf("Failed to get images %s: %w", report_id, err)
	}

	if !water.OrganizationID.IsNull() {
		org_id := water.OrganizationID.MustGet()
		org, err := models.FindOrganization(ctx, db.PGInstance.BobDB, org_id)
		if err != nil {
			return result, fmt.Errorf("Failed to get district %d information: %w", org_id, err)
		}
		result.District = newContentDistrict(org)
	}

	result.Report.ID = report_id
	result.Report.Address = water.AddressRaw
	result.Report.Created = water.Created
	result.Report.ImageCount = len(images)
	result.Report.Status = strings.Title(water.Status.String())
	result.Report.Type = "Mosquito Nuisance"
	result.Report.Details = []DetailEntry{
		DetailEntry{
			Name:  "Has a gate that affects access?",
			Value: strconv.FormatBool(water.AccessGate),
		},
		DetailEntry{
			Name:  "Has dog that affects access?",
			Value: strconv.FormatBool(water.AccessDog),
		},
		DetailEntry{
			Name:  "Has a fence that affects access?",
			Value: strconv.FormatBool(water.AccessFence),
		},
		DetailEntry{
			Name:  "Has a locked entrace that affects access?",
			Value: strconv.FormatBool(water.AccessLocked),
		},
		DetailEntry{
			Name:  "Reporter observed larvae (wigglers)?",
			Value: strconv.FormatBool(water.HasLarvae),
		},
		DetailEntry{
			Name:  "Reporter observed pupae (tumblers)?",
			Value: strconv.FormatBool(water.HasPupae),
		},
		DetailEntry{
			Name:  "Reporter observed adult mosquitoes?",
			Value: strconv.FormatBool(water.HasAdult),
		},
	}
	result.Timeline = []TimelineEntry{
		TimelineEntry{
			At:     water.Created,
			Detail: "Initial report was submitted",
			Title:  "Created",
		},
	}
	location, err := bob.One(ctx, db.PGInstance.BobDB, psql.Select(
		sm.Columns(
			psql.F("ST_AsGeoJSON", "location"),
		),
		sm.From("publicreport.water"),
		sm.Where(psql.Quote("public_id").EQ(psql.Arg(report_id))),
	), scan.SingleColumnMapper[string])
	if err != nil {
		return result, fmt.Errorf("Failed to query water %s: %w", report_id, err)
	}
	result.Report.Location = location

	return result, err
}

func getStatusByID(w http.ResponseWriter, r *http.Request) {
	report_id := chi.URLParam(r, "report_id")
	ctx := r.Context()

	location, err := models.PublicreportReportLocations.Query(
		models.SelectWhere.PublicreportReportLocations.PublicID.EQ(report_id),
	).One(ctx, db.PGInstance.BobDB)
	if err != nil {
		respondError(w, "Failed to find report", err, http.StatusBadRequest)
		return
	}
	var content ContentStatusByID
	switch location.TableName.MustGet() {
	case "nuisance":
		content, err = contentFromNuisance(ctx, report_id)
	case "water":
		content, err = contentFromWater(ctx, report_id)
	default:
		err = fmt.Errorf("table name %s not in switch", location.TableName.MustGet())
	}
	if err != nil {
		respondError(w, "Failed to generate report content", err, http.StatusInternalServerError)
		return
	}
	content.URL = makeContentURL(nil)
	html.RenderOrError(
		w,
		"rmo/status-by-id.html",
		content,
	)
}

func sanitizeReportID(r string) string {
	result := ""
	for _, char := range r {
		if char != '-' {
			result += string(char)
		}
	}
	return strings.ToUpper(result)
}
