package html

import (
	"context"
	"net/http"

	nhttp "github.com/Gleipnir-Technology/nidus-sync/http"
	"github.com/gorilla/schema"
	"github.com/rs/zerolog/log"
)

var decoder = schema.NewDecoder()

type handlerFunctionGet[T any] func(context.Context, *http.Request) (*Response[T], *nhttp.ErrorWithStatus)

func MakeGet[T any](f handlerFunctionGet[T]) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		resp, e := f(ctx, r)
		if e != nil {
			log.Warn().Int("status", e.Status)
			http.Error(w, e.Error(), e.Status)
			return
		}
		RenderOrError(w, resp.Template, Content[T]{
			C:   resp.Content,
			URL: NewContentURL(),
		})
	}
}

type handlerFunctionPost[T any] func(context.Context, *http.Request, T) (string, *nhttp.ErrorWithStatus)

func MakePost[T any](f handlerFunctionPost[T]) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := r.ParseForm()
		if err != nil {
			RespondError(w, "Failed to parse form", err, http.StatusBadRequest)
			return
		}

		var content T

		err = decoder.Decode(&content, r.PostForm)
		if err != nil {
			RespondError(w, "Failed to decode form", err, http.StatusBadRequest)
			return
		}
		ctx := r.Context()
		path, e := f(ctx, r, content)
		if e != nil {
			http.Error(w, e.Error(), e.Status)
			return
		}
		http.Redirect(w, r, path, http.StatusFound)
	}
}
