package sync

import (
	"context"
	"net/http"

	"github.com/Gleipnir-Technology/nidus-sync/auth"
	"github.com/Gleipnir-Technology/nidus-sync/html"
	nhttp "github.com/Gleipnir-Technology/nidus-sync/http"
	"github.com/Gleipnir-Technology/nidus-sync/platform"
	"github.com/gorilla/schema"
	"github.com/rs/zerolog/log"
)

var decoder = schema.NewDecoder()

type handlerFunctionGet[T any] func(context.Context, *http.Request, platform.User) (*html.Response[T], *nhttp.ErrorWithStatus)
type wrappedHandler func(http.ResponseWriter, *http.Request)
type contentAuthenticated[T any] struct {
	C            T
	Config       html.ContentConfig
	Organization platform.Organization
	URL          html.ContentURL
	User         platform.User
}

// w http.ResponseWriter, r *http.Request, u platform.User) {
func authenticatedHandler[T any](f handlerFunctionGet[T]) http.Handler {
	return auth.NewEnsureAuth(func(w http.ResponseWriter, r *http.Request, u platform.User) {
		ctx := r.Context()
		resp, e := f(ctx, r, u)
		//log.Info().Str("template", template).Err(e).Msg("handler done")
		if e != nil {
			log.Warn().Int("status", e.Status).Err(e).Str("user message", e.Message).Msg("Responding with an error from sync pages")
			http.Error(w, e.Error(), e.Status)
			return
		}
		html.RenderOrError(w, resp.Template, contentAuthenticated[T]{
			C:            resp.Content,
			Config:       html.NewContentConfig(),
			Organization: u.Organization,
			URL:          html.NewContentURL(),
			User:         u,
		})
	})
}

type handlerFunctionPost[T any] func(context.Context, *http.Request, platform.User, T) (string, *nhttp.ErrorWithStatus)

func authenticatedHandlerPost[T any](f handlerFunctionPost[T]) http.Handler {
	return auth.NewEnsureAuth(func(w http.ResponseWriter, r *http.Request, u platform.User) {
		err := r.ParseForm()
		if err != nil {
			respondError(w, "Failed to parse form", err, http.StatusBadRequest)
			return
		}

		var content T

		err = decoder.Decode(&content, r.PostForm)
		if err != nil {
			respondError(w, "Failed to decode form", err, http.StatusBadRequest)
			return
		}
		ctx := r.Context()
		path, e := f(ctx, r, u, content)
		if e != nil {
			http.Error(w, e.Error(), e.Status)
			return
		}
		http.Redirect(w, r, path, http.StatusFound)
	})
}
func authenticatedHandlerPostMultipart[T any](f handlerFunctionPost[T]) http.Handler {
	return auth.NewEnsureAuth(func(w http.ResponseWriter, r *http.Request, u platform.User) {
		err := r.ParseMultipartForm(32 << 10) // 32 MB buffer
		if err != nil {
			respondError(w, "Failed to parse form", err, http.StatusBadRequest)
			return
		}

		var content T

		err = decoder.Decode(&content, r.PostForm)
		if err != nil {
			respondError(w, "Failed to decode form", err, http.StatusBadRequest)
			return
		}
		ctx := r.Context()
		path, e := f(ctx, r, u, content)
		if e != nil {
			http.Error(w, e.Error(), e.Status)
			return
		}
		http.Redirect(w, r, path, http.StatusFound)
	})
}
