package publicreport

import (
	"bytes"
	"io"
	"net/http"

	"github.com/Gleipnir-Technology/nidus-sync/htmlpage"
	"github.com/Gleipnir-Technology/nidus-sync/htmlpage/public-reports"
	"github.com/go-chi/chi/v5"
	"github.com/rs/zerolog/log"
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
	latitude := r.FormValue("latitude")
	longitude := r.FormValue("longitude")
	created := r.FormValue("created")
	//photos := r.FormValue("photos")

	log.Info().Str("latitude", latitude).Str("longitude", longitude).Str("created", created).Msg("Got upload")
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
		}
	}
	http.Redirect(w, r, "/quick-submit-complete?report=123", http.StatusFound)
}

// Respond with an error that is visible to the user
func respondError(w http.ResponseWriter, m string, e error, s int) {
	log.Warn().Int("status", s).Err(e).Str("user message", m).Msg("Responding with an error")
	http.Error(w, m, s)
}
