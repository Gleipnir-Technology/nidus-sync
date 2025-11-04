package main

import (
	"log"
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

func postSignup(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		log.Printf("Error parsing form: %v", err)
		http.Error(w, "Failed to process form", http.StatusBadRequest)
		return
	}
	
	email := r.FormValue("email")
	name := r.FormValue("name")
	terms := r.FormValue("terms")
	
	log.Printf("Signup - Email: %s, Name: %s, Terms: %s", email, name, terms)
	
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Form received"))
}
