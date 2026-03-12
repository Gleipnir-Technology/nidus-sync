package api

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/Gleipnir-Technology/nidus-sync/auth"
	"github.com/Gleipnir-Technology/nidus-sync/html"
	nhttp "github.com/Gleipnir-Technology/nidus-sync/http"
	"github.com/Gleipnir-Technology/nidus-sync/platform"
	"github.com/gorilla/schema"
	"github.com/rs/zerolog/log"
)

var decoder = schema.NewDecoder()

type handlerFunctionGet[T any] func(context.Context, *http.Request, platform.User, queryParams) (*T, *nhttp.ErrorWithStatus)
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
	return auth.NewEnsureAuth(func(w http.ResponseWriter, r *http.Request, u platform.User) {
		ctx := r.Context()
		var body []byte
		var params queryParams
		err := decoder.Decode(&params, r.URL.Query())
		if err != nil {
			log.Error().Err(err).Msg("decode query failure")
			http.Error(w, "failed to decode query", http.StatusInternalServerError)
			return
		}
		resp, e := f(ctx, r, u, params)
		w.Header().Set("Content-Type", "application/json")
		//log.Info().Str("template", template).Err(e).Msg("handler done")
		if e != nil {
			log.Warn().Int("status", e.Status).Err(e).Str("user message", e.Message).Msg("Responding with an error from api")
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

type handlerFunctionPost[ReqType any, ResponseType any] func(context.Context, *http.Request, platform.User, ReqType) (ResponseType, *nhttp.ErrorWithStatus)

func authenticatedHandlerJSONPost[ReqType any, ResponseType any](f handlerFunctionPost[ReqType, ResponseType]) http.Handler {
	return auth.NewEnsureAuth(func(w http.ResponseWriter, r *http.Request, u platform.User) {
		w.Header().Set("Content-Type", "application/json")
		var req ReqType
		body, err := io.ReadAll(r.Body)
		if err != nil {
			respondError(w, http.StatusInternalServerError, "Failed to read body: %w", err)
			return
		}
		err = json.Unmarshal(body, &req)
		if err != nil {
			respondError(w, http.StatusBadRequest, "Failed to decode request: %w", err)
			return
		}
		ctx := r.Context()
		response, e := f(ctx, r, u, req)
		if e != nil {
			log.Warn().Int("status", e.Status).Err(e).Str("user message", e.Message).Msg("Responding with an error from api")
			body, err = json.Marshal(ErrorAPI{Message: e.Error()})
			if err != nil {
				log.Error().Err(err).Msg("failed to marshal error")
				http.Error(w, "{\"message\": \"boom. I can't even tell you what went wrong\"}", http.StatusInternalServerError)
				return
			}
			http.Error(w, string(body), e.Status)
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
