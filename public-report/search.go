package publicreport

import (
	"net/http"

	"github.com/Gleipnir-Technology/nidus-sync/htmlpage"
)

type ContextSearch struct{}

var (
	Search = buildTemplate("search", "base")
)

func getSearch(w http.ResponseWriter, r *http.Request) {
	htmlpage.RenderOrError(
		w,
		Search,
		ContextSearch{},
	)
}
