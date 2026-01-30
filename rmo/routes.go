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
	//r.Get("/district/{slug}", renderMock(mockDistrictRootT))
	//r.Get("/district/{slug}/nuisance", renderMock(mockNuisanceT))
	//r.Get("/district/{slug}/nuisance-submit-complete", renderMock(mockNuisanceSubmitCompleteT))
	//r.Get("/district/{slug}/status", renderMock(mockStatusT))
	//r.Get("/district/{slug}/water", renderMock(mockWaterT))

	r.Get("/privacy", getPrivacy)
	r.Get("/robots.txt", getRobots)
	r.Get("/email", getEmailByCode)
	r.Get("/image/{uuid}", getImageByUUID)
	r.Route("/mock", addMockRoutes)
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
	html.AddStaticRoute(r, "/static")
	return r
}
