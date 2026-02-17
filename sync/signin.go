package sync

import (
	"errors"
	"net/http"
	"strings"

	"github.com/Gleipnir-Technology/nidus-sync/auth"
	"github.com/Gleipnir-Technology/nidus-sync/config"
	"github.com/Gleipnir-Technology/nidus-sync/db/models"
	"github.com/Gleipnir-Technology/nidus-sync/html"
	"github.com/rs/zerolog/log"
)

type ContentSignin struct {
	InvalidCredentials bool
	Next               string
}
type ContentSignup struct{}

func getSignin(w http.ResponseWriter, r *http.Request) {
	errorCode := r.URL.Query().Get("error")
	next := r.URL.Query().Get("next")
	signin(w, errorCode, next)
}

func getSignout(w http.ResponseWriter, r *http.Request, user *models.User) {
	auth.SignoutUser(r, user)
	http.Redirect(w, r, "/signin", http.StatusFound)
}

func getSignup(w http.ResponseWriter, r *http.Request) {
	signup(w, r.URL.Path)
}

func postSignin(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		respondError(w, "Could not parse form", err, http.StatusBadRequest)
		return
	}

	next := r.FormValue("next")
	username := r.FormValue("username")
	password := r.FormValue("password")

	log.Info().Str("username", username).Str("next", next).Msg("HTML Signin")

	_, err := auth.SigninUser(r, username, password)
	if err != nil {
		if errors.Is(err, auth.InvalidCredentials{}) {
			http.Redirect(w, r, "/signin?error=invalid-credentials", http.StatusFound)
			return
		}
		if errors.Is(err, auth.InvalidUsername{}) {
			http.Redirect(w, r, "/signin?error=invalid-credentials", http.StatusFound)
			return
		}
		respondError(w, "Failed to signin user", err, http.StatusInternalServerError)
		return
	}
	if next == "" {
		next = "/"
	}
	location := config.MakeURLNidus(next)
	http.Redirect(w, r, location, http.StatusFound)
}

func postSignup(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		respondError(w, "Could not parse form", err, http.StatusBadRequest)
		return
	}

	username := r.FormValue("username")
	name := r.FormValue("name")
	password := r.FormValue("password")
	terms := r.FormValue("terms")

	log.Info().Str("username", username).Str("name", name).Str("password", strings.Repeat("*", len(password))).Msg("Signup")

	if terms != "on" {
		log.Warn().Msg("Terms not agreed")
		http.Error(w, "You must agree to the terms to register", http.StatusBadRequest)
		return
	}

	user, err := auth.SignupUser(r.Context(), username, name, password)
	if err != nil {
		respondError(w, "Failed to signup user", err, http.StatusInternalServerError)
		return
	}

	auth.AddUserSession(r, user)

	http.Redirect(w, r, "/", http.StatusFound)
}

func signin(w http.ResponseWriter, errorCode string, next string) {
	if next == "" {
		next = "/"
	}
	data := ContentSignin{
		InvalidCredentials: errorCode == "invalid-credentials",
		Next:               next,
	}
	html.RenderOrError(w, "sync/signin.html", data)
}

func signup(w http.ResponseWriter, path string) {
	data := ContentSignup{}
	html.RenderOrError(w, "sync/signup.html", data)
}
