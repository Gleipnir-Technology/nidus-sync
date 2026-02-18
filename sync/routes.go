package sync

import (
	"context"
	"fmt"
	"net/http"

	"github.com/Gleipnir-Technology/nidus-sync/api"
	"github.com/Gleipnir-Technology/nidus-sync/auth"
	"github.com/Gleipnir-Technology/nidus-sync/db/models"
	"github.com/Gleipnir-Technology/nidus-sync/html"
	"github.com/go-chi/chi/v5"
	"github.com/rs/zerolog/log"
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
	addMock(r, "/mock/report", "sync/mock/report.html")
	addMock(r, "/mock/report/{code}", "sync/mock/report-detail.html")
	addMock(r, "/mock/report/{code}/confirm", "sync/mock/report-confirmation.html")
	addMock(r, "/mock/report/{code}/contribute", "sync/mock/report-contribute.html")
	addMock(r, "/mock/report/{code}/evidence", "sync/mock/report-evidence.html")
	addMock(r, "/mock/report/{code}/schedule", "sync/mock/report-schedule.html")
	addMock(r, "/mock/report/{code}/update", "sync/mock/report-update.html")

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
	r.Method("GET", "/admin", authenticatedHandler(getAdminDash))
	r.Method("GET", "/cell/{cell}", auth.NewEnsureAuth(getCellDetails))
	r.Method("GET", "/download", authenticatedHandler(getDownloadList))
	r.Method("GET", "/layout-test", auth.NewEnsureAuth(getLayoutTest))
	r.Method("GET", "/message", authenticatedHandler(getMessageList))
	r.Method("GET", "/notification", auth.NewEnsureAuth(getNotificationList))
	r.Method("GET", "/pool", auth.NewEnsureAuth(getPoolList))
	r.Method("GET", "/pool/upload", auth.NewEnsureAuth(getPoolUpload))
	r.Method("GET", "/pool/upload/{id}", auth.NewEnsureAuth(getPoolUploadByID))
	r.Method("POST", "/pool/upload", auth.NewEnsureAuth(postPoolUpload))
	r.Method("GET", "/radar", authenticatedHandler(getRadar))
	r.Method("GET", "/service-request", authenticatedHandler(getServiceRequestList))
	r.Method("GET", "/service-request/{id}", authenticatedHandler(getServiceRequestDetail))
	r.Method("GET", "/setting", auth.NewEnsureAuth(getSetting))
	r.Method("GET", "/setting/organization", auth.NewEnsureAuth(getSettingOrganization))
	r.Method("GET", "/setting/integration", auth.NewEnsureAuth(getSettingIntegration))
	r.Method("GET", "/setting/pesticide", authenticatedHandler(getSettingPesticide))
	r.Method("GET", "/setting/pesticide/add", authenticatedHandler(getSettingPesticideAdd))
	r.Method("GET", "/setting/user", authenticatedHandler(getSettingUserList))
	r.Method("GET", "/setting/user/add", authenticatedHandler(getSettingUserAdd))
	r.Method("GET", "/signout", auth.NewEnsureAuth(getSignout))
	r.Method("GET", "/source/{globalid}", auth.NewEnsureAuth(getSource))
	r.Method("GET", "/stadia", auth.NewEnsureAuth(getStadia))
	r.Method("GET", "/sudo", authenticatedHandler(getSudo))
	r.Method("POST", "/sudo/email", authenticatedHandlerPost(postSudoEmail))
	r.Method("POST", "/sudo/sms", authenticatedHandlerPost(postSudoSMS))
	r.Method("GET", "/trap/{globalid}", auth.NewEnsureAuth(getTrap))
	r.Method("GET", "/text/{destination}", auth.NewEnsureAuth(getTextMessages))
	r.Method("GET", "/upload", authenticatedHandler(getUploadList))

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
func newError(mesg_format string, args ...interface{}) *errorWithStatus {
	w := fmt.Errorf(mesg_format, args...)
	return &errorWithStatus{
		Message: w.Error(),
		Status:  http.StatusInternalServerError,
	}
}

type handlerFunctionGet[T any] func(context.Context, *models.User) (string, T, *errorWithStatus)
type wrappedHandler func(http.ResponseWriter, *http.Request)
type contentAuthenticated[T any] struct {
	C    T
	URL  ContentURL
	User User
}

// w http.ResponseWriter, r *http.Request, u *models.User) {
func authenticatedHandler[T any](f handlerFunctionGet[T]) http.Handler {
	return auth.NewEnsureAuth(func(w http.ResponseWriter, r *http.Request, u *models.User) {
		ctx := r.Context()
		userContent, err := contentForUser(ctx, u)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		template, content, e := f(ctx, u)
		//log.Info().Str("template", template).Err(e).Msg("handler done")
		if e != nil {
			log.Warn().Int("status", e.Status).Err(e).Str("user message", e.Message).Msg("Responding with an error from sync pages")
			http.Error(w, e.Error(), e.Status)
			return
		}
		html.RenderOrError(w, template, contentAuthenticated[T]{
			C:    content,
			URL:  newContentURL(),
			User: userContent,
		})
	})
}

type handlerFunctionPost[T any] func(context.Context, *models.User, T) (string, *errorWithStatus)

func authenticatedHandlerPost[T any](f handlerFunctionPost[T]) http.Handler {
	return auth.NewEnsureAuth(func(w http.ResponseWriter, r *http.Request, u *models.User) {
		err := r.ParseForm()
		if err != nil {
			respondError(w, "Failed to parse form", err, http.StatusBadRequest)
			return
		}

		var content T

		err = decoder.Decode(&content, r.PostForm)
		if err != nil {
			respondError(w, "Failed to decode form", err, http.StatusBadRequest)
			return
		}
		ctx := r.Context()
		path, e := f(ctx, u, content)
		if e != nil {
			http.Error(w, e.Error(), e.Status)
			return
		}
		http.Redirect(w, r, path, http.StatusFound)
	})
}
