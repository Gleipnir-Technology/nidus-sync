package api

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/Gleipnir-Technology/nidus-sync/auth"
	"github.com/go-chi/render"
	"github.com/rs/zerolog/log"
)

func postSignin(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		render.Render(w, r, errRender(fmt.Errorf("Failed to parse POST form: %w", err)))
		return
	}

	username := r.FormValue("username")
	password := r.FormValue("password")

	if password == "" || username == "" {
		w.Header().Set("WWW-Authenticate-Error", "no-credentials")
		http.Error(w, "invalid-credentials", http.StatusUnauthorized)
		return
	}
	log.Info().Str("username", username).Msg("API Signin")
	_, err := auth.SigninUser(r, username, password)
	if err != nil {
		if errors.Is(err, auth.InvalidCredentials{}) {
			w.Header().Set("WWW-Authenticate-Error", "invalid-credentials")
			http.Error(w, "invalid-credentials", http.StatusUnauthorized)
			return
		}
		if errors.Is(err, auth.InvalidUsername{}) {
			w.Header().Set("WWW-Authenticate-Error", "invalid-credentials")
			http.Error(w, "invalid-credentials", http.StatusUnauthorized)
			return
		}
		log.Error().Err(err).Str("username", username).Msg("Login server error")
		http.Error(w, "signin-server-error", http.StatusInternalServerError)
		return
	}

	http.Error(w, "", http.StatusAccepted)
}
