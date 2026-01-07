package report

import (
	"net/http"

	"github.com/Gleipnir-Technology/nidus-sync/htmlpage"
	"github.com/Gleipnir-Technology/nidus-sync/htmlpage/public-reports"
	"github.com/go-chi/chi/v5"
)

func Router() chi.Router {
	r := chi.NewRouter()
	r.Get("/", getRoot)
	r.Get("/nuisance", getNuisance)
	r.Get("/pool", getPool)
	r.Get("/quick", getQuick)
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
func getStatus(w http.ResponseWriter, r *http.Request) {
	htmlpage.RenderOrError(
		w,
		publicreports.Status,
		publicreports.ContextStatus{},
	)
}
