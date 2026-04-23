package api

import (
	"context"
	"errors"
	"net/http"

	"github.com/Gleipnir-Technology/nidus-sync/auth"
	nhttp "github.com/Gleipnir-Technology/nidus-sync/http"
	"github.com/Gleipnir-Technology/nidus-sync/platform"
	"github.com/rs/zerolog/log"
)

type reqSignin struct {
	Password string `schema:"password"`
	Username string `schema:"username"`
}

func postSignin(ctx context.Context, r *http.Request, req reqSignin) (string, *nhttp.ErrorWithStatus) {
	if req.Password == "" {
		return "", nhttp.NewBadRequest("Empty password")
	}
	if req.Username == "" {
		return "", nhttp.NewBadRequest("Empty username")
	}
	log.Info().Str("username", req.Username).Msg("API Signin")
	_, err := auth.SigninUser(r, req.Username, req.Password)
	if err != nil {
		if errors.Is(err, auth.InvalidCredentials{}) {
			return "", nhttp.NewUnauthorized("invalid credentials")
		}
		if errors.Is(err, auth.InvalidUsername{}) {
			return "", nhttp.NewUnauthorized("invalid credentials")
		}
		if errors.Is(err, platform.NoUserError{}) {
			return "", nhttp.NewUnauthorized("invalid credentials")
		}
		log.Error().Err(err).Str("username", req.Username).Msg("Login server error")
		return "", nhttp.NewError("login server error")
	}
	return "/", nil
}
func postSignout(ctx context.Context, w http.ResponseWriter, r *http.Request, u platform.User) *nhttp.ErrorWithStatus {
	auth.SignoutUser(r, u)
	return nil
}
