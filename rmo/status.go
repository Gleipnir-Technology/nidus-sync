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
	"github.com/Gleipnir-Technology/nidus-sync/config"
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
	District    *ContentDistrict
	Error       string
	MapboxToken string
	ReportID    string
	URL         ContentURL
}
type ContentStatusByID struct {
	District    *ContentDistrict
	MapboxToken string
	Report      Report
	Timeline    []TimelineEntry
	URL         ContentURL
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
		Error:       "",
		MapboxToken: config.MapboxToken,
		ReportID:    "",
		URL:         makeContentURL(nil),
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
	result.Report.Address = nuisance.Address
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
func contentFromPool(ctx context.Context, report_id string) (result ContentStatusByID, err error) {
	pool, err := models.PublicreportPools.Query(
		models.SelectWhere.PublicreportPools.PublicID.EQ(report_id),
	).One(ctx, db.PGInstance.BobDB)
	if err != nil {
		return result, fmt.Errorf("Failed to query pool %s: %w", report_id, err)
	}

	images, err := sql.PublicreportImageWithJSONByPoolID(pool.ID).All(ctx, db.PGInstance.BobDB)
	if err != nil {
		return result, fmt.Errorf("Failed to get images %s: %w", report_id, err)
	}

	if !pool.OrganizationID.IsNull() {
		org_id := pool.OrganizationID.MustGet()
		org, err := models.FindOrganization(ctx, db.PGInstance.BobDB, org_id)
		if err != nil {
			return result, fmt.Errorf("Failed to get district %d information: %w", org_id, err)
		}
		result.District = newContentDistrict(org)
	}

	result.Report.ID = report_id
	result.Report.Address = pool.Address
	result.Report.Created = pool.Created
	result.Report.ImageCount = len(images)
	result.Report.Status = strings.Title(pool.Status.String())
	result.Report.Type = "Mosquito Nuisance"
	result.Report.Details = []DetailEntry{
		DetailEntry{
			Name:  "Has a gate that affects access?",
			Value: strconv.FormatBool(pool.AccessGate),
		},
		DetailEntry{
			Name:  "Has dog that affects access?",
			Value: strconv.FormatBool(pool.AccessDog),
		},
		DetailEntry{
			Name:  "Has a fence that affects access?",
			Value: strconv.FormatBool(pool.AccessFence),
		},
		DetailEntry{
			Name:  "Has a locked entrace that affects access?",
			Value: strconv.FormatBool(pool.AccessLocked),
		},
		DetailEntry{
			Name:  "Reporter observed larvae (wigglers)?",
			Value: strconv.FormatBool(pool.HasLarvae),
		},
		DetailEntry{
			Name:  "Reporter observed pupae (tumblers)?",
			Value: strconv.FormatBool(pool.HasPupae),
		},
		DetailEntry{
			Name:  "Reporter observed adult mosquitoes?",
			Value: strconv.FormatBool(pool.HasAdult),
		},
	}
	result.Timeline = []TimelineEntry{
		TimelineEntry{
			At:     pool.Created,
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
		sm.From("publicreport.pool"),
		sm.Where(psql.Quote("public_id").EQ(psql.Arg(report_id))),
	), scan.SingleColumnMapper[string])
	if err != nil {
		return result, fmt.Errorf("Failed to query pool %s: %w", report_id, err)
	}
	result.Report.Location = location

	return result, err
}

func contentFromQuick(ctx context.Context, report_id string) (result ContentStatusByID, err error) {
	quick, err := models.PublicreportQuicks.Query(
		models.SelectWhere.PublicreportQuicks.PublicID.EQ(report_id),
	).One(ctx, db.PGInstance.BobDB)
	if err != nil {
		return result, fmt.Errorf("Failed to query nuisance %s: %w", report_id, err)
	}

	result.Report.ID = report_id
	result.Report.Address = quick.Address
	result.Report.Comments = quick.Comments
	result.Report.Created = quick.Created
	result.Report.Type = "Quick"

	type LocationGeoJSON struct {
		Location string
	}
	location, err := bob.One(ctx, db.PGInstance.BobDB, psql.Select(
		sm.Columns(
			psql.F("ST_AsGeoJSON", "location"),
		),
		sm.From("publicreport.quick"),
		sm.Where(psql.Quote("public_id").EQ(psql.Arg(report_id))),
	), scan.SingleColumnMapper[string])
	if err != nil {
		return result, fmt.Errorf("Failed to query nuisance %s: %w", report_id, err)
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
	case "pool":
		content, err = contentFromPool(ctx, report_id)
	}
	if err != nil {
		respondError(w, "Failed to generate report content", err, http.StatusInternalServerError)
		return
	}
	content.MapboxToken = config.MapboxToken
	content.URL = makeContentURL(nil)
	html.RenderOrError(
		w,
		"rmo/status-by-id.html",
		content,
	)
}

/*
	func postQuick(w http.ResponseWriter, r *http.Request) {
		err := r.ParseMultipartForm(32 << 10) // 32 MB buffer
		if err != nil {
			respondError(w, "Failed to parse form", err, http.StatusBadRequest)
			return
		}
		lat := r.FormValue("latitude")
		lng := r.FormValue("longitude")
		comments := r.FormValue("comments")
		//photos := r.FormValue("photos")

		latitude, err := strconv.ParseFloat(lat, 64)
		if err != nil {
			respondError(w, "Failed to create parse latitude", err, http.StatusBadRequest)
			return
		}
		longitude, err := strconv.ParseFloat(lng, 64)
		if err != nil {
			respondError(w, "Failed to create parse longitude", err, http.StatusBadRequest)
			return
		}
		u, err := GenerateReportID()
		if err != nil {
			respondError(w, "Failed to create quick report public ID", err, http.StatusInternalServerError)
			return
		}
		c, err := h3utils.GetCell(longitude, latitude, 15)
		setter := models.PublicreportQuickSetter{
			Created:  omit.From(time.Now()),
			Comments: omit.From(comments),
			//Location: omitnull.From(fmt.Sprintf("ST_GeometryFromText(Point(%s %s))", longitude, latitude)),
			H3cell:		omitnull.From(c.String()),
			PublicID:	  omit.From(u),
			ReporterEmail: omit.From(""),
			ReporterPhone: omit.From(""),
		}
		quick, err := models.PublicreportQuicks.Insert(&setter).One(r.Context(), db.PGInstance.BobDB)
		if err != nil {
			respondError(w, "Failed to create database record", err, http.StatusInternalServerError)
			return
		}
		_, err = psql.Update(
			um.Table("publicreport.quick"),
			um.SetCol("location").To(fmt.Sprintf("ST_GeometryFromText('Point(%f %f)')", longitude, latitude)),
			um.Where(psql.Quote("id").EQ(psql.Arg(quick.ID))),
		).Exec(r.Context(), db.PGInstance.BobDB)
		if err != nil {
			respondError(w, "Failed to insert publicreport", err, http.StatusInternalServerError)
			return
		}
		log.Info().Float64("latitude", latitude).Float64("longitude", longitude).Msg("Got upload")
		photoSetters := make([]*models.PublicreportQuickPhotoSetter, 0)
		uploads, err := extractPhotoUploads(r)
		if err != nil {
			respondError(w, "Failed to extract photo uploads", err, http.StatusInternalServerError)
			return
		}
		for _, u := range uploads {
			photoSetters = append(photoSetters, &models.PublicreportQuickPhotoSetter{
				Filename: omit.From(u.Filename),
				Size:	 omit.From(u.Size),
				UUID:	 omit.From(u.UUID),
			})
		}
		err = quick.InsertQuickPhotos(r.Context(), db.PGInstance.BobDB, photoSetters...)
		if err != nil {
			respondError(w, "Failed to create photo records", err, http.StatusInternalServerError)
			return
		}
		http.Redirect(w, r, fmt.Sprintf("/quick-submit-complete?report=%s", u), http.StatusFound)
	}
*/
func sanitizeReportID(r string) string {
	result := ""
	for _, char := range r {
		if char != '-' {
			result += string(char)
		}
	}
	return strings.ToUpper(result)
}
