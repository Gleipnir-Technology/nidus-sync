package html

import (
	"context"
	"net/http"

	nhttp "github.com/Gleipnir-Technology/nidus-sync/http"
	"github.com/rs/zerolog/log"
)

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
			C: resp.Content,
		})
	}
}
