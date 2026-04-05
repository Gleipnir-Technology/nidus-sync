package resource

import (
// "github.com/gorilla/schema"
)

type QueryParams struct {
	Limit          *int    `schema:"limit"`
	OrganizationID *int    `schema:"org"`
	Query          *string `schema:"query"`
	Sort           *string `schema:"sort"`
	Type           *string `schema:"type"`
}

func (qp QueryParams) SortOrDefault(default_name string, ascending bool) (string, bool) {
	if qp.Sort == nil {
		return default_name, ascending
	}
	s := *qp.Sort
	if s == "" {
		return default_name, ascending
	}
	a := true
	if s[0] == '-' {
		a = false
	}
	if s[0] == '+' || s[0] == '-' {
		s = s[1:]
	}
	return s, a
}
