package publicreport

import (
	"net/http"

	"github.com/Gleipnir-Technology/nidus-sync/config"
	"github.com/Gleipnir-Technology/nidus-sync/htmlpage"
)

type ContentSearch struct {
	MapboxToken string
	URLTegola   string
}

var (
	Search = buildTemplate("search", "base")
)

func getSearch(w http.ResponseWriter, r *http.Request) {
	htmlpage.RenderOrError(
		w,
		Search,
		ContentSearch{
			MapboxToken: config.MapboxToken,
			URLTegola:   config.MakeURLTegola("/"),
		},
	)
}
