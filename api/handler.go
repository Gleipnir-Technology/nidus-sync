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
	"github.com/Gleipnir-Technology/nidus-sync/platform/file"
	"github.com/Gleipnir-Technology/nidus-sync/resource"
	"github.com/google/uuid"
	"github.com/gorilla/schema"
	"github.com/rs/zerolog/log"
)

var decoder = schema.NewDecoder()

type handlerFunctionDelete func(context.Context, *http.Request, platform.User) *nhttp.ErrorWithStatus

func authenticatedHandlerDelete(f handlerFunctionDelete) http.Handler {
	return auth.NewEnsureAuth(func(w http.ResponseWriter, r *http.Request, u platform.User) {
		ctx := r.Context()
		e := f(ctx, r, u)
		if e != nil {
			log.Warn().Int("status", e.Status).Err(e).Str("user message", e.Message).Msg("Responding with an error from api")
			body, err := json.Marshal(ErrorAPI{Message: e.Error()})
			if err != nil {
				log.Error().Err(err).Msg("failed to marshal error")
				http.Error(w, "{\"message\": \"boom. I can't even tell you what went wrong\"}", http.StatusInternalServerError)
				return
			}
			http.Error(w, string(body), e.Status)
			return
		}
		http.Error(w, "", http.StatusNoContent)
		return
	})
}

type handlerFunctionGetImage func(context.Context, *http.Request, platform.User) (file.Collection, uuid.UUID, *nhttp.ErrorWithStatus)

func authenticatedHandlerGetImage(f handlerFunctionGetImage) http.Handler {
	return auth.NewEnsureAuth(func(w http.ResponseWriter, r *http.Request, u platform.User) {
		ctx := r.Context()
		collection, uid, e := f(ctx, r, u)
		if e != nil {
			log.Warn().Int("status", e.Status).Err(e).Str("user message", e.Message).Msg("Responding with an error from api")
			body, err := json.Marshal(ErrorAPI{Message: e.Error()})
			if err != nil {
				log.Error().Err(err).Msg("failed to marshal error")
				http.Error(w, "{\"message\": \"boom. I can't even tell you what went wrong\"}", http.StatusInternalServerError)
				return
			}
			http.Error(w, string(body), e.Status)
			return
		}
		file.ImageFileToWriter(collection, uid, w)
	})
}

type handlerFunctionGet[T any] func(context.Context, *http.Request, platform.User, resource.QueryParams) (*T, *nhttp.ErrorWithStatus)
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
		var params resource.QueryParams
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
			log.Error().Err(err).Msg("failed to marshal json")
			http.Error(w, "{\"message\": \"failed to marshal json\"}", http.StatusInternalServerError)
			return
		}
		w.Write(body)
	})
}

type handlerFunctionGetSlice[T any] func(context.Context, *http.Request, resource.QueryParams) ([]*T, *nhttp.ErrorWithStatus)
type handlerFunctionGetSliceAuthenticated[T any] func(context.Context, *http.Request, platform.User, resource.QueryParams) ([]*T, *nhttp.ErrorWithStatus)

func authenticatedHandlerJSONSlice[T any](f handlerFunctionGetSliceAuthenticated[T]) http.Handler {
	return auth.NewEnsureAuth(func(w http.ResponseWriter, r *http.Request, u platform.User) {
		ctx := r.Context()
		var body []byte
		var params resource.QueryParams
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
			log.Error().Err(err).Msg("failed to marshal json")
			http.Error(w, "{\"message\": \"failed to marshal json\"}", http.StatusInternalServerError)
			return
		}
		w.Write(body)
	})
}
func handlerJSONSlice[T any](f handlerFunctionGetSlice[T]) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		var body []byte
		var params resource.QueryParams
		err := decoder.Decode(&params, r.URL.Query())
		if err != nil {
			log.Error().Err(err).Msg("decode query failure")
			http.Error(w, "failed to decode query", http.StatusInternalServerError)
			return
		}
		resp, e := f(ctx, r, params)
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
			log.Error().Err(err).Msg("failed to marshal json")
			http.Error(w, "{\"message\": \"failed to marshal json\"}", http.StatusInternalServerError)
			return
		}
		w.Write(body)
	}
}

type handlerFunctionPost[ReqType any] func(context.Context, *http.Request, ReqType) (string, *nhttp.ErrorWithStatus)
type handlerFunctionPostAuthenticated[RequestType any, ResponseType any] func(context.Context, *http.Request, platform.User, RequestType) (ResponseType, *nhttp.ErrorWithStatus)

func authenticatedHandlerJSONPost[RequestType any, ResponseType any](f handlerFunctionPostAuthenticated[RequestType, ResponseType]) http.Handler {
	return auth.NewEnsureAuth(func(w http.ResponseWriter, r *http.Request, u platform.User) {
		w.Header().Set("Content-Type", "application/json")
		req, e := parseRequest[RequestType](r)
		if e != nil {
			serializeError(w, e)
			return
		}
		ctx := r.Context()
		resp, e := f(ctx, r, u, *req)
		if e != nil {
			serializeError(w, e)
			return
		}
		body, err := json.Marshal(resp)
		if err != nil {
			log.Error().Err(err).Msg("failed to marshal json")
			http.Error(w, "{\"message\": \"failed to marshal json\"}", http.StatusInternalServerError)
			return
		}
		w.Write(body)
	})
}

type handlerFunctionPutAuthenticated[ReqType any] func(context.Context, *http.Request, platform.User, ReqType) (string, *nhttp.ErrorWithStatus)

func authenticatedHandlerJSONPut[ReqType any](f handlerFunctionPutAuthenticated[ReqType]) http.Handler {
	return auth.NewEnsureAuth(func(w http.ResponseWriter, r *http.Request, u platform.User) {
		w.Header().Set("Content-Type", "application/json")
		req, e := parseRequest[ReqType](r)
		if e != nil {
			serializeError(w, e)
			return
		}
		ctx := r.Context()
		path, e := f(ctx, r, u, *req)
		if e != nil {
			serializeError(w, e)
			return
		}
		if path == "" {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		w.Header().Set("Location", path)
		http.Redirect(w, r, path, http.StatusCreated)
	})
}
func parseRequest[ReqType any](r *http.Request) (*ReqType, *nhttp.ErrorWithStatus) {
	var req ReqType
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return nil, nhttp.NewError("Failed to read body: %w", err)
	}
	err = json.Unmarshal(body, &req)
	if err != nil {
		return nil, nhttp.NewErrorStatus(http.StatusBadRequest, "Failed to decode request: %w", err)
	}
	return &req, nil
}
func serializeError(w http.ResponseWriter, e *nhttp.ErrorWithStatus) {
	log.Warn().Int("status", e.Status).Err(e).Str("user message", e.Message).Msg("Responding with an error from api")
	body, err := json.Marshal(ErrorAPI{Message: e.Error()})
	if err != nil {
		log.Error().Err(err).Msg("failed to marshal error")
		http.Error(w, "{\"message\": \"boom. I can't even tell you what went wrong\"}", http.StatusInternalServerError)
		return
	}
	http.Error(w, string(body), e.Status)
	return
}
func handlerJSONPost[ReqType any](f handlerFunctionPost[ReqType]) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		req, e := parseRequest[ReqType](r)
		if e != nil {
			serializeError(w, e)
			return
		}
		ctx := r.Context()
		path, e := f(ctx, r, *req)
		if e != nil {
			return
		}
		http.Redirect(w, r, path, http.StatusFound)
	}
}

type postMultipartResponse struct {
	URI string `json:"uri"`
}

func authenticatedHandlerPostMultipart[ResponseType any](f handlerFunctionPostAuthenticated[[]file.Upload, ResponseType], collection file.Collection) http.Handler {
	return auth.NewEnsureAuth(func(w http.ResponseWriter, r *http.Request, u platform.User) {
		err := r.ParseMultipartForm(32 << 10) // 32 MB buffer
		if err != nil {
			respondError(w, http.StatusBadRequest, "Failed to parse form: %w ", err)
			return
		}
		uploads, err := file.SaveFileUploads(r, collection)
		if err != nil {
			respondError(w, http.StatusInternalServerError, "failed to save uploads: %w", err)
			return
		}

		/*
			err = decoder.Decode(&content, r.PostForm)
			if err != nil {
				respondError(w, http.StatusBadRequest, "Failed to decode form: %w", err)
				return
			}
		*/
		ctx := r.Context()
		resp, e := f(ctx, r, u, uploads)
		if e != nil {
			http.Error(w, e.Error(), e.Status)
			return
		}
		body, err := json.Marshal(resp)
		if err != nil {
			log.Error().Err(err).Msg("failed to marshal json")
			http.Error(w, "{\"message\": \"failed to marshal json\"}", http.StatusInternalServerError)
			return
		}
		w.Write(body)
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
