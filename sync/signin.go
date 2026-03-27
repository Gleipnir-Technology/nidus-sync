package sync

import (
	"errors"
	"net/http"

	"github.com/Gleipnir-Technology/nidus-sync/auth"
	"github.com/Gleipnir-Technology/nidus-sync/config"
	"github.com/Gleipnir-Technology/nidus-sync/html"
	"github.com/Gleipnir-Technology/nidus-sync/platform"
	"github.com/rs/zerolog/log"
)

type contentSignin struct {
	InvalidCredentials bool
	Next               string
}

func getSignin(w http.ResponseWriter, r *http.Request) {
	errorCode := r.URL.Query().Get("error")
	next := r.URL.Query().Get("next")
	signin(w, errorCode, next)
}

func getSignout(w http.ResponseWriter, r *http.Request, user platform.User) {
	auth.SignoutUser(r, user)
	http.Redirect(w, r, "/signin", http.StatusFound)
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

type contentUnauthenticated[T any] struct {
	C      T
	Config html.ContentConfig
	URL    html.ContentURL
}

func signin(w http.ResponseWriter, errorCode string, next string) {
	if next == "" {
		next = "/"
	}
	data := contentUnauthenticated[contentSignin]{
		C: contentSignin{
			InvalidCredentials: errorCode == "invalid-credentials",
			Next:               next,
		},
		Config: html.NewContentConfig(),
		URL:    html.NewContentURL(),
	}
	html.RenderOrError(w, "sync/signin.html", data)
}
