package sync

import (
	"net/http"

	"github.com/Gleipnir-Technology/nidus-sync/config"
	"github.com/Gleipnir-Technology/nidus-sync/html"
)

type ContentPrivacy struct {
	Address string
	Company string
	Site    string
	URLSync string
}

func getPrivacy(w http.ResponseWriter, r *http.Request) {
	html.RenderOrError(
		w,
		"sync/privacy.html",
		ContentPrivacy{
			Address: "2726 S Quinn Ave, Gilbert, AZ, USA",
			Company: "Gleipnir LLC",
			Site:    "Nidus Sync",
			URLSync: config.MakeURLNidus("/"),
		},
	)
}
