package main

import (
	"log/slog"
	"net/http"
)
func getFavicon(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "image/x-icon")

	http.ServeFile(w, r, "static/favicon.ico")
}

func getRoot(w http.ResponseWriter, r *http.Request) {
	err := htmlRoot(w, r.URL.Path)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
func getSignup(w http.ResponseWriter, r *http.Request) {
	err := htmlSignup(w, r.URL.Path)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func respondError(w http.ResponseWriter, m string, e error, s int) {
	slog.Error(m, slog.Int("status", s), slog.String("err", e.Error()))
	http.Error(w, m, http.StatusBadRequest)
}

func postSignup(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		respondError(w, "Could not parse form", err, http.StatusBadRequest)
		return
	}
	
	username := r.FormValue("username")
	name := r.FormValue("name")
	password := r.FormValue("password")
	terms := r.FormValue("terms")
	
	slog.Info("Signup", 
		slog.String("username", username),
		slog.String("name", name))
	
	if terms != "on" {
		slog.Error("Terms not agreed", slog.String("terms", terms))
		http.Error(w, "You must agree to the terms to register", http.StatusBadRequest)
		return
	}

	if err := signupUser(username, name, password); err != nil {
		respondError(w, "Failed to signup user", err, http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/", http.StatusFound)
}
