package main

import (
	"context"
	"log"
	"net/http"
	"os"
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
		log.Println("You must specify a non-empty ARCGIS_CLIENT_ID")
		os.Exit(1)
	}
	ClientSecret = os.Getenv("ARCGIS_CLIENT_SECRET")
	if ClientSecret == "" {
		log.Println("You must specify a non-empty ARCGIS_CLIENT_SECRET")
		os.Exit(1)
	}
	BaseURL = os.Getenv("BASE_URL")
	if BaseURL == "" {
		log.Println("You must specify a non-empty BASE_URL")
		os.Exit(1)
	}
	bind := os.Getenv("BIND")
	if bind == "" {
		bind = ":9001"
	}
	pg_dsn := os.Getenv("POSTGRES_DSN")
	if pg_dsn == "" {
		log.Println("You must specify a non-empty POSTGRES_DSN")
		os.Exit(1)
	}

	log.Println("Starting...")
	err := initializeDatabase(context.TODO(), pg_dsn)
	if err != nil {
		log.Printf("Failed to connect to database: %v", err)
		os.Exit(2)
	}
	sessionManager = scs.New()
	sessionManager.Store = pgxstore.New(PGInstance.PGXPool)
	sessionManager.Lifetime = 24 * time.Hour

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(sessionManager.LoadAndSave)

	r.Get("/", getRoot)
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

	log.Printf("Serving on %s", bind)
	log.Fatal(http.ListenAndServe(bind, r))
}
