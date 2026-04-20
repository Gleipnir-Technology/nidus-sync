package main

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"runtime/debug"
	"syscall"
	"time"

	"github.com/Gleipnir-Technology/nidus-sync/api"
	"github.com/Gleipnir-Technology/nidus-sync/auth"
	"github.com/Gleipnir-Technology/nidus-sync/config"
	"github.com/Gleipnir-Technology/nidus-sync/db"
	"github.com/Gleipnir-Technology/nidus-sync/html"
	"github.com/Gleipnir-Technology/nidus-sync/llm"
	"github.com/Gleipnir-Technology/nidus-sync/middleware"
	"github.com/Gleipnir-Technology/nidus-sync/platform"
	"github.com/Gleipnir-Technology/nidus-sync/rmo"
	nidussync "github.com/Gleipnir-Technology/nidus-sync/sync"
	"github.com/getsentry/sentry-go"
	sentryhttp "github.com/getsentry/sentry-go/http"
	"github.com/getsentry/sentry-go/zerolog"
	"github.com/gorilla/mux"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	err := config.Parse()
	if err != nil {
		log.Error().Err(err).Msg("Failed to parse config")
		os.Exit(1)
	}

	var prod = flag.Bool("prod", false, "Force into production mode")
	flag.Parse()
	if prod != nil && *prod {
		log.Warn().Msg("Forcing production mode for testing templates")
		config.Environment = "PRODUCTION"
	}
	log.Info().Str("environment", config.Environment).Bool("is-prod", config.IsProductionEnvironment()).Str("version", Version).Str("commit", Commit).Msg("Starting")
	err = sentry.Init(sentry.ClientOptions{
		Debug:            false, //!config.IsProductionEnvironment(),
		Dsn:              config.SentryDSN,
		EnableTracing:    true,
		SendDefaultPII:   true,
		TracesSampleRate: 1.0,
	})
	if err != nil {
		log.Error().Err(err).Msg("Failed to start sentry connection")
		os.Exit(2)
	}
	sentryWriter, err := sentryzerolog.New(sentryzerolog.Config{
		ClientOptions: sentry.ClientOptions{
			Dsn: config.SentryDSN,
		},
		Options: sentryzerolog.Options{
			Levels:          []zerolog.Level{zerolog.ErrorLevel, zerolog.FatalLevel, zerolog.PanicLevel},
			WithBreadcrumbs: true,
			FlushTimeout:    3 * time.Second,
		},
	})
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to create sentry writer")
		os.Exit(2)
	}
	defer sentryWriter.Close()

	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	log.Logger = log.Output(zerolog.MultiLevelWriter(zerolog.ConsoleWriter{Out: os.Stderr}, sentryWriter))
	if os.Getenv("VERBOSE") != "" {
		log.Logger = log.Logger.Level(zerolog.DebugLevel)
	} else {
		log.Logger = log.Logger.Level(zerolog.InfoLevel)
	}

	// Defer cleanup in reverse order - these will execute LAST (LIFO)
	defer func() {
		log.Info().Msg("Final cleanup")
		os.Stderr.Sync()
		sentryWriter.Close()
		sentry.Flush(2 * time.Second)
	}()

	err = db.InitializeDatabase(context.TODO(), config.PGDSN)
	if err != nil {
		log.Error().Err(err).Msg("Failed to connect to database")
		os.Exit(3)
	}

	err = html.LoadTemplates()
	if err != nil {
		log.Error().Err(err).Msg("Failed to load html templates")
		os.Exit(4)
	}
	// Start up background processes
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	err = platform.StartAll(ctx)
	if err != nil {
		log.Error().Err(err).Msg("Failed at platform.StartAll")
		os.Exit(5)
	}
	router_logger := log.With().Logger()
	sentryMiddleware := sentryhttp.New(sentryhttp.Options{
		Repanic: true,
	})
	//r := chi.NewRouter()
	r := mux.NewRouter()

	r.Use(LoggerMiddleware(&router_logger))
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	//r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(sentryMiddleware.Handle)
	r.Use(auth.NewSessionManager().LoadAndSave)

	sync_router := r.Host(config.DomainNidus).Subrouter()
	rmo_router := r.Host(config.DomainRMO).Subrouter()

	// Set up routing by hostname
	sync_api_router := sync_router.PathPrefix("/api").Subrouter()
	api.AddRoutes(sync_api_router)
	nidussync.Router(sync_router)

	rmo_api_router := rmo_router.PathPrefix("/api").Subrouter()
	api.AddRoutes(rmo_api_router)
	rmo.Router(rmo_router)

	//hr.Map("", sr)                         // default
	//hr.Map("*", sr)                        // default

	log.Debug().Str("report url", config.DomainRMO).Str("sync url", config.DomainNidus).Msg("Serving at URLs")

	openai_logger := log.With().Logger()
	err = llm.CreateOpenAIClient(ctx, &openai_logger)
	if err != nil {
		log.Error().Err(err).Msg("Failed to start openAI client")
		os.Exit(8)
	}
	server := &http.Server{
		Addr:    config.Bind,
		Handler: r,
	}
	go func() {
		log.Info().Str("address", config.Bind).Msg("Serving HTTP requests")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Error().Str("err", err.Error()).Msg("HTTP Server Error")
		}
		log.Debug().Msg("Exiting listen-and-serve goroutine")
	}()

	chan_envelope := make(chan platform.Envelope, 10)
	api.SetVersion(Version)
	platform.SetEventChannel(chan_envelope)
	api.SetEventChannel(chan_envelope)

	// Wait for the interrupt signal to gracefully shut down
	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, syscall.SIGINT, syscall.SIGTERM)
	<-signalCh
	log.Info().Msg("Received shutdown signal, shutting down...")
	// Ensure logs are flushed
	os.Stderr.Sync()

	platform.EventShutdown()
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdownCancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Error().Err(err).Msg("Error on HTTP server shutdown")
	}

	cancel()
	close(chan_envelope)
	platform.WaitForExit()

	log.Info().Msg("Shutdown complete")
	os.Stderr.Sync()
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
				if os.Getenv("VERBOSE") != "" {
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
				}
			}()

			next.ServeHTTP(ww, r)
		}
		return http.HandlerFunc(fn)
	}
}
