package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"runtime/debug"
	"sync"
	"syscall"
	"time"

	"github.com/Gleipnir-Technology/nidus-sync/db"
	"github.com/alexedwards/scs/pgxstore"
	"github.com/alexedwards/scs/v2"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

var sessionManager *scs.SessionManager

var BaseURL, ClientID, ClientSecret, Environment, MapboxToken string

func main() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	ClientID = os.Getenv("ARCGIS_CLIENT_ID")
	if ClientID == "" {
		log.Error().Msg("You must specify a non-empty ARCGIS_CLIENT_ID")
		os.Exit(1)
	}
	ClientSecret = os.Getenv("ARCGIS_CLIENT_SECRET")
	if ClientSecret == "" {
		log.Error().Msg("You must specify a non-empty ARCGIS_CLIENT_SECRET")
		os.Exit(1)
	}
	BaseURL = os.Getenv("BASE_URL")
	if BaseURL == "" {
		log.Error().Msg("You must specify a non-empty BASE_URL")
		os.Exit(1)
	}
	bind := os.Getenv("BIND")
	if bind == "" {
		bind = ":9001"
	}
	Environment = os.Getenv("ENVIRONMENT")
	if Environment == "" {
		log.Error().Msg("You must specify a non-empty ENVIRONMENT")
		os.Exit(1)
	}
	if !(Environment == "PRODUCTION" || Environment == "DEVELOPMENT") {
		log.Error().Str("ENVIRONMENT", Environment).Msg("ENVIRONMENT should be either DEVELOPMENT or PRODUCTION")
		os.Exit(2)
	}
	MapboxToken = os.Getenv("MAPBOX_TOKEN")
	if MapboxToken == "" {
		log.Error().Msg("You must specify a non-empty MAPBOX_TOKEN")
		os.Exit(1)
	}
	pg_dsn := os.Getenv("POSTGRES_DSN")
	if pg_dsn == "" {
		log.Error().Msg("You must specify a non-empty POSTGRES_DSN")
		os.Exit(1)
	}

	log.Info().Msg("Starting...")
	err := db.InitializeDatabase(context.TODO(), pg_dsn)
	if err != nil {
		log.Error().Str("err", err.Error()).Msg("Failed to connect to database")
		os.Exit(2)
	}
	sessionManager = scs.New()
	sessionManager.Store = pgxstore.New(db.PGInstance.PGXPool)
	sessionManager.Lifetime = 24 * time.Hour

	router_logger := log.With().Logger()
	r := chi.NewRouter()
	r.Use(LoggerMiddleware(&router_logger))
	r.Use(sessionManager.LoadAndSave)

	// Root is a special endpoint that is neither authenticated nor unauthenticated
	r.Get("/", getRoot)

	// Unauthenticated endpoints
	r.Get("/arcgis/oauth/begin", getArcgisOauthBegin)
	r.Get("/arcgis/oauth/callback", getArcgisOauthCallback)
	r.Get("/favicon.ico", getFavicon)

	r.Get("/oauth/refresh", getOAuthRefresh)

	r.Get("/phone-call", getPhoneCall)
	r.Get("/qr-code/report/{code}", getQRCodeReport)
	r.Get("/report", getReport)
	r.Get("/report/{code}", getReportDetail)
	r.Get("/report/{code}/confirm", getReportConfirmation)
	r.Get("/report/{code}/contribute", getReportContribute)
	r.Get("/report/{code}/evidence", getReportEvidence)
	r.Get("/report/{code}/schedule", getReportSchedule)
	r.Get("/report/{code}/update", getReportUpdate)
	r.Get("/service-request", getServiceRequest)
	r.Get("/service-request/{code}", getServiceRequestDetail)
	r.Get("/service-request-location", getServiceRequestLocation)
	r.Get("/service-request-mosquito", getServiceRequestMosquito)
	r.Get("/service-request-pool", getServiceRequestPool)
	r.Get("/service-request-quick", getServiceRequestQuick)
	r.Get("/service-request-quick-confirmation", getServiceRequestQuickConfirmation)
	r.Get("/service-request-updates", getServiceRequestUpdates)
	r.Get("/signin", getSignin)
	r.Post("/signin", postSignin)
	r.Get("/signup", getSignup)
	r.Post("/signup", postSignup)

	// Authenticated endpoints
	r.Method("GET", "/cell/{cell}", NewEnsureAuth(getCellDetails))
	r.Method("GET", "/settings", NewEnsureAuth(getSettings))
	r.Method("GET", "/source/{globalid}", NewEnsureAuth(getSource))
	r.Method("GET", "/vector-tiles/{org_id}/{tileset_id}/{zoom}/{x}/{y}.{format}", NewEnsureAuth(getVectorTiles))

	localFS := http.Dir("./static")
	FileServer(r, "/static", localFS, embeddedStaticFS, "static")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	NewOAuthTokenChannel = make(chan struct{}, 10)

	var waitGroup sync.WaitGroup

	waitGroup.Add(1)
	go func() {
		defer waitGroup.Done()
		refreshFieldseekerData(ctx, NewOAuthTokenChannel)
	}()

	server := &http.Server{
		Addr:    bind,
		Handler: r,
	}
	go func() {
		log.Info().Str("address", bind).Msg("Serving HTTP requests")
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

func IsProductionEnvironment() bool {
	return Environment == "PRODUCTION"
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
					http.Error(ww, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
				}

				remote_addr := r.RemoteAddr
				forwarded_for := r.Header.Get("X-Forwarded-For")
				if forwarded_for != "" {
					remote_addr = forwarded_for
				}
				// log end request
				log.Info().
					Str("type", "access").
					Timestamp().
					Fields(map[string]interface{}{
						"remote_ip":  remote_addr,
						"url":        r.URL.Path,
						"proto":      r.Proto,
						"method":     r.Method,
						"user_agent": r.Header.Get("User-Agent"),
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
