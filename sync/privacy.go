package sync

import (
	"net/http"

	"github.com/Gleipnir-Technology/nidus-sync/config"
	"github.com/Gleipnir-Technology/nidus-sync/htmlpage"
)

type ContentPrivacy struct {
	Address string
	Company string
	Site    string
	URLSync string
}

var (
	PrivacyT = buildTemplate("privacy", "base")
)

func getPrivacy(w http.ResponseWriter, r *http.Request) {
	htmlpage.RenderOrError(
		w,
		PrivacyT,
		ContentPrivacy{
			Address: "2726 S Quinn Ave, Gilbert, AZ, USA",
			Company: "Gleipnir LLC",
			Site:    "Nidus Sync",
			URLSync: config.MakeURLNidus("/"),
		},
	)
}
