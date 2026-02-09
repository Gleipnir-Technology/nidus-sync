package rmo

import (
	"github.com/Gleipnir-Technology/nidus-sync/html"
	"github.com/go-chi/chi/v5"
)

func Router() chi.Router {
	r := chi.NewRouter()
	r.Get("/", getRoot)
	r.Get("/nuisance", getNuisance)
	r.Post("/nuisance", postNuisance)
	r.Get("/submit-complete", getSubmitComplete)
	r.Get("/water", getWater)
	r.Post("/water", postWater)

	r.Get("/district", getDistrictList)
	r.Get("/district/{slug}", getRootDistrict)
	r.Get("/district/{slug}/nuisance", getNuisanceDistrict)
	//r.Get("/district/{slug}/nuisance-submit-complete", renderMock(mockNuisanceSubmitCompleteT))
	//r.Get("/district/{slug}/status", renderMock(mockStatusT))
	r.Get("/district/{slug}/water", getWaterDistrict)
	//r.Post("/district/{slug}/water", postWaterDistrict)
	r.Get("/error", getError)

	r.Get("/privacy", getPrivacy)
	r.Get("/robots.txt", getRobots)
	r.Get("/email/render/{code}", getEmailByCode)
	r.Get("/email/confirm", getEmailConfirm)
	r.Post("/email/confirm", postEmailConfirm)
	r.Get("/email/confirm/complete", getEmailConfirmComplete)
	r.Get("/email/unsubscribe", getEmailUnsubscribe)
	r.Get("/email/unsubscribe/report/{report_id}", getEmailReportUnsubscribe)
	r.Get("/image/{uuid}", getImageByUUID)
	r.Route("/mock", addMockRoutes)
	r.Get("/quick", getQuick)
	r.Post("/quick-submit", postQuick)
	r.Get("/quick-submit-complete", getQuickSubmitComplete)
	r.Post("/register-notifications", postRegisterNotifications)
	r.Get("/register-notifications-complete", getRegisterNotificationsComplete)
	r.Get("/report/suggest", getReportSuggestion)
	r.Get("/search", getSearch)
	r.Get("/scss/*", getScssDebug)
	r.Get("/status", getStatus)
	r.Get("/status/{report_id}", getStatusByID)
	r.Get("/terms-of-service", getTerms)
	html.AddStaticRoute(r, "/static")
	return r
}
