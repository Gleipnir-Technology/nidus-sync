package html

import (
	"embed"
	"net/http"

	"github.com/go-chi/chi/v5"
)

//go:embed static/*
var EmbeddedStaticFS embed.FS

var localFS http.Dir

func AddStaticRoute(r chi.Router, path string) {
	if localFS == "" {
		localFS = http.Dir("./html/static")
	}
	FileServer(r, "/static", localFS, EmbeddedStaticFS, "static")
}
