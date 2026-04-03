package html

import (
	"net/http"
)

func BoolFromForm(r *http.Request, k string) bool {
	s := r.PostFormValue(k)
	if s == "on" {
		return true
	}
	return false
}
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

func PostFormValueOrNone(r *http.Request, k string) string {
	v := r.PostFormValue(k)
	if v == "" {
		return "none"
	}
	return v
}
