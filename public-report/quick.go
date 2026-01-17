package publicreport

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/Gleipnir-Technology/nidus-sync/comms"
	"github.com/Gleipnir-Technology/nidus-sync/db"
	"github.com/Gleipnir-Technology/nidus-sync/db/enums"
	"github.com/Gleipnir-Technology/nidus-sync/db/models"
	"github.com/Gleipnir-Technology/nidus-sync/h3utils"
	"github.com/Gleipnir-Technology/nidus-sync/htmlpage"
	"github.com/aarondl/opt/omit"
	"github.com/aarondl/opt/omitnull"
	"github.com/rs/zerolog/log"
	"github.com/stephenafamo/bob"
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
func getRegisterNotificationsComplete(w http.ResponseWriter, r *http.Request) {
	report := r.URL.Query().Get("report")
	htmlpage.RenderOrError(
		w,
		RegisterNotificationsComplete,
		ContextRegisterNotificationsComplete{
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
	ctx := r.Context()
	tx, err := db.PGInstance.BobDB.BeginTx(ctx, nil)
	if err != nil {
		respondError(w, "Failed to create transaction", err, http.StatusInternalServerError)
		return
	}
	defer tx.Rollback(ctx)

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
	quick, err := models.PublicreportQuicks.Insert(&setter).One(ctx, tx)
	if err != nil {
		respondError(w, "Failed to create database record", err, http.StatusInternalServerError)
		return
	}
	_, err = psql.Update(
		um.Table("publicreport.quick"),
		um.SetCol("location").To(fmt.Sprintf("ST_GeometryFromText('Point(%f %f)')", longitude, latitude)),
		um.Where(psql.Quote("id").EQ(psql.Arg(quick.ID))),
	).Exec(ctx, tx)
	if err != nil {
		respondError(w, "Failed to insert publicreport", err, http.StatusInternalServerError)
		return
	}
	log.Info().Float64("latitude", latitude).Float64("longitude", longitude).Msg("Got upload")
	uploads, err := extractImageUploads(r)
	log.Info().Int("len", len(uploads)).Msg("extracted uploads")
	if err != nil {
		respondError(w, "Failed to extract image uploads", err, http.StatusInternalServerError)
		return
	}
	images, err := saveImageUploads(ctx, tx, uploads)
	if err != nil {
		respondError(w, "Failed to save image uploads", err, http.StatusInternalServerError)
		return
	}
	log.Info().Int("len", len(images)).Msg("saved uploads")
	setters := make([]*models.PublicreportQuickImageSetter, 0)
	for _, image := range images {
		setters = append(setters, &models.PublicreportQuickImageSetter{
			ImageID: omit.From(int32(image.ID)),
			QuickID: omit.From(int32(quick.ID)),
		})
	}
	_, err = models.PublicreportQuickImages.Insert(bob.ToMods(setters...)).Exec(ctx, tx)
	if err != nil {
		respondError(w, "Failed to save reference to images", err, http.StatusInternalServerError)
		return
	}
	log.Info().Int("len", len(images)).Msg("saved uploads")
	tx.Commit(ctx)
	http.Redirect(w, r, fmt.Sprintf("/quick-submit-complete?report=%s", u), http.StatusFound)
}
func postRegisterNotifications(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		respondError(w, "Failed to parse form", err, http.StatusBadRequest)
		return
	}
	consent := r.PostFormValue("consent")
	email := r.PostFormValue("email")
	phone := r.PostFormValue("phone")
	report_id := r.PostFormValue("report_id")
	if consent != "on" {
		respondError(w, "You must consent", nil, http.StatusBadRequest)
		return
	}
	if email == "" && phone == "" {
		http.Redirect(w, r, fmt.Sprintf("/quick-submit-complete?report=%s", report_id), http.StatusFound)
		return
	}
	result, err := psql.Update(
		um.Table("publicreport.quick"),
		um.SetCol("reporter_email").ToArg(email),
		um.SetCol("reporter_phone").ToArg(phone),
		um.Where(psql.Quote("public_id").EQ(psql.Arg(report_id))),
	).Exec(r.Context(), db.PGInstance.BobDB)
	if err != nil {
		respondError(w, "Failed to update report", err, http.StatusInternalServerError)
		return
	}
	rowcount, err := result.RowsAffected()
	if err != nil {
		respondError(w, "Failed to get rows affected", err, http.StatusInternalServerError)
		return
	}
	if email != "" {
		comms.SendEmail(comms.EmailRequest{
			From: "website@mosquitoes.online",
			To: email,
			Subject: "test email",
			Text: "This is just testing that I can send email",
		})
	}
	if phone != "" {
		err := comms.SendSMS(phone, "testing 1 2 3")
		if err != nil {
			log.Error().Err(err).Msg("Failed to send SMS")
		}
	}
	if rowcount == 0 {
		http.Redirect(w, r, fmt.Sprintf("/error?code=no-rows-affected&report=%s", report_id), http.StatusFound)
	} else {
		http.Redirect(w, r, fmt.Sprintf("/register-notifications-complete?report=%s", report_id), http.StatusFound)
	}
}

