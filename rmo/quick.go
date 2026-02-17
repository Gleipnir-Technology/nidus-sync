package rmo

import (
	"net/http"

	"github.com/Gleipnir-Technology/nidus-sync/html"
)

type ContentRegisterNotificationsComplete struct {
	ReportID string
}
type District struct {
	LogoURL string
	Name    string
}

func getRegisterNotificationsComplete(w http.ResponseWriter, r *http.Request) {
	report := r.URL.Query().Get("report")
	html.RenderOrError(
		w,
		"rmo/register-notifications-complete.html",
		ContentRegisterNotificationsComplete{
			ReportID: report,
		},
	)
}
