package api

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Gleipnir-Technology/nidus-sync/auth"
	"github.com/Gleipnir-Technology/nidus-sync/db"
	"github.com/Gleipnir-Technology/nidus-sync/db/models"
	"github.com/Gleipnir-Technology/nidus-sync/html"
	nhttp "github.com/Gleipnir-Technology/nidus-sync/http"
	"github.com/Gleipnir-Technology/nidus-sync/platform"
	"github.com/gorilla/schema"
	"github.com/rs/zerolog/log"
)

var decoder = schema.NewDecoder()

type handlerFunctionGet[T any] func(context.Context, *http.Request, *models.Organization, *models.User) (*T, *nhttp.ErrorWithStatus)
type wrappedHandler func(http.ResponseWriter, *http.Request)
type contentAuthenticated[T any] struct {
	C      T
	Config html.ContentConfig
	User   platform.User
}

type ErrorAPI struct {
	Message string `json:"message"`
}

func authenticatedHandlerJSON[T any](f handlerFunctionGet[T]) http.Handler {
	return auth.NewEnsureAuth(func(w http.ResponseWriter, r *http.Request, u *models.User) {
		ctx := r.Context()
		org, err := u.Organization().One(ctx, db.PGInstance.BobDB)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if org == nil {
			http.Error(w, "nil org", http.StatusInternalServerError)
			return
		}
		var body []byte
		resp, e := f(ctx, r, org, u)
		w.Header().Set("Content-Type", "application/json")
		//log.Info().Str("template", template).Err(e).Msg("handler done")
		if e != nil {
			log.Warn().Int("status", e.Status).Err(e).Str("user message", e.Message).Msg("Responding with an error from sync pages")
			body, err = json.Marshal(ErrorAPI{Message: e.Error()})
			if err != nil {
				log.Error().Err(err).Msg("failed to marshal error")
				http.Error(w, "{\"message\": \"boom. I can't even tell you what went wrong\"}", http.StatusInternalServerError)
				return
			}
			http.Error(w, string(body), e.Status)
			return
		}
		body, err = json.Marshal(resp)
		if err != nil {
			http.Error(w, "{\"message\": \"failed to marshal json\"}", http.StatusInternalServerError)
			return
		}
		w.Write(body)
	})
}

type handlerFunctionPost[FormType any, ResponseType any] func(context.Context, *http.Request, *models.Organization, *models.User, FormType) (ResponseType, *nhttp.ErrorWithStatus)

func authenticatedHandlerJSONPost[FormType any, ResponseType any](f handlerFunctionPost[FormType, ResponseType]) http.Handler {
	return auth.NewEnsureAuth(func(w http.ResponseWriter, r *http.Request, u *models.User) {
		w.Header().Set("Content-Type", "application/json")
		err := r.ParseForm()
		if err != nil {
			respondError(w, http.StatusBadRequest, "failed to parse form: %w", err)
			return
		}

		var form FormType
		err = decoder.Decode(&form, r.PostForm)
		if err != nil {
			respondError(w, http.StatusBadRequest, "Failed to decode form: %w", err)
			return
		}
		ctx := r.Context()
		org, err := u.Organization().One(ctx, db.PGInstance.BobDB)
		if err != nil {
			respondError(w, http.StatusInternalServerError, "Failed to get org: %w", err)
			return
		}
		response, e := f(ctx, r, org, u, form)
		if e != nil {
			http.Error(w, e.Error(), e.Status)
			return
		}
		resp_body, err := json.Marshal(response)
		if err != nil {
			respondError(w, http.StatusInternalServerError, "Failed to marshal json response: %w", err)
			return
		}
		w.Write(resp_body)
	})
}

func respondError(w http.ResponseWriter, status int, format string, args ...any) {
	outer_err := fmt.Errorf(format, args...)
	body, err := json.Marshal(ErrorAPI{
		Message: outer_err.Error(),
	})
	if err != nil {
		http.Error(w, "{\"message\": \"failed to marshal json\"}", http.StatusInternalServerError)
		return
	}
	http.Error(w, string(body), status)
}
