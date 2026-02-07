package rmo

import (
	"net/http"

	"github.com/Gleipnir-Technology/nidus-sync/config"
	"github.com/Gleipnir-Technology/nidus-sync/html"
)

type ContentSearch struct {
	MapboxToken string
	URLTegola   string
}

func getSearch(w http.ResponseWriter, r *http.Request) {
	html.RenderOrError(
		w,
		"rmo/search.html",
		ContentSearch{
			MapboxToken: config.MapboxToken,
			URLTegola:   config.MakeURLTegola("/"),
		},
	)
}
