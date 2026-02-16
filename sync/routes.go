package sync

import (
	"context"
	"net/http"

	"github.com/Gleipnir-Technology/nidus-sync/api"
	"github.com/Gleipnir-Technology/nidus-sync/auth"
	"github.com/Gleipnir-Technology/nidus-sync/db/models"
	"github.com/Gleipnir-Technology/nidus-sync/html"
	"github.com/go-chi/chi/v5"
)

func Router() chi.Router {
	r := chi.NewRouter()
	// Root is a special endpoint that is neither authenticated nor unauthenticated
	r.Get("/", getRoot)

	// Unauthenticated endpoints
	r.Get("/arcgis/oauth/begin", getArcgisOauthBegin)
	r.Get("/arcgis/oauth/callback", getArcgisOauthCallback)
	r.Get("/district", getDistrict)

	// Mock endpoints
	r.Get("/mock", renderMockList)
	addMock(r, "/mock/admin", "sync/mock/admin.html")
	addMock(r, "/mock/dispatch", "sync/mock/dispatch.html")
	addMock(r, "/mock/dispatch-results", "sync/mock/dispatch-results.html")
	addMock(r, "/mock/report", "sync/mock/report.html")
	addMock(r, "/mock/report/{code}", "sync/mock/report-detail.html")
	addMock(r, "/mock/report/{code}/confirm", "sync/mock/report-confirmation.html")
	addMock(r, "/mock/report/{code}/contribute", "sync/mock/report-contribute.html")
	addMock(r, "/mock/report/{code}/evidence", "sync/mock/report-evidence.html")
	addMock(r, "/mock/report/{code}/schedule", "sync/mock/report-schedule.html")
	addMock(r, "/mock/report/{code}/update", "sync/mock/report-update.html")
	addMock(r, "/mock/service-request/{code}", "sync/mock/service-request-detail.html")
	addMock(r, "/mock/setting", "sync/mock/setting.html")
	addMock(r, "/mock/setting/pesticide/add", "sync/mock/setting-pesticide-add.html")
	addMock(r, "/mock/setting/user", "sync/mock/setting-user.html")
	addMock(r, "/mock/setting/user/add", "sync/mock/setting-user-add.html")

	// Utility endpoints
	r.Get("/oauth/refresh", getOAuthRefresh)
	r.Get("/privacy", getPrivacy)
	r.Get("/qr-code/report/{code}", getQRCodeReport)
	r.Get("/signin", getSignin)
	r.Post("/signin", postSignin)
	r.Get("/signup", getSignup)
	r.Post("/signup", postSignup)
	r.Get("/template-test", getTemplateTest)

	// Authenticated endpoints
	r.Route("/api", api.AddRoutes)
	r.Method("GET", "/cell/{cell}", auth.NewEnsureAuth(getCellDetails))
	r.Method("GET", "/layout-test", auth.NewEnsureAuth(getLayoutTest))
	r.Method("GET", "/notification", auth.NewEnsureAuth(getNotificationList))
	r.Method("GET", "/pool", auth.NewEnsureAuth(getPoolList))
	r.Method("GET", "/pool/upload", auth.NewEnsureAuth(getPoolUpload))
	r.Method("GET", "/pool/upload/{id}", auth.NewEnsureAuth(getPoolUploadByID))
	r.Method("POST", "/pool/upload", auth.NewEnsureAuth(postPoolUpload))
	r.Method("GET", "/radar", auth.NewEnsureAuth(getRadar))
	r.Method("GET", "/setting", auth.NewEnsureAuth(getSetting))
	r.Method("GET", "/setting/district", auth.NewEnsureAuth(getSettingDistrict))
	r.Method("GET", "/setting/integration", auth.NewEnsureAuth(getSettingIntegration))
	r.Method("GET", "/setting/pesticide", authenticatedHandler(getSettingPesticide))
	r.Method("GET", "/signout", auth.NewEnsureAuth(getSignout))
	r.Method("GET", "/source/{globalid}", auth.NewEnsureAuth(getSource))
	r.Method("GET", "/stadia", auth.NewEnsureAuth(getStadia))
	r.Method("GET", "/trap/{globalid}", auth.NewEnsureAuth(getTrap))
	r.Method("GET", "/text/{destination}", auth.NewEnsureAuth(getTextMessages))

	html.AddStaticRoute(r, "/static")
	return r
}

type errorWithStatus struct {
	Message string
	Status  int
}

func (e *errorWithStatus) Error() string {
	return e.Message
}

type handlerFunction func(context.Context, *models.User) (string, interface{}, *errorWithStatus)
type wrappedHandler func(http.ResponseWriter, *http.Request)
type contentAuthenticated struct {
	C    interface{}
	URL  ContentURL
	User User
}

// w http.ResponseWriter, r *http.Request, u *models.User) {
func authenticatedHandler(f handlerFunction) http.Handler {
	return auth.NewEnsureAuth(func(w http.ResponseWriter, r *http.Request, u *models.User) {
		ctx := r.Context()
		userContent, err := contentForUser(ctx, u)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		template, content, e := f(ctx, u)
		if err != nil {
			//log.Warn().Int("status", s).Err(e).Str("user message", m).Msg("Responding with an error from sync pages")
			http.Error(w, err.Error(), e.Status)
			return
		}
		html.RenderOrError(w, template, contentAuthenticated{
			C:    content,
			URL:  newContentURL(),
			User: userContent,
		})
	})
}
