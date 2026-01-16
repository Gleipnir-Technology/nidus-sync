package publicreport

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/Gleipnir-Technology/nidus-sync/db"
	"github.com/Gleipnir-Technology/nidus-sync/db/enums"
	"github.com/Gleipnir-Technology/nidus-sync/db/models"
	"github.com/Gleipnir-Technology/nidus-sync/h3utils"
	"github.com/Gleipnir-Technology/nidus-sync/htmlpage"
	"github.com/aarondl/opt/omit"
	"github.com/aarondl/opt/omitnull"
	"github.com/rs/zerolog/log"
	"github.com/stephenafamo/bob/dialect/psql"
	"github.com/stephenafamo/bob/dialect/psql/um"
)

type ContextQuick struct{}
type ContextQuickSubmitComplete struct {
	ReportID string
}

var (
	Quick               = buildTemplate("quick", "base")
	QuickSubmitComplete = buildTemplate("quick-submit-complete", "base")
)

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
		Address:  omit.From(""),
		Created:  omit.From(time.Now()),
		Comments: omit.From(comments),
		//Location: omitnull.From(fmt.Sprintf("ST_GeometryFromText(Point(%s %s))", longitude, latitude)),
		H3cell:        omitnull.From(c.String()),
		PublicID:      omit.From(u),
		ReporterEmail: omit.From(""),
		ReporterPhone: omit.From(""),
		Status:        omit.From(enums.PublicreportReportstatustypeReported),
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
			Size:     omit.From(u.Size),
			UUID:     omit.From(u.UUID),
		})
	}
	err = quick.InsertQuickPhotos(r.Context(), db.PGInstance.BobDB, photoSetters...)
	if err != nil {
		respondError(w, "Failed to create photo records", err, http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, fmt.Sprintf("/quick-submit-complete?report=%s", u), http.StatusFound)
}
