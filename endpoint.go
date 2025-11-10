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

func getArcgisOauthBegin(w http.ResponseWriter, r *http.Request) {
	expiration := 60
	authURL := buildArcGISAuthURL(ClientID, redirectURL(), expiration)
	http.Redirect(w, r, authURL, http.StatusFound)
}

func getArcgisOauthCallback(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")
	slog.Info("Handling oauth callback", slog.String("code", code))
	if code == "" {
		respondError(w, "Access code is empty", nil, http.StatusBadRequest)
		return
	}
	user, err := getAuthenticatedUser(r)
	if err != nil {
		respondError(w, "You're not currently authenticated, which really shouldn't happen.", err, http.StatusUnauthorized)
		return
	}
	err = handleOauthAccessCode(r.Context(), user, code)
	if err != nil {
		respondError(w, "Failed to handle access code", err, http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, BaseURL+"/", http.StatusFound)
}
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

func getReportConfirmation(w http.ResponseWriter, r *http.Request) {
	code := chi.URLParam(r, "code")
	err := htmlReportConfirmation(w, code)
	if err != nil {
		respondError(w, "Failed to generate report page", err, http.StatusInternalServerError)
	}
}

func getReportContribute(w http.ResponseWriter, r *http.Request) {
	code := chi.URLParam(r, "code")
	err := htmlReportContribute(w, code)
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

func getReportEvidence(w http.ResponseWriter, r *http.Request) {
	code := chi.URLParam(r, "code")
	err := htmlReportEvidence(w, code)
	if err != nil {
		respondError(w, "Failed to generate report page", err, http.StatusInternalServerError)
	}
}

func getReportSchedule(w http.ResponseWriter, r *http.Request) {
	code := chi.URLParam(r, "code")
	err := htmlReportSchedule(w, code)
	if err != nil {
		respondError(w, "Failed to generate report page", err, http.StatusInternalServerError)
	}
}

func getReportUpdate(w http.ResponseWriter, r *http.Request) {
	code := chi.URLParam(r, "code")
	err := htmlReportUpdate(w, code)
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
		has, err := hasFieldseekerConnection(r.Context(), user)
		if err != nil {
			respondError(w, "Failed to check for ArcGIS connection", err, http.StatusInternalServerError)
			return
		}
		if has {
			htmlDashboard(r.Context(), w, user)
			return
		} else {
			err = htmlOauthPrompt(w, user)
		}
	}
	if err != nil {
		respondError(w, "Failed to render root", err, http.StatusInternalServerError)
	}
}

func getServiceRequest(w http.ResponseWriter, r *http.Request) {
	err := htmlServiceRequest(w)
	if err != nil {
		respondError(w, "Failed to generate service request page", err, http.StatusInternalServerError)
	}
}

func getServiceRequestDetail(w http.ResponseWriter, r *http.Request) {
	code := chi.URLParam(r, "code")
	err := htmlServiceRequestDetail(w, code)
	if err != nil {
		respondError(w, "Failed to generate service request page", err, http.StatusInternalServerError)
	}
}

func getServiceRequestLocation(w http.ResponseWriter, r *http.Request) {
	err := htmlServiceRequestLocation(w)
	if err != nil {
		respondError(w, "Failed to generate service request location page", err, http.StatusInternalServerError)
	}
}

func getServiceRequestMosquito(w http.ResponseWriter, r *http.Request) {
	err := htmlServiceRequestMosquito(w)
	if err != nil {
		respondError(w, "Failed to generate service request mosquito page", err, http.StatusInternalServerError)
	}
}

func getServiceRequestPool(w http.ResponseWriter, r *http.Request) {
	err := htmlServiceRequestPool(w)
	if err != nil {
		respondError(w, "Failed to generate service request pool page", err, http.StatusInternalServerError)
	}
}

func getServiceRequestQuick(w http.ResponseWriter, r *http.Request) {
	err := htmlServiceRequestQuick(w)
	if err != nil {
		respondError(w, "Failed to generate service request quick page", err, http.StatusInternalServerError)
	}
}

func getServiceRequestQuickConfirmation(w http.ResponseWriter, r *http.Request) {
	err := htmlServiceRequestQuickConfirmation(w)
	if err != nil {
		respondError(w, "Failed to generate service request quick confirmation page", err, http.StatusInternalServerError)
	}
}

func getServiceRequestUpdates(w http.ResponseWriter, r *http.Request) {
	err := htmlServiceRequestUpdates(w)
	if err != nil {
		respondError(w, "Failed to generate service request updates page", err, http.StatusInternalServerError)
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
