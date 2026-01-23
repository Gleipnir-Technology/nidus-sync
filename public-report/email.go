package publicreport

import (
	"net/http"

	"github.com/Gleipnir-Technology/nidus-sync/comms"
	"github.com/go-chi/chi/v5"
)

func getEmailInitial(w http.ResponseWriter, r *http.Request) {
	email := chi.URLParam(r, "email")
	comms.RenderEmailInitial(w, email)
}
func getEmailReportSubscriptionConfirmation(w http.ResponseWriter, r *http.Request) {
	report_id := chi.URLParam(r, "report_id")
	comms.RenderEmailReportConfirmation(w, report_id)
}
