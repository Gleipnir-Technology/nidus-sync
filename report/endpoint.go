package report

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func Router() chi.Router {
	r := chi.NewRouter()
	r.Get("/", getRoot)
	return r
}

func getRoot(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Herro.")
}
