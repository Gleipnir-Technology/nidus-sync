package resource

import (
	"fmt"
	"net/http"
	"strconv"

	nhttp "github.com/Gleipnir-Technology/nidus-sync/http"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/rs/zerolog/log"
)

type router struct {
	router *mux.Router
}

func NewRouter(r *mux.Router) *router {
	return &router{
		router: r,
	}
}
func (r *router) IDFromMux(req *http.Request) (int, *nhttp.ErrorWithStatus) {
	vars := mux.Vars(req)
	comm_id_str := vars["id"]
	if comm_id_str == "" {
		return 0, nhttp.NewBadRequest("no id provided")
	}
	comm_id, err := strconv.Atoi(comm_id_str)
	if err != nil {
		return 0, nhttp.NewBadRequest("can't turn report ID into an int: %w", err)
	}
	return comm_id, nil
}
func (r *router) IDFromURI(route string, uri string) (*int, error) {
	var match mux.RouteMatch
	req, _ := http.NewRequest("GET", uri, nil)
	if !r.router.Match(req, &match) {
		return nil, fmt.Errorf("URI does not match any known route: %s", uri)
	}

	route_name := match.Route.GetName()
	if route_name != route {
		return nil, fmt.Errorf("URI is not for the correct resource '%s', but for '%s'", route, route_name)
	}
	vars := match.Vars
	id_str, ok := vars["id"]
	if !ok {
		entry := log.Debug()
		for k, v := range vars {
			entry = entry.Str(k, v)
		}
		entry.Msg("current URI values")
		return nil, fmt.Errorf("No id found in URI %s", uri)
	}
	id, err := strconv.Atoi(id_str)
	if err != nil {
		return nil, fmt.Errorf("parse id: %w", err)
	}
	return &id, nil

}
func (r *router) UUIDFromURI(route string, uri string) (*uuid.UUID, error) {
	var match mux.RouteMatch
	req, _ := http.NewRequest("GET", uri, nil)
	if !r.router.Match(req, &match) {
		return nil, fmt.Errorf("URI does not match any known route: %s", uri)
	}

	route_name := match.Route.GetName()
	if route_name != route {
		return nil, fmt.Errorf("URI is not for the correct resource '%s', but for '%s'", route, route_name)
	}
	vars := match.Vars
	uuid_str, ok := vars["uuid"]
	if !ok {
		entry := log.Debug()
		for k, v := range vars {
			entry = entry.Str(k, v)
		}
		entry.Msg("current URI values")
		return nil, fmt.Errorf("No uuid found in URI %s", uri)
	}
	uid, err := uuid.Parse(uuid_str)
	if err != nil {
		return nil, fmt.Errorf("parse uuid: %w", err)
	}
	return &uid, nil
}
func (r *router) IDToURI(route string, id int) (string, error) {
	i := strconv.FormatInt(int64(id), 10)
	return r.IDStrToURI(route, i)
}
func (r *router) IDStrToURI(route string, id string) (string, error) {
	handler := r.router.Get(route)
	if handler == nil {
		return "", fmt.Errorf("nil handler '%s'", route)
	}
	uri, err := handler.URL("id", id)
	if err != nil {
		return "", fmt.Errorf("build uri: %w", err)
	}
	uri.Scheme = "https"
	return uri.String(), nil
}
func (r *router) SlugToURI(route string, slug string) (string, error) {
	handler := r.router.Get(route)
	if handler == nil {
		return "", fmt.Errorf("nil handler '%s'", route)
	}
	uri, err := handler.URL("slug", slug)
	if err != nil {
		return "", fmt.Errorf("build uri: %w", err)
	}
	uri.Scheme = "https"
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
	uri.Scheme = "https"
	result := uri.String()
	return &result, nil
}
