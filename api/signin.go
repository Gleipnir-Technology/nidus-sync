package api

import (
	"context"
	"errors"
	"net/http"

	"github.com/Gleipnir-Technology/nidus-sync/auth"
	nhttp "github.com/Gleipnir-Technology/nidus-sync/http"
	"github.com/rs/zerolog/log"
)

type reqSignin struct {
	Password string `json:"password"`
	Username string `json:"username"`
}

func postSignin(ctx context.Context, r *http.Request, req reqSignin) (string, *nhttp.ErrorWithStatus) {
	if req.Password == "" {
		return "", nhttp.NewErrorStatus(http.StatusBadRequest, "Empty password")
	}
	if req.Username == "" {
		return "", nhttp.NewErrorStatus(http.StatusBadRequest, "Empty username")
	}
	log.Info().Str("username", req.Username).Msg("API Signin")
	_, err := auth.SigninUser(r, req.Username, req.Password)
	if err != nil {
		if errors.Is(err, auth.InvalidCredentials{}) {
			return "", nhttp.NewErrorStatus(http.StatusUnauthorized, "invalid credentials")
		}
		if errors.Is(err, auth.InvalidUsername{}) {
			return "", nhttp.NewErrorStatus(http.StatusUnauthorized, "invalid credentials")
		}
		log.Error().Err(err).Str("username", req.Username).Msg("Login server error")
		return "", nhttp.NewError("login server error")
	}
	return "/", nil
}
