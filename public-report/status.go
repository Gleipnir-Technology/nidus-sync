package publicreport

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/Gleipnir-Technology/nidus-sync/config"
	"github.com/Gleipnir-Technology/nidus-sync/db"
	"github.com/Gleipnir-Technology/nidus-sync/db/models"
	"github.com/Gleipnir-Technology/nidus-sync/db/sql"
	"github.com/Gleipnir-Technology/nidus-sync/htmlpage"
	"github.com/go-chi/chi/v5"
	"github.com/rs/zerolog/log"
	"github.com/stephenafamo/bob"
	"github.com/stephenafamo/bob/dialect/psql"
	"github.com/stephenafamo/bob/dialect/psql/sm"
	"github.com/stephenafamo/scan"
	/*
		"strconv"

		"github.com/Gleipnir-Technology/nidus-sync/db"
		"github.com/Gleipnir-Technology/nidus-sync/h3utils"
		"github.com/aarondl/opt/omit"
		"github.com/aarondl/opt/omitnull"
	*/)

type Contact struct {
	Email string
	Name  string
	Phone string
}
type Image struct {
	URL string
}
type Report struct {
	Address   string
	Comments  string
	Created   time.Time
	ID        string
	Images    []Image
	Location  string // GeoJSON
	Reporter  Contact
	SiteOwner Contact
	Type      string
}

type ContentStatus struct {
	Error    string
	ReportID string
}
type ContentStatusByID struct {
	MapboxToken string
	Report      Report
}

var (
	Status     = buildTemplate("status", "base")
	StatusByID = buildTemplate("status-by-id", "base")
)

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
	if report_id_str == "" {
		htmlpage.RenderOrError(
			w,
			Status,
			ContentStatus{
				Error:    "",
				ReportID: "",
			},
		)
		return
	}
	report_id := sanitizeReportID(report_id_str)
	report_id_str = formatReportID(report_id)
	results, err := sql.PublicreportIDTable(report_id).All(r.Context(), db.PGInstance.BobDB)
	if err != nil {
		respondError(w, "Failed to query for report", err, http.StatusInternalServerError)
		return
	}
	if len(results) != 1 {
		log.Error().Int("count", len(results)).Str("report_id", report_id_str).Msg("Got too many results for report id. This is a programmer error.")
		htmlpage.RenderOrError(
			w,
			Status,
			ContentStatus{
				Error:    "Sorry, server's confused",
				ReportID: report_id_str,
			},
		)
	}
	result := results[0]
	if result.ExistsSomewhere {
		http.Redirect(w, r, fmt.Sprintf("/status/%s", report_id), http.StatusFound)
		return
	}
	htmlpage.RenderOrError(
		w,
		Status,
		ContentStatus{
			Error:    "Sorry, we can't find that report",
			ReportID: report_id_str,
		},
	)
}
func contentFromNuisance(ctx context.Context, report_id string) (result ContentStatusByID, err error) {
	nuisance, err := models.PublicreportNuisances.Query(
		models.SelectWhere.PublicreportNuisances.PublicID.EQ(report_id),
	).One(ctx, db.PGInstance.BobDB)
	if err != nil {
		return result, fmt.Errorf("Failed to query nuisance %s: %w", report_id, err)
	}
	result.Report.ID = report_id
	result.Report.Address = nuisance.Address
	result.Report.Created = nuisance.Created
	result.Report.Reporter.Email = nuisance.ReporterEmail
	result.Report.Reporter.Name = nuisance.ReporterName
	result.Report.Reporter.Phone = nuisance.ReporterPhone

	type LocationGeoJSON struct {
		Location string
	}
	row, err := bob.One(ctx, db.PGInstance.BobDB, psql.Select(
		sm.From(
			psql.F("ST_AsGeoJSON", "location"),
		).As("location"),
		sm.Where(psql.Quote("public_id").EQ(psql.Arg(report_id))),
	), scan.StructMapper[LocationGeoJSON]())
	if err != nil {
		return result, fmt.Errorf("Failed to query nuisance %s: %w", report_id, err)
	}
	result.Report.Location = row.Location

	return result, err
}
func contentFromPool(ctx context.Context, report_id string) (result ContentStatusByID, err error) {
	return result, err
}
func contentFromQuick(ctx context.Context, report_id string) (result ContentStatusByID, err error) {
	quick, err := models.PublicreportQuicks.Query(
		models.SelectWhere.PublicreportQuicks.PublicID.EQ(report_id),
	).One(ctx, db.PGInstance.BobDB)
	if err != nil {
		return result, fmt.Errorf("Failed to query nuisance %s: %w", report_id, err)
	}

	images, err := quick.Images().All(ctx, db.PGInstance.BobDB)
	if err != nil {
		return result, fmt.Errorf("Failed to get images %s: %w", report_id, err)
	}

	result.Report.ID = report_id
	result.Report.Address = quick.Address
	result.Report.Comments = quick.Comments
	result.Report.Created = quick.Created
	result.Report.Reporter.Email = quick.ReporterEmail
	result.Report.Reporter.Name = "-"
	result.Report.Reporter.Phone = quick.ReporterPhone
	result.Report.Type = "Quick"

	for _, image := range images {
		result.Report.Images = append(result.Report.Images, Image{
			URL: fmt.Sprintf("https://%s/image/%s", config.RMODomain, image.StorageUUID),
		})
	}
	type LocationGeoJSON struct {
		Location string
	}
	row, err := bob.One(ctx, db.PGInstance.BobDB, psql.Select(
		sm.Columns(
			psql.F("ST_AsGeoJSON", "location"),
		),
		sm.From("publicreport.quick"),
		sm.Where(psql.Quote("public_id").EQ(psql.Arg(report_id))),
	), scan.StructMapper[LocationGeoJSON]())
	if err != nil {
		return result, fmt.Errorf("Failed to query nuisance %s: %w", report_id, err)
	}
	result.Report.Location = row.Location

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
	case "quick":
		content, err = contentFromQuick(ctx, report_id)
	}
	content.MapboxToken = config.MapboxToken
	htmlpage.RenderOrError(
		w,
		StatusByID,
		content,
	)
}

/*
	func getQuick(w http.ResponseWriter, r *http.Request) {
		htmlpage.RenderOrError(
			w,
			Quick,
			ContentQuick{},
		)
	}

	func getQuickSubmitComplete(w http.ResponseWriter, r *http.Request) {
		report := r.URL.Query().Get("report")
		htmlpage.RenderOrError(
			w,
			QuickSubmitComplete,
			ContentQuickSubmitComplete{
				ReportID: report,
			},
		)
	}

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
