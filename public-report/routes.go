package publicreport

import (
	"github.com/Gleipnir-Technology/nidus-sync/htmlpage"
	"github.com/go-chi/chi/v5"
)

func Router() chi.Router {
	r := chi.NewRouter()
	r.Get("/", getRoot)
	r.Get("/privacy", getPrivacy)
	r.Get("/robots.txt", getRobots)
	r.Get("/email/report/{report_id}/subscription-confirmation", getEmailReportSubscriptionConfirmation)
	r.Get("/nuisance", getNuisance)
	r.Post("/nuisance-submit", postNuisance)
	r.Get("/nuisance-submit-complete", getNuisanceSubmitComplete)
	r.Get("/pool", getPool)
	r.Post("/pool-submit", postPool)
	r.Get("/pool-submit-complete", getPoolSubmitComplete)
	r.Get("/quick", getQuick)
	r.Post("/quick-submit", postQuick)
	r.Get("/quick-submit-complete", getQuickSubmitComplete)
	r.Post("/register-notifications", postRegisterNotifications)
	r.Get("/register-notifications-complete", getRegisterNotificationsComplete)
	r.Get("/search", getSearch)
	r.Get("/status", getStatus)
	r.Get("/status/{report_id}", getStatusByID)
	r.Get("/terms-of-service", getTerms)
	htmlpage.AddStaticRoute(r, "/static")
	return r
}
