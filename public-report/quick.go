package publicreport

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/Gleipnir-Technology/nidus-sync/background"
	"github.com/Gleipnir-Technology/nidus-sync/comms"
	"github.com/Gleipnir-Technology/nidus-sync/config"
	"github.com/Gleipnir-Technology/nidus-sync/db"
	"github.com/Gleipnir-Technology/nidus-sync/db/enums"
	"github.com/Gleipnir-Technology/nidus-sync/db/models"
	"github.com/Gleipnir-Technology/nidus-sync/h3utils"
	"github.com/Gleipnir-Technology/nidus-sync/htmlpage"
	"github.com/Gleipnir-Technology/nidus-sync/platform"
	"github.com/aarondl/opt/omit"
	"github.com/aarondl/opt/omitnull"
	"github.com/rs/zerolog/log"
	"github.com/stephenafamo/bob"
	"github.com/stephenafamo/bob/dialect/psql"
	"github.com/stephenafamo/bob/dialect/psql/um"
)

type ContentQuick struct{}
type ContentQuickSubmitComplete struct {
	District *District
	ReportID string
}
type ContentRegisterNotificationsComplete struct {
	ReportID string
}
type District struct {
	LogoURL string
	Name    string
}

var (
	quickT                         = buildTemplate("quick", "base")
	quickSubmitCompleteT           = buildTemplate("quick-submit-complete", "base")
	registerNotificationsCompleteT = buildTemplate("register-notifications-complete", "base")
)

func getQuick(w http.ResponseWriter, r *http.Request) {
	htmlpage.RenderOrError(
		w,
		quickT,
		ContentQuick{},
	)
}
func getQuickSubmitComplete(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	report_id := r.URL.Query().Get("report")
	report, err := models.PublicreportQuicks.Query(
		models.SelectWhere.PublicreportQuicks.PublicID.EQ(report_id),
	).One(ctx, db.PGInstance.BobDB)
	if err != nil {
		respondError(w, "Failed to get report", err, http.StatusInternalServerError)
		return
	}
	var district *District
	if !report.OrganizationID.IsNull() {
		org_id := report.OrganizationID.MustGet()
		org, err := models.Organizations.Query(
			models.Preload.Organization.ImportDistrictGidDistrict(),
			models.SelectWhere.Organizations.ID.EQ(org_id),
		).One(ctx, db.PGInstance.BobDB)
		if err != nil {
			respondError(w, "Failed to get org", err, http.StatusInternalServerError)
			return
		}
		d := org.R.ImportDistrictGidDistrict
		log.Debug().Int32("org_id", org.ID).Int32("d_gid", d.Gid).Msg("Getting district")
		if d != nil {
			district = &District{
				LogoURL: config.MakeURLNidus("/api/district/%s/logo", org.Slug.GetOr("placeholder")),
				Name:    d.Agency.GetOr("Unknown"),
			}
		}
	}
	htmlpage.RenderOrError(
		w,
		quickSubmitCompleteT,
		ContentQuickSubmitComplete{
			District: district,
			ReportID: report.PublicID,
		},
	)
}
func getRegisterNotificationsComplete(w http.ResponseWriter, r *http.Request) {
	report := r.URL.Query().Get("report")
	htmlpage.RenderOrError(
		w,
		registerNotificationsCompleteT,
		ContentRegisterNotificationsComplete{
			ReportID: report,
		},
	)
}
func matchDistrict(ctx context.Context, longitude, latitude float64, images []ImageUpload) (*int32, error) {
	for _, image := range images {
		if image.Exif.GPS == nil {
			continue
		}
		_, org, err := platform.DistrictForLocation(ctx, image.Exif.GPS.Longitude, image.Exif.GPS.Latitude)
		if err != nil {
			log.Warn().Err(err).Msg("Failed to get district for location")
			continue
		}
		if org != nil {
			return &org.ID, nil
		}
	}
	_, org, err := platform.DistrictForLocation(ctx, longitude, latitude)
	if err != nil {
		log.Warn().Err(err).Msg("Failed to get district for location")
		return nil, fmt.Errorf("Failed to get district for location: %w", err)
	}
	if org == nil {
		log.Debug().Err(err).Float64("lng", longitude).Float64("lat", latitude).Msg("No district match by report location")
		return nil, nil
	}
	log.Debug().Err(err).Int32("org_id", org.ID).Float64("lng", longitude).Float64("lat", latitude).Msg("Found district match by report location")
	return &org.ID, nil
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

	organization_id, err := matchDistrict(ctx, longitude, latitude, uploads)
	if err != nil {
		log.Warn().Err(err).Msg("Failed to match district")
	}

	log.Info().Int("len", len(images)).Msg("saved uploads")
	c, err := h3utils.GetCell(longitude, latitude, 15)
	setter := models.PublicreportQuickSetter{
		Address:        omit.From(""),
		Created:        omit.From(time.Now()),
		Comments:       omit.From(comments),
		OrganizationID: omitnull.FromPtr(organization_id),
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
	if len(images) > 0 {
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
	}
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
	phone_str := r.PostFormValue("phone")
	report_id := r.PostFormValue("report_id")
	if consent != "on" {
		respondError(w, "You must consent", nil, http.StatusBadRequest)
		return
	}
	if email == "" && phone_str == "" {
		http.Redirect(w, r, fmt.Sprintf("/quick-submit-complete?report=%s", report_id), http.StatusFound)
		return
	}
	phone, err := comms.ParsePhoneNumber(phone_str)
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
		background.ReportSubscriptionConfirmationEmail(email, report_id)
	}
	if phone_str != "" {
		background.ReportSubscriptionConfirmationText(*phone, report_id)
	}
	if rowcount == 0 {
		http.Redirect(w, r, fmt.Sprintf("/error?code=no-rows-affected&report=%s", report_id), http.StatusFound)
	} else {
		http.Redirect(w, r, fmt.Sprintf("/register-notifications-complete?report=%s", report_id), http.StatusFound)
	}
}
