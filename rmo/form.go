package rmo

import (
	"net/http"
)

func postFormBool(r *http.Request, k string) *bool {
	v := r.PostFormValue(k)
	if v == "" {
		return nil
	}
	result := false
	if v == "on" {
		result = true
		return &result
	}
	return &result
}

func postFormValueOrNone(r *http.Request, k string) string {
	v := r.PostFormValue(k)
	if v == "" {
		return "none"
	}
	return v
}
