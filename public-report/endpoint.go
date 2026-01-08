package publicreport

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/Gleipnir-Technology/nidus-sync/db"
	"github.com/Gleipnir-Technology/nidus-sync/db/models"
	"github.com/Gleipnir-Technology/nidus-sync/h3utils"
	"github.com/Gleipnir-Technology/nidus-sync/htmlpage"
	"github.com/Gleipnir-Technology/nidus-sync/htmlpage/public-reports"
	"github.com/Gleipnir-Technology/nidus-sync/userfile"
	"github.com/aarondl/opt/omit"
	"github.com/aarondl/opt/omitnull"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"github.com/stephenafamo/bob/dialect/psql"
	//"github.com/stephenafamo/bob/dialect/psql/dialect"
	//"github.com/stephenafamo/bob/dialect/psql/im"
	"github.com/stephenafamo/bob/dialect/psql/um"
)

func Router() chi.Router {
	r := chi.NewRouter()
	r.Get("/", getRoot)
	r.Get("/nuisance", getNuisance)
	r.Get("/pool", getPool)
	r.Get("/quick", getQuick)
	r.Post("/quick-submit", postQuick)
	r.Get("/quick-submit-complete", getQuickSubmitComplete)
	r.Get("/status", getStatus)
	localFS := http.Dir("./static")
	htmlpage.FileServer(r, "/static", localFS, publicreports.EmbeddedStaticFS, "static")
	return r
}

func getRoot(w http.ResponseWriter, r *http.Request) {
	htmlpage.RenderOrError(
		w,
		publicreports.Root,
		publicreports.ContextRoot{},
	)
}

func getNuisance(w http.ResponseWriter, r *http.Request) {
	htmlpage.RenderOrError(
		w,
		publicreports.Nuisance,
		publicreports.ContextNuisance{},
	)
}
func getPool(w http.ResponseWriter, r *http.Request) {
	htmlpage.RenderOrError(
		w,
		publicreports.Pool,
		publicreports.ContextPool{},
	)
}
func getQuick(w http.ResponseWriter, r *http.Request) {
	htmlpage.RenderOrError(
		w,
		publicreports.Quick,
		publicreports.ContextQuick{},
	)
}
func getQuickSubmitComplete(w http.ResponseWriter, r *http.Request) {
	report := r.URL.Query().Get("report")
	htmlpage.RenderOrError(
		w,
		publicreports.QuickSubmitComplete,
		publicreports.ContextQuickSubmitComplete{
			ReportID: report,
		},
	)
}
func getStatus(w http.ResponseWriter, r *http.Request) {
	htmlpage.RenderOrError(
		w,
		publicreports.Status,
		publicreports.ContextStatus{},
	)
}
func postQuick(w http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(32 << 10) // 32 MB buffer
	if err != nil {
		log.Error().Err(err).Msg("Failed to parse form")
		respondError(w, "Failed to parse form", err, http.StatusBadRequest)
		return
	}
	lat := r.FormValue("latitude")
	lng := r.FormValue("longitude")
	created := r.FormValue("created")
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
	u, err := uuid.NewUUID()
	if err != nil {
		respondError(w, "Failed to create quick report uuid", err, http.StatusInternalServerError)
		return
	}
	c, err := h3utils.GetCell(longitude, latitude, 15)
	setter := models.PublicreportQuickSetter{
		Created:  omit.From(time.Now()),
		Comments: omit.From(comments),
		//Location: omitnull.From(fmt.Sprintf("ST_GeometryFromText(Point(%s %s))", longitude, latitude)),
		H3cell: omitnull.From(c.String()),
		UUID:   omit.From(u),
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
	log.Info().Float64("latitude", latitude).Float64("longitude", longitude).Str("created", created).Msg("Got upload")
	photoSetters := make([]*models.PublicreportQuickPhotoSetter, 0)
	for _, fheaders := range r.MultipartForm.File {
		for _, headers := range fheaders {
			file, err := headers.Open()

			if err != nil {
				respondError(w, "Failed to open header", err, http.StatusInternalServerError)
				return
			}

			defer file.Close()

			buff := make([]byte, 512)
			file.Read(buff)

			file.Seek(0, 0)
			contentType := http.DetectContentType(buff)
			var sizeBuff bytes.Buffer
			fileSize, err := sizeBuff.ReadFrom(file)
			if err != nil {
				respondError(w, "Failed to read file", err, http.StatusInternalServerError)
				return
			}
			file.Seek(0, 0)
			contentBuf := bytes.NewBuffer(nil)
			if _, err := io.Copy(contentBuf, file); err != nil {
				respondError(w, "Failed to save file", err, http.StatusInternalServerError)
				return
			}
			log.Info().Int64("size", fileSize).Str("filename", headers.Filename).Str("content-type", contentType).Msg("Got an uploaded file")
			u, err := uuid.NewUUID()
			if err != nil {
				respondError(w, "Failed to create quick report photo uuid", err, http.StatusInternalServerError)
				continue
			}
			err = userfile.PublicImageFileContentWrite(u, file)
			photoSetters = append(photoSetters, &models.PublicreportQuickPhotoSetter{
				Size:     omit.From(fileSize),
				Filename: omit.From(headers.Filename),
				UUID:     omit.From(u),
			})
		}
	}
	/*err = quick.InsertQuickPhotos(r.Context(), db.PGInstance.BobDB, photoSetters...)
	if err != nil {
		respondError(w, "Failed to create photo records", err, http.StatusInternalServerError)
		return
	}*/
	http.Redirect(w, r, "/quick-submit-complete?report=123", http.StatusFound)
}

// Respond with an error that is visible to the user
func respondError(w http.ResponseWriter, m string, e error, s int) {
	log.Warn().Int("status", s).Err(e).Str("user message", m).Msg("Responding with an error")
	http.Error(w, m, s)
}
