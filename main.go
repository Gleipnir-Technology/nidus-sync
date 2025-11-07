package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/alexedwards/scs/pgxstore"
	"github.com/alexedwards/scs/v2"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

var sessionManager *scs.SessionManager

var BaseURL, ClientID, ClientSecret string

func main() {
	ClientID = os.Getenv("ARCGIS_CLIENT_ID")
	if ClientID == "" {
		slog.Error("You must specify a non-empty ARCGIS_CLIENT_ID")
		os.Exit(1)
	}
	ClientSecret = os.Getenv("ARCGIS_CLIENT_SECRET")
	if ClientSecret == "" {
		slog.Error("You must specify a non-empty ARCGIS_CLIENT_SECRET")
		os.Exit(1)
	}
	BaseURL = os.Getenv("BASE_URL")
	if BaseURL == "" {
		slog.Error("You must specify a non-empty BASE_URL")
		os.Exit(1)
	}
	bind := os.Getenv("BIND")
	if bind == "" {
		bind = ":9001"
	}
	pg_dsn := os.Getenv("POSTGRES_DSN")
	if pg_dsn == "" {
		slog.Error("You must specify a non-empty POSTGRES_DSN")
		os.Exit(1)
	}

	slog.Info("Starting...")
	err := initializeDatabase(context.TODO(), pg_dsn)
	if err != nil {
		slog.Error("Failed to connect to database", slog.String("err", err.Error()))
		os.Exit(2)
	}
	sessionManager = scs.New()
	sessionManager.Store = pgxstore.New(PGInstance.PGXPool)
	sessionManager.Lifetime = 24 * time.Hour

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(sessionManager.LoadAndSave)

	r.Get("/", getRoot)
	r.Get("/arcgis/oauth/begin", getArcgisOauthBegin)
	r.Get("/arcgis/oauth/callback", getArcgisOauthCallback)
	r.Get("/qr-code/report/{code}", getQRCodeReport)
	r.Get("/report", getReport)
	r.Get("/report/{code}", getReportDetail)
	r.Get("/report/{code}/confirm", getReportConfirmation)
	r.Get("/report/{code}/contribute", getReportContribute)
	r.Get("/report/{code}/evidence", getReportEvidence)
	r.Get("/report/{code}/schedule", getReportSchedule)
	r.Get("/report/{code}/update", getReportUpdate)
	r.Post("/signin", postSignin)
	r.Get("/signup", getSignup)
	r.Post("/signup", postSignup)
	r.Get("/favicon.ico", getFavicon)

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
		slog.Info("Serving HTTP requests", slog.String("address", bind))
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("HTTP Server Error", slog.String("err", err.Error()))
		}
	}()

	// Wait for the interrupt signal to gracefully shut down
	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, syscall.SIGINT, syscall.SIGTERM)
	<-signalCh

	slog.Info("Received shutdown signal, shutting down...")
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdownCancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		slog.Error("HTTP server shutdown error", slog.String("err", err.Error()))
	}

	cancel()

	waitGroup.Wait()

	slog.Info("Shutdown complete")
}
