package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"runtime/debug"
	"syscall"
	"time"

	"github.com/Gleipnir-Technology/nidus-sync/auth"
	"github.com/Gleipnir-Technology/nidus-sync/background"
	"github.com/Gleipnir-Technology/nidus-sync/comms/email"
	"github.com/Gleipnir-Technology/nidus-sync/comms/text"
	"github.com/Gleipnir-Technology/nidus-sync/config"
	"github.com/Gleipnir-Technology/nidus-sync/db"
	"github.com/Gleipnir-Technology/nidus-sync/llm"
	"github.com/Gleipnir-Technology/nidus-sync/public-report"
	nidussync "github.com/Gleipnir-Technology/nidus-sync/sync"
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

	err = email.LoadTemplates()
	if err != nil {
		log.Error().Err(err).Msg("Failed to load email templates")
		os.Exit(3)
	}

	err = text.StoreSources()
	if err != nil {
		log.Error().Err(err).Msg("Failed to store text source phone numbers")
		os.Exit(4)
	}

	router_logger := log.With().Logger()
	r := chi.NewRouter()

	r.Use(LoggerMiddleware(&router_logger))
	r.Use(middleware.RealIP)
	r.Use(auth.NewSessionManager().LoadAndSave)

	hr := hostrouter.New()

	// Set up routing by hostname
	sr := nidussync.Router()
	hr.Map("", sr)                                  // default
	hr.Map("*", sr)                                 // default
	hr.Map(config.DomainRMO, publicreport.Router()) // report.mosquitoes.online
	hr.Map(config.DomainNidus, sr)
	r.Mount("/", hr)

	log.Info().Str("report url", config.DomainRMO).Str("sync url", config.DomainNidus).Msg("Serving at URLs")

	// Start up background processes
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	err = llm.CreateOpenAIClient(ctx)
	if err != nil {
		log.Error().Err(err).Msg("Failed to start openAI client")
		os.Exit(5)
	}
	background.Start(ctx)
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
	background.WaitForExit()

	log.Info().Msg("Shutdown complete")
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
