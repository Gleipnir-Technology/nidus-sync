package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

var sessionManager *scs.SessionManager

var BaseURL, ClientID, ClientSecret string

func main() {
	ClientID = os.Getenv("ARCGIS_CLIENT_ID")
	if ClientID == "" {
		log.Println("You must specify a non-empty CLIENT_ID")
		os.Exit(1)
	}
	ClientSecret = os.Getenv("ARCGIS_CLIENT_SECRET")
	if ClientSecret == "" {
		log.Println("You must specify a non-empty CLIENT_SECRET")
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

	log.Println("Starting...")
	sessionManager = scs.New()
	sessionManager.Lifetime = 24 * time.Hour

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(sessionManager.LoadAndSave)

	r.Get("/", getRoot)
	r.Get("/signup", getSignup)
	r.Get("/favicon.ico", getFavicon)

	localFS := http.Dir("./static")
	FileServer(r, "/static", localFS, embeddedStaticFS, "static")

	log.Printf("Serving on %s", bind)
	log.Fatal(http.ListenAndServe(bind, r))
}
