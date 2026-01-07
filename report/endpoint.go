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
	localFS := http.Dir("./static")
	htmlpage.FileServer(r, "/static", localFS, publicreports.EmbeddedStaticFS, "static")
	return r
}

func getRoot(w http.ResponseWriter, r *http.Request) {
	htmlpage.RenderOrError(
		w,
		publicreports.Root,
		publicreports.RootContext{},
	)
}
