package publicreport

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/Gleipnir-Technology/nidus-sync/db"
	"github.com/Gleipnir-Technology/nidus-sync/db/sql"
	"github.com/Gleipnir-Technology/nidus-sync/htmlpage"
	"github.com/go-chi/chi/v5"
	"github.com/rs/zerolog/log"
	/*
			"strconv"
			"time"

		"github.com/Gleipnir-Technology/nidus-sync/db/models"
			"github.com/Gleipnir-Technology/nidus-sync/db"
			"github.com/Gleipnir-Technology/nidus-sync/h3utils"
			"github.com/aarondl/opt/omit"
			"github.com/aarondl/opt/omitnull"
			"github.com/stephenafamo/bob/dialect/psql"
			"github.com/stephenafamo/bob/dialect/psql/um"
	*/)

type Report struct {
	ID string
}

type ContextStatus struct {
	Error    string
	ReportID string
}
type ContextStatusByID struct {
	Report Report
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
			ContextStatus{
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
			ContextStatus{
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
		ContextStatus{
			Error:    "Sorry, we can't find that report",
			ReportID: report_id_str,
		},
	)
}
func getStatusByID(w http.ResponseWriter, r *http.Request) {
	report_id := chi.URLParam(r, "report_id")
	htmlpage.RenderOrError(
		w,
		StatusByID,
		ContextStatusByID{
			Report: Report{
				ID: report_id,
			},
		},
	)
}

/*
	func getQuick(w http.ResponseWriter, r *http.Request) {
		htmlpage.RenderOrError(
			w,
			Quick,
			ContextQuick{},
		)
	}

	func getQuickSubmitComplete(w http.ResponseWriter, r *http.Request) {
		report := r.URL.Query().Get("report")
		htmlpage.RenderOrError(
			w,
			QuickSubmitComplete,
			ContextQuickSubmitComplete{
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
