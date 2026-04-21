package api

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/Gleipnir-Technology/nidus-sync/auth"
	nhttp "github.com/Gleipnir-Technology/nidus-sync/http"
	"github.com/Gleipnir-Technology/nidus-sync/platform"
	"github.com/Gleipnir-Technology/nidus-sync/platform/file"
	"github.com/Gleipnir-Technology/nidus-sync/resource"
	"github.com/google/uuid"
	"github.com/gorilla/schema"
	"github.com/rs/zerolog/log"
)

type ErrorAPI struct {
	Message string `json:"message"`
}

var decoder = schema.NewDecoder()

type handlerBase func(context.Context, http.ResponseWriter, *http.Request) *nhttp.ErrorWithStatus
type handlerBaseAuthenticated func(context.Context, http.ResponseWriter, *http.Request, platform.User) *nhttp.ErrorWithStatus
type handlerFunctionDelete func(context.Context, *http.Request, platform.User) *nhttp.ErrorWithStatus
type handlerFunctionGet[T any] func(context.Context, *http.Request, resource.QueryParams) (*T, *nhttp.ErrorWithStatus)
type handlerFunctionGetAuthenticated[T any] func(context.Context, *http.Request, platform.User, resource.QueryParams) (*T, *nhttp.ErrorWithStatus)
type handlerFunctionGetImage func(context.Context, *http.Request, platform.User) (file.Collection, uuid.UUID, *nhttp.ErrorWithStatus)
type handlerFunctionGetSlice[T any] func(context.Context, *http.Request, resource.QueryParams) ([]*T, *nhttp.ErrorWithStatus)
type handlerFunctionGetSliceAuthenticated[T any] func(context.Context, *http.Request, platform.User, resource.QueryParams) ([]*T, *nhttp.ErrorWithStatus)
type handlerFunctionPost[RequestType any, ResponseType any] func(context.Context, *http.Request, RequestType) (ResponseType, *nhttp.ErrorWithStatus)
type handlerFunctionPostAuthenticated[RequestType any, ResponseType any] func(context.Context, *http.Request, platform.User, RequestType) (ResponseType, *nhttp.ErrorWithStatus)
type handlerFunctionPostFormMultipart[RequestType any, ResponseType any] func(context.Context, *http.Request, RequestType) (*ResponseType, *nhttp.ErrorWithStatus)
type handlerFunctionPutAuthenticated[RequestType any] func(context.Context, *http.Request, platform.User, RequestType) (string, *nhttp.ErrorWithStatus)

func authenticatedHandlerBasic(f handlerBaseAuthenticated) http.Handler {
	return auth.NewEnsureAuth(func(w http.ResponseWriter, r *http.Request, u platform.User) {
		ctx := r.Context()
		e := f(ctx, w, r, u)
		if e != nil {
			respondErrorStatus(w, e)
			return
		}
		return
	})
}
func authenticatedHandlerDelete(f handlerFunctionDelete) http.Handler {
	return auth.NewEnsureAuth(func(w http.ResponseWriter, r *http.Request, u platform.User) {
		ctx := r.Context()
		e := f(ctx, r, u)
		if e != nil {
			respondErrorStatus(w, e)
			return
		}
		http.Error(w, "", http.StatusNoContent)
		return
	})
}

func authenticatedHandlerGetImage(f handlerFunctionGetImage) http.Handler {
	return auth.NewEnsureAuth(func(w http.ResponseWriter, r *http.Request, u platform.User) {
		ctx := r.Context()
		collection, uid, e := f(ctx, r, u)
		if e != nil {
			respondErrorStatus(w, e)
			return
		}
		file.ImageFileToWriter(collection, uid, w)
	})
}

func authenticatedHandlerJSON[T any](f handlerFunctionGetAuthenticated[T]) http.Handler {
	return auth.NewEnsureAuth(func(w http.ResponseWriter, r *http.Request, u platform.User) {
		ctx := r.Context()
		var body []byte
		var params resource.QueryParams
		err := decoder.Decode(&params, r.URL.Query())
		if err != nil {
			respondErrorStatus(w, nhttp.NewBadRequest("failed to decode query: %w", err))
			return
		}
		resp, e := f(ctx, r, u, params)
		w.Header().Set("Content-Type", "application/json")
		//log.Info().Str("template", template).Err(e).Msg("handler done")
		if e != nil {
			respondErrorStatus(w, e)
			return
		}
		body, err = json.Marshal(resp)
		if err != nil {
			respondErrorStatus(w, nhttp.NewError("failed to marshal json: %w", err))
			return
		}
		w.Write(body)
	})
}

func authenticatedHandlerJSONSlice[T any](f handlerFunctionGetSliceAuthenticated[T]) http.Handler {
	return auth.NewEnsureAuth(func(w http.ResponseWriter, r *http.Request, u platform.User) {
		ctx := r.Context()
		var body []byte
		var params resource.QueryParams
		err := decoder.Decode(&params, r.URL.Query())
		if err != nil {
			respondErrorStatus(w, nhttp.NewBadRequest("failed to decode query: %w", err))
			return
		}
		resp, e := f(ctx, r, u, params)
		w.Header().Set("Content-Type", "application/json")
		//log.Info().Str("template", template).Err(e).Msg("handler done")
		if e != nil {
			respondErrorStatus(w, e)
			return
		}
		if resp == nil {
			body, err = json.Marshal([]struct{}{})
		} else {
			body, err = json.Marshal(resp)
		}
		if err != nil {
			respondErrorStatus(w, nhttp.NewError("failed to marshal json: %w", err))
			return
		}
		w.Write(body)
	})
}
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
			respondErrorStatus(w, nhttp.NewError("failed to marshal json: %w", err))
			return
		}
		w.Write(body)
	})
}

func authenticatedHandlerJSONPut[RequestType any](f handlerFunctionPutAuthenticated[RequestType]) http.Handler {
	return auth.NewEnsureAuth(func(w http.ResponseWriter, r *http.Request, u platform.User) {
		w.Header().Set("Content-Type", "application/json")
		req, e := parseRequest[RequestType](r)
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
func handlerBasic(f handlerBase) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		e := f(ctx, w, r)
		if e != nil {
			respondErrorStatus(w, e)
			return
		}
	}
}
func handlerJSON[T any](f handlerFunctionGet[T]) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		var body []byte
		var params resource.QueryParams
		err := decoder.Decode(&params, r.URL.Query())
		if err != nil {
			respondErrorStatus(w, nhttp.NewBadRequest("failed to decode query: %w", err))
			return
		}
		resp, e := f(ctx, r, params)
		w.Header().Set("Content-Type", "application/json")
		//log.Info().Str("template", template).Err(e).Msg("handler done")
		if e != nil {
			respondErrorStatus(w, e)
			return
		}
		body, err = json.Marshal(resp)
		if err != nil {
			respondErrorStatus(w, nhttp.NewError("failed to marshal json: %w", err))
			return
		}
		w.Write(body)
	}
}
func handlerJSONSlice[T any](f handlerFunctionGetSlice[T]) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		var body []byte
		var params resource.QueryParams
		err := decoder.Decode(&params, r.URL.Query())
		if err != nil {
			respondErrorStatus(w, nhttp.NewBadRequest("failed to decode query: %w", err))
			return
		}
		resp, e := f(ctx, r, params)
		w.Header().Set("Content-Type", "application/json")
		//log.Info().Str("template", template).Err(e).Msg("handler done")
		if e != nil {
			respondErrorStatus(w, e)
			return
		}
		body, err = json.Marshal(resp)
		if err != nil {
			respondErrorStatus(w, nhttp.NewError("failed to marshal json: %w", err))
			return
		}
		w.Write(body)
	}
}

func handlerJSONPost[RequestType any, ResponseType any](f handlerFunctionPost[RequestType, ResponseType]) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		req, e := parseRequest[RequestType](r)
		if e != nil {
			serializeError(w, e)
			return
		}
		ctx := r.Context()
		resp, e := f(ctx, r, *req)
		if e != nil {
			serializeError(w, e)
			return
		}
		body, err := json.Marshal(resp)
		if err != nil {
			respondErrorStatus(w, nhttp.NewError("failed to marshal json: %w", err))
			return
		}
		w.Write(body)
	}
}

func handlerJSONPut[RequestType any, ResponseType any](f handlerFunctionPost[RequestType, ResponseType]) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		req, e := parseRequest[RequestType](r)
		if e != nil {
			serializeError(w, e)
			return
		}
		ctx := r.Context()
		resp, e := f(ctx, r, *req)
		if e != nil {
			serializeError(w, e)
			return
		}
		body, err := json.Marshal(resp)
		if err != nil {
			respondErrorStatus(w, nhttp.NewError("failed to marshal json: %w", err))
			return
		}
		w.Write(body)
	}
}
func handlerFormPost[RequestType any, ResponseType any](f handlerFunctionPostFormMultipart[RequestType, ResponseType]) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		err := r.ParseMultipartForm(32 << 12) // 128 MB buffer
		if err != nil {
			respondErrorStatus(w, nhttp.NewBadRequest("bad form: %w", err))
			return
		}
		var req RequestType
		err = decoder.Decode(&req, r.PostForm)
		if err != nil {
			respondErrorStatus(w, nhttp.NewBadRequest("decode form: %w", err))
			return
		}
		ctx := r.Context()
		resp, e := f(ctx, r, req)
		if e != nil {
			serializeError(w, e)
			return
		}
		body, err := json.Marshal(resp)
		if err != nil {
			respondErrorStatus(w, nhttp.NewError("failed to marshal json: %w", err))
			return
		}
		w.Write(body)
	}
}
func parseRequest[RequestType any](r *http.Request) (*RequestType, *nhttp.ErrorWithStatus) {
	var err error
	var req RequestType
	content_type := r.Header.Get("Content-Type")
	if content_type == "application/json" {
		body, e := io.ReadAll(r.Body)
		if e != nil {
			return nil, nhttp.NewError("Failed to read body: %w", err)
		}
		err = json.Unmarshal(body, &req)
	} else if content_type == "application/x-www-form-urlencoded" {
		e := r.ParseForm()
		if err != nil {
			return nil, nhttp.NewBadRequest("parsing form: %w", e)
		}
		err = decoder.Decode(&req, r.PostForm)
	} else {
		return nil, nhttp.NewBadRequest("unrecognized content type '%s'", content_type)
	}
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
func respondErrorStatus(w http.ResponseWriter, e *nhttp.ErrorWithStatus) {
	log.Warn().Int("status", e.Status).Err(e).Str("user message", e.Message).Msg("Responding with an error from api")
	body, err := json.Marshal(ErrorAPI{Message: e.Error()})
	if err != nil {
		log.Error().Err(err).Msg("failed to marshal error")
		http.Error(w, "{\"message\": \"boom. I can't even tell you what went wrong\"}", http.StatusInternalServerError)
		return
	}
	http.Error(w, string(body), e.Status)
}
