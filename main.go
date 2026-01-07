package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"runtime/debug"
	"sync"
	"syscall"
	"time"

	"github.com/Gleipnir-Technology/nidus-sync/api"
	"github.com/Gleipnir-Technology/nidus-sync/auth"
	"github.com/Gleipnir-Technology/nidus-sync/background"
	"github.com/Gleipnir-Technology/nidus-sync/config"
	"github.com/Gleipnir-Technology/nidus-sync/db"
	"github.com/Gleipnir-Technology/nidus-sync/queue"
	"github.com/Gleipnir-Technology/nidus-sync/report"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/hostrouter"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	err := config.Parse()
	if err != nil {
		log.Error().Err(err).Msg("Failed to parse config")
		os.Exit(1)
	}
	log.Info().Msg("Starting...")
	err = db.InitializeDatabase(context.TODO(), config.PGDSN)
	if err != nil {
		log.Error().Err(err).Msg("Failed to connect to database")
		os.Exit(2)
	}

	router_logger := log.With().Logger()
	r := chi.NewRouter()

	r.Use(LoggerMiddleware(&router_logger))
	r.Use(middleware.RealIP)
	r.Use(auth.NewSessionManager().LoadAndSave)

	hr := hostrouter.New()

	// Set up routing by hostname
	sr := syncRouter()
	hr.Map("", sr)                            // default
	hr.Map("*", sr)                           // default
	hr.Map(config.URLReport, report.Router()) // report.mosquitoes.online
	hr.Map(config.URLSync, sr)
	r.Mount("/", hr)

	log.Info().Str("report url", config.URLReport).Str("sync url", config.URLSync).Msg("Serving at URLs")
	// Start up background processes
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	background.NewOAuthTokenChannel = make(chan struct{}, 10)

	var waitGroup sync.WaitGroup

	waitGroup.Add(1)
	go func() {
		defer waitGroup.Done()
		background.RefreshFieldseekerData(ctx, background.NewOAuthTokenChannel)
	}()

	waitGroup.Add(1)
	go func() {
		defer waitGroup.Done()
		queue.StartAudioWorker(ctx)
	}()

	server := &http.Server{
		Addr:    config.Bind,
		Handler: r,
	}
	go func() {
		log.Info().Str("address", config.Bind).Msg("Serving HTTP requests")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Error().Str("err", err.Error()).Msg("HTTP Server Error")
		}
	}()

	// Wait for the interrupt signal to gracefully shut down
	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, syscall.SIGINT, syscall.SIGTERM)
	<-signalCh

	log.Info().Msg("Received shutdown signal, shutting down...")
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdownCancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Error().Str("err", err.Error()).Msg("HTTP server shutdown error")
	}

	cancel()

	waitGroup.Wait()

	log.Info().Msg("Shutdown complete")
}
func syncRouter() chi.Router {
	r := chi.NewRouter()
	// Root is a special endpoint that is neither authenticated nor unauthenticated
	r.Get("/", getRoot)

	// Unauthenticated endpoints
	r.Get("/arcgis/oauth/begin", getArcgisOauthBegin)
	r.Get("/arcgis/oauth/callback", getArcgisOauthCallback)
	r.Get("/favicon.ico", getFavicon)

	r.Get("/mock", renderMock("mock-root"))
	r.Get("/mock/admin", renderMock("admin"))
	r.Get("/mock/admin/service-request", renderMock("admin-service-request"))
	r.Get("/mock/data-entry", renderMock("data-entry"))
	r.Get("/mock/data-entry/bad", renderMock("data-entry-bad"))
	r.Get("/mock/data-entry/good", renderMock("data-entry-good"))
	r.Get("/mock/dispatch", renderMock("dispatch"))
	r.Get("/mock/dispatch-results", renderMock("dispatch-results"))
	r.Get("/mock/report", renderMock("report"))
	r.Get("/mock/report/{code}", renderMock("report-detail"))
	r.Get("/mock/report/{code}/confirm", renderMock("report-confirmation"))
	r.Get("/mock/report/{code}/contribute", renderMock("report-contribute"))
	r.Get("/mock/report/{code}/evidence", renderMock("report-evidence"))
	r.Get("/mock/report/{code}/schedule", renderMock("report-schedule"))
	r.Get("/mock/report/{code}/update", renderMock("report-update"))
	r.Get("/mock/service-request", renderMock("service-request"))
	r.Get("/mock/service-request/{code}", renderMock("service-request-detail"))
	r.Get("/mock/service-request-location", renderMock("service-request-location"))
	r.Get("/mock/service-request-mosquito", renderMock("service-request-mosquito"))
	r.Get("/mock/service-request-pool", renderMock("service-request-pool"))
	r.Get("/mock/service-request-quick", renderMock("service-request-quick"))
	r.Get("/mock/service-request-quick-confirmation", renderMock("service-request-quick-confirmation"))
	r.Get("/mock/service-request-updates", renderMock("service-request-updates"))
	r.Get("/mock/setting", renderMock("setting-mock"))
	r.Get("/mock/setting/integration", renderMock("setting-integration"))
	r.Get("/mock/setting/pesticide", renderMock("setting-pesticide"))
	r.Get("/mock/setting/pesticide/add", renderMock("setting-pesticide-add"))
	r.Get("/mock/setting/user", renderMock("setting-user"))
	r.Get("/mock/setting/user/add", renderMock("setting-user-add"))

	r.Get("/oauth/refresh", getOAuthRefresh)

	r.Get("/qr-code/report/{code}", getQRCodeReport)
	r.Get("/signin", getSignin)
	r.Post("/signin", postSignin)
	r.Get("/signup", getSignup)
	r.Post("/signup", postSignup)
	r.Get("/sms", getSMS)
	r.Post("/sms", postSMS)
	r.Get("/sms.php", getSMS)
	r.Get("/sms/{org}", getSMS)
	r.Post("/sms/{org}", postSMS)

	// Authenticated endpoints
	r.Route("/api", api.AddRoutes)
	r.Method("GET", "/cell/{cell}", auth.NewEnsureAuth(getCellDetails))
	r.Method("GET", "/settings", auth.NewEnsureAuth(getSettings))
	r.Method("GET", "/source/{globalid}", auth.NewEnsureAuth(getSource))
	//r.Method("GET", "/vector-tiles/{org_id}/{tileset_id}/{zoom}/{x}/{y}.{format}", auth.NewEnsureAuth(getVectorTiles))

	localFS := http.Dir("./static")
	FileServer(r, "/static", localFS, embeddedStaticFS, "static")
	return r
}

func LoggerMiddleware(logger *zerolog.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			log := logger.With().Logger()

			ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)

			t1 := time.Now()
			defer func() {
				t2 := time.Now()

				// Recover and record stack traces in case of a panic
				if rec := recover(); rec != nil {
					log.Error().
						Str("type", "error").
						Timestamp().
						Interface("recover_info", rec).
						Bytes("debug_stack", debug.Stack()).
						Msg("log system error")
					fmt.Println("Stack:", string(debug.Stack()))
					http.Error(ww, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
				}

				remote_addr := r.RemoteAddr
				forwarded_for := r.Header.Get("X-Forwarded-For")
				if forwarded_for != "" {
					remote_addr = forwarded_for
				}
				// log end request
				log.Info().
					//Str("type", "access").
					Timestamp().
					Fields(map[string]interface{}{
						"remote_ip": remote_addr,
						"url":       r.URL.Path,
						//"proto":      r.Proto,
						"method": r.Method,
						//"user_agent": r.Header.Get("User-Agent"),
						"status":     ww.Status(),
						"latency_ms": float64(t2.Sub(t1).Nanoseconds()) / 1000000.0,
						"bytes_in":   r.Header.Get("Content-Length"),
						"bytes_out":  ww.BytesWritten(),
					}).
					Msg("incoming_request")
			}()

			next.ServeHTTP(ww, r)
		}
		return http.HandlerFunc(fn)
	}
}
