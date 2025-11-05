package main

import (
	"errors"
	"log/slog"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/skip2/go-qrcode"
)
func getFavicon(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "image/x-icon")

	http.ServeFile(w, r, "static/favicon.ico")
}

func getQRCodeReport(w http.ResponseWriter, r *http.Request) {
	code := chi.URLParam(r, "code")
	if code == "" {
		respondError(w, "There should always be a code", nil, http.StatusBadRequest)
	}
	content := BaseURL + "/report/" + code
	// Get optional size parameter (default to 256)
	size := 256
	if sizeStr := r.URL.Query().Get("size"); sizeStr != "" {
		var err error
		size, err = strconv.Atoi(sizeStr)
		if err != nil {
			http.Error(w, "Invalid 'size' parameter, must be an integer", http.StatusBadRequest)
			return
		}
	}

	// Get optional error correction level (default to Medium)
	level := qrcode.Medium
	if levelStr := r.URL.Query().Get("level"); levelStr != "" {
		switch levelStr {
		case "L", "l":
			level = qrcode.Low
		case "M", "m":
			level = qrcode.Medium
		case "Q", "q":
			level = qrcode.High
		case "H", "h":
			level = qrcode.Highest
		default:
			respondError(w, "Invalid 'level' parameter, must be L, M, Q, or H", nil, http.StatusBadRequest)
			return
		}
	}

	// Generate the QR code
	var qr *qrcode.QRCode
	var err error
	qr, err = qrcode.New(content, level)
	if err != nil {
		respondError(w, "Error generating QR code", err, http.StatusInternalServerError)
		return
	}

	// Set the appropriate content type
	w.Header().Set("Content-Type", "image/png")
	
	// Generate PNG and write directly to the response writer
	png, err := qr.PNG(size)
	if err != nil {
		respondError(w, "Error encoding QR code to PNG", err, http.StatusInternalServerError)
		return
	}

	_, err = w.Write(png)
	if err != nil {
		respondError(w, "Error writing response", err, http.StatusInternalServerError)
	}
}
func getReport(w http.ResponseWriter, r *http.Request) {
	//org := r.URL.Query().Get("org")
	err := htmlReport(w)
	if err != nil {
		respondError(w, "Failed to generate report page", err, http.StatusInternalServerError)
	}
}
func getReportDetail(w http.ResponseWriter, r *http.Request) {
	code := chi.URLParam(r, "code")
	err := htmlReportDetail(w, code)
	if err != nil {
		respondError(w, "Failed to generate report page", err, http.StatusInternalServerError)
	}
}

func getRoot(w http.ResponseWriter, r *http.Request) {
	user, err := getAuthenticatedUser(r)
	if err != nil && !errors.Is(err, &NoCredentialsError{}) {
		respondError(w, "Failed to get root", err, http.StatusInternalServerError)
		return
	}
	if user == nil {
		errorCode := r.URL.Query().Get("error")
		err = htmlSignin(w, errorCode)
	} else {
		err = htmlDashboard(w, user)
	}
	if err != nil {
		respondError(w, "Failed to render root", err, http.StatusInternalServerError)
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

func postSignin(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		respondError(w, "Could not parse form", err, http.StatusBadRequest)
		return
	}
	
	username := r.FormValue("username")
	password := r.FormValue("password")
	
	slog.Info("Signin", 
		slog.String("username", username),
		slog.String("password", strings.Repeat("*", len(password))))
	

	_, err := signinUser(r, username, password)
	if err != nil {
		if errors.Is(err, InvalidCredentials{}) {
			http.Redirect(w, r, "/?error=invalid-credentials", http.StatusFound)
			return
		}
		respondError(w, "Failed to signin user", err, http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/", http.StatusFound)
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
		slog.String("name", name),
		slog.String("password", strings.Repeat("*", len(password))))
	
	if terms != "on" {
		slog.Error("Terms not agreed", slog.String("terms", terms))
		http.Error(w, "You must agree to the terms to register", http.StatusBadRequest)
		return
	}

	user, err := signupUser(username, name, password)
	if err != nil {
		respondError(w, "Failed to signup user", err, http.StatusInternalServerError)
		return
	}

	addUserSession(r, user)

	http.Redirect(w, r, "/", http.StatusFound)
}
