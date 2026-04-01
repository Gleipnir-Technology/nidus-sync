package resource

import (
	"fmt"
	"strconv"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

type router struct {
	router *mux.Router
}

func NewRouter(r *mux.Router) *router {
	return &router{
		router: r,
	}
}
func (r *router) IDToURI(route string, id int) (string, error) {
	i := strconv.FormatInt(int64(id), 10)
	handler := r.router.Get(route)
	if handler == nil {
		return "", fmt.Errorf("nil handler '%s'", route)
	}
	uri, err := handler.URL("id", i)
	if err != nil {
		return "", fmt.Errorf("build uri: %w", err)
	}
	return uri.String(), nil
}

func (r *router) UUIDToURI(route string, u *uuid.UUID) (*string, error) {
	if u == nil {
		return nil, nil
	}
	handler := r.router.Get(route)
	if handler == nil {
		return nil, fmt.Errorf("nil handler '%s'", route)
	}
	uri, err := handler.URL("uuid", u.String())
	if err != nil {
		return nil, fmt.Errorf("build uri: %w", err)
	}
	result := uri.String()
	return &result, nil
}
