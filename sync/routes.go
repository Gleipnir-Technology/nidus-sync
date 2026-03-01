package sync

import (
	"context"
	"fmt"
	"net/http"

	"github.com/Gleipnir-Technology/nidus-sync/api"
	"github.com/Gleipnir-Technology/nidus-sync/auth"
	"github.com/Gleipnir-Technology/nidus-sync/db"
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
	r.Get("/mailer/{code}", getMailer)
	r.Get("/mailer/{code}/preview", getMailerPreview)
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
	r.Get("/privacy", getPrivacy)
	r.Get("/qr-code/report/{code}", getQRCodeReport)
	r.Get("/qr-code/mailer/{code}", getQRCodeMailer)
	r.Get("/signin", getSignin)
	r.Post("/signin", postSignin)
	r.Get("/signup", getSignup)
	r.Post("/signup", postSignup)
	r.Get("/template-test", getTemplateTest)

	// Authenticated endpoints
	r.Route("/api", api.AddRoutes)

	r.Method("GET", "/admin", authenticatedHandler(getAdminDash))
	r.Method("GET", "/cell/{cell}", authenticatedHandler(getCellDetails))
	r.Method("GET", "/communication", authenticatedHandler(getCommunicationRoot))
	r.Method("GET", "/configuration", authenticatedHandler(getConfigurationRoot))
	r.Method("GET", "/configuration/integration", authenticatedHandler(getConfigurationIntegration))
	r.Method("GET", "/configuration/integration/arcgis", authenticatedHandler(getConfigurationIntegrationArcgis))
	r.Method("GET", "/configuration/organization", authenticatedHandler(getConfigurationOrganization))
	r.Method("GET", "/configuration/pesticide", authenticatedHandler(getConfigurationPesticide))
	r.Method("GET", "/configuration/pesticide/add", authenticatedHandler(getConfigurationPesticideAdd))
	r.Method("GET", "/configuration/upload", authenticatedHandler(getUploadList))
	r.Method("GET", "/configuration/upload/pool", authenticatedHandler(getUploadPool))
	r.Method("GET", "/configuration/upload/pool/bob", authenticatedHandler(getUploadPoolBobCreate))
	r.Method("POST", "/configuration/upload/pool/bob", authenticatedHandlerPostMultipart(postUploadPoolBobCreate))
	r.Method("GET", "/configuration/upload/pool/custom", authenticatedHandler(getUploadPoolCustomCreate))
	r.Method("POST", "/configuration/upload/pool/custom", authenticatedHandlerPostMultipart(postUploadPoolCustomCreate))
	r.Method("GET", "/configuration/upload/{id}", authenticatedHandler(getUploadByID))
	r.Method("POST", "/configuration/upload/{id}/discard", authenticatedHandlerPost(postUploadDiscard))
	r.Method("GET", "/configuration/user", authenticatedHandler(getConfigurationUserList))
	r.Method("GET", "/configuration/user/add", authenticatedHandler(getConfigurationUserAdd))
	r.Method("GET", "/download", authenticatedHandler(getDownloadList))
	r.Method("GET", "/intelligence", authenticatedHandler(getIntelligenceRoot))
	r.Method("GET", "/layout-test", authenticatedHandler(getLayoutTest))
	r.Method("GET", "/message", authenticatedHandler(getMessageList))
	r.Method("GET", "/notification", authenticatedHandler(getNotificationList))
	r.Method("GET", "/oauth/refresh", authenticatedHandler(getOAuthRefresh))
	r.Method("GET", "/operations", authenticatedHandler(getOperationsRoot))
	r.Method("GET", "/parcel", authenticatedHandler(getParcel))
	r.Method("GET", "/planning", authenticatedHandler(getPlanningRoot))
	r.Method("GET", "/pool", authenticatedHandler(getPoolList))
	r.Method("GET", "/pool/create", authenticatedHandler(getPoolCreate))
	r.Method("GET", "/pool/{id}", authenticatedHandler(getPoolByID))
	r.Method("GET", "/radar", authenticatedHandler(getRadar))
	r.Method("GET", "/review", authenticatedHandler(getReviewRoot))
	r.Method("GET", "/service-request", authenticatedHandler(getServiceRequestList))
	r.Method("GET", "/service-request/{id}", authenticatedHandler(getServiceRequestDetail))
	r.Method("GET", "/signout", auth.NewEnsureAuth(getSignout))
	r.Method("GET", "/source/{globalid}", authenticatedHandler(getSource))
	r.Method("GET", "/stadia", authenticatedHandler(getStadia))
	r.Method("GET", "/sudo", authenticatedHandler(getSudo))
	r.Method("POST", "/sudo/email", authenticatedHandlerPost(postSudoEmail))
	r.Method("POST", "/sudo/sms", authenticatedHandlerPost(postSudoSMS))
	r.Method("GET", "/trap/{globalid}", authenticatedHandler(getTrap))
	r.Method("GET", "/text/{destination}", authenticatedHandler(getTextMessages))
	r.Method("GET", "/tile/gps", auth.NewEnsureAuth(getTileGPS))

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
	return newErrorStatus(http.StatusInternalServerError, mesg_format, args...)
}
func newErrorMaybe(mesg_format string, err error, args ...interface{}) *errorWithStatus {
	if err == nil {
		return nil
	}
	allArgs := append([]interface{}{err}, args...)
	return newErrorStatus(http.StatusInternalServerError, mesg_format, allArgs...)
}
func newErrorStatus(status int, mesg_format string, args ...interface{}) *errorWithStatus {
	w := fmt.Errorf(mesg_format, args...)
	return &errorWithStatus{
		Message: w.Error(),
		Status:  status,
	}
}

type response[T any] struct {
	content  T
	template string
}

func newResponse[T any](template string, content T) *response[T] {
	return &response[T]{
		content:  content,
		template: template,
	}
}

type handlerFunctionGet[T any] func(context.Context, *http.Request, *models.Organization, *models.User) (*response[T], *errorWithStatus)
type wrappedHandler func(http.ResponseWriter, *http.Request)
type contentAuthenticated[T any] struct {
	C            T
	Config       contentConfig
	Organization *models.Organization
	URL          contentURL
	User         User
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
		org, err := u.Organization().One(ctx, db.PGInstance.BobDB)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		resp, e := f(ctx, r, org, u)
		//log.Info().Str("template", template).Err(e).Msg("handler done")
		if e != nil {
			log.Warn().Int("status", e.Status).Err(e).Str("user message", e.Message).Msg("Responding with an error from sync pages")
			http.Error(w, e.Error(), e.Status)
			return
		}
		if org == nil {
			http.Error(w, "nil org", http.StatusInternalServerError)
			return
		}
		html.RenderOrError(w, resp.template, contentAuthenticated[T]{
			C:            resp.content,
			Config:       newContentConfig(),
			Organization: org,
			URL:          newContentURL(),
			User:         userContent,
		})
	})
}

type handlerFunctionPost[T any] func(context.Context, *http.Request, *models.Organization, *models.User, T) (string, *errorWithStatus)

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
		org, err := u.Organization().One(ctx, db.PGInstance.BobDB)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		path, e := f(ctx, r, org, u, content)
		if e != nil {
			http.Error(w, e.Error(), e.Status)
			return
		}
		http.Redirect(w, r, path, http.StatusFound)
	})
}
func authenticatedHandlerPostMultipart[T any](f handlerFunctionPost[T]) http.Handler {
	return auth.NewEnsureAuth(func(w http.ResponseWriter, r *http.Request, u *models.User) {
		err := r.ParseMultipartForm(32 << 10) // 32 MB buffer
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
		org, err := u.Organization().One(ctx, db.PGInstance.BobDB)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		path, e := f(ctx, r, org, u, content)
		if e != nil {
			http.Error(w, e.Error(), e.Status)
			return
		}
		http.Redirect(w, r, path, http.StatusFound)
	})
}
