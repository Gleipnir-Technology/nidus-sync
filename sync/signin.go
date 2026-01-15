package sync

import (
	"errors"
	"net/http"
	"strings"

	"github.com/Gleipnir-Technology/nidus-sync/auth"
	"github.com/Gleipnir-Technology/nidus-sync/db/models"
	"github.com/Gleipnir-Technology/nidus-sync/htmlpage"
	"github.com/rs/zerolog/log"
)

var (
	signinT = buildTemplate("signin", "base")
	signupT = buildTemplate("signup", "base")
)

func getSignin(w http.ResponseWriter, r *http.Request) {
	errorCode := r.URL.Query().Get("error")
	signin(w, errorCode)
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

	username := r.FormValue("username")
	password := r.FormValue("password")

	log.Info().Str("username", username).Msg("Signin")

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

	http.Redirect(w, r, "/", http.StatusFound)
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

func signin(w http.ResponseWriter, errorCode string) {
	data := ContentSignin{
		InvalidCredentials: errorCode == "invalid-credentials",
	}
	htmlpage.RenderOrError(w, signinT, data)
}

func signup(w http.ResponseWriter, path string) {
	data := ContentSignup{}
	htmlpage.RenderOrError(w, signupT, data)
}
