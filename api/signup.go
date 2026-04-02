package api

import (
	"context"
	"net/http"
	"strings"

	"github.com/Gleipnir-Technology/nidus-sync/auth"
	nhttp "github.com/Gleipnir-Technology/nidus-sync/http"
	"github.com/rs/zerolog/log"
)

type reqSignup struct {
	Username string `json:"username"`
	Name     string `json:"name"`
	Password string `json:"password"`
	Terms    bool   `json:"terms"`
}

func postSignup(ctx context.Context, r *http.Request, signup reqSignup) (string, *nhttp.ErrorWithStatus) {

	log.Info().Str("username", signup.Username).Str("name", signup.Name).Str("password", strings.Repeat("*", len(signup.Password))).Msg("Signup")

	if !signup.Terms {
		log.Warn().Msg("Terms not agreed")
		return "", nhttp.NewErrorStatus(http.StatusBadRequest, "You must agree to the terms to register")
	}

	user, err := auth.SignupUser(r.Context(), signup.Username, signup.Name, signup.Password)
	if err != nil {
		return "", nhttp.NewError("Failed to signup user", err)
	}

	auth.AddUserSession(ctx, user)

	return "/", nil
}
