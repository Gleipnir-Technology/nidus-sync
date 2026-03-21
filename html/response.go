package html

import (
	"net/http"

	"github.com/rs/zerolog/log"
)

// Respond with an error that is visible to the user
func RespondError(w http.ResponseWriter, m string, e error, s int) {
	log.Warn().Int("status", s).Err(e).Str("user message", m).Msg("Responding with an error")
	http.Error(w, m, s)
}

type Response[T any] struct {
	Content  T
	Template string
}

func NewResponse[T any](template string, content T) *Response[T] {
	return &Response[T]{
		Content:  content,
		Template: template,
	}
}
