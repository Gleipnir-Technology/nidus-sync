package sync

import (
	"context"
	"net/http"

	"github.com/Gleipnir-Technology/nidus-sync/auth"
	"github.com/Gleipnir-Technology/nidus-sync/db"
	"github.com/Gleipnir-Technology/nidus-sync/db/models"
	"github.com/Gleipnir-Technology/nidus-sync/html"
	nhttp "github.com/Gleipnir-Technology/nidus-sync/http"
	"github.com/Gleipnir-Technology/nidus-sync/platform"
	"github.com/rs/zerolog/log"
)

type handlerFunctionGet[T any] func(context.Context, *http.Request, *models.Organization, *models.User) (*html.Response[T], *nhttp.ErrorWithStatus)
type wrappedHandler func(http.ResponseWriter, *http.Request)
type contentAuthenticated[T any] struct {
	C            T
	Config       html.ContentConfig
	Organization *models.Organization
	URL          html.ContentURL
	User         platform.User
}

// w http.ResponseWriter, r *http.Request, u *models.User) {
func authenticatedHandler[T any](f handlerFunctionGet[T]) http.Handler {
	return auth.NewEnsureAuth(func(w http.ResponseWriter, r *http.Request, u *models.User) {
		ctx := r.Context()
		userContent, err := auth.ContentForUser(ctx, u)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		org, err := u.Organization().One(ctx, db.PGInstance.BobDB)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if org == nil {
			http.Error(w, "nil org", http.StatusInternalServerError)
			return
		}
		resp, e := f(ctx, r, org, u)
		//log.Info().Str("template", template).Err(e).Msg("handler done")
		if e != nil {
			log.Warn().Int("status", e.Status).Err(e).Str("user message", e.Message).Msg("Responding with an error from sync pages")
			http.Error(w, e.Error(), e.Status)
			return
		}
		html.RenderOrError(w, resp.Template, contentAuthenticated[T]{
			C:            resp.Content,
			Config:       html.NewContentConfig(),
			Organization: org,
			URL:          html.NewContentURL(),
			User:         userContent,
		})
	})
}

type handlerFunctionPost[T any] func(context.Context, *http.Request, *models.Organization, *models.User, T) (string, *nhttp.ErrorWithStatus)

func authenticatedHandlerPost[T any](f handlerFunctionPost[T]) http.Handler {
	return auth.NewEnsureAuth(func(w http.ResponseWriter, r *http.Request, u *models.User) {
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
		org, err := u.Organization().One(ctx, db.PGInstance.BobDB)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		path, e := f(ctx, r, org, u, content)
		if e != nil {
			http.Error(w, e.Error(), e.Status)
			return
		}
		http.Redirect(w, r, path, http.StatusFound)
	})
}
func authenticatedHandlerPostMultipart[T any](f handlerFunctionPost[T]) http.Handler {
	return auth.NewEnsureAuth(func(w http.ResponseWriter, r *http.Request, u *models.User) {
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
		org, err := u.Organization().One(ctx, db.PGInstance.BobDB)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		path, e := f(ctx, r, org, u, content)
		if e != nil {
			http.Error(w, e.Error(), e.Status)
			return
		}
		http.Redirect(w, r, path, http.StatusFound)
	})
}
