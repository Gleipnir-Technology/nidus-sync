package main

import (
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/Gleipnir-Technology/nidus-sync/db/models"
	"github.com/go-chi/chi/v5"
	"github.com/rs/zerolog/log"
	"github.com/skip2/go-qrcode"
)

func getArcgisOauthBegin(w http.ResponseWriter, r *http.Request) {
	authURL := buildArcGISAuthURL(ClientID)
	http.Redirect(w, r, authURL, http.StatusFound)
}

func getArcgisOauthCallback(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")
	log.Info().Str("code", code).Msg("Handling oauth callback")
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

func getCellDetails(w http.ResponseWriter, r *http.Request, user *models.User) {
	cell_str := chi.URLParam(r, "cell")
	if cell_str == "" {
		respondError(w, "There should always be a cell", nil, http.StatusBadRequest)
		return
	}
	cell, err := HexToInt64(cell_str)
	if err != nil {
		respondError(w, "Cannot convert provided cell to uint64", err, http.StatusBadRequest)
		return
	}
	htmlCell(r.Context(), w, user, cell)
}

func getFavicon(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "image/x-icon")

	http.ServeFile(w, r, "static/favicon.ico")
}

func getOAuthRefresh(w http.ResponseWriter, r *http.Request) {
	user, err := getAuthenticatedUser(r)
	if err != nil {
		http.Redirect(w, r, "/?next=/oauth/refresh", http.StatusFound)
	}
	htmlOauthPrompt(w, user)
}

func getPhoneCall(w http.ResponseWriter, r *http.Request) {
	htmlPhoneCall(w)
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
	htmlReport(w)
}

func getReportConfirmation(w http.ResponseWriter, r *http.Request) {
	code := chi.URLParam(r, "code")
	htmlReportConfirmation(w, code)
}

func getReportContribute(w http.ResponseWriter, r *http.Request) {
	code := chi.URLParam(r, "code")
	htmlReportContribute(w, code)
}

func getReportDetail(w http.ResponseWriter, r *http.Request) {
	code := chi.URLParam(r, "code")
	htmlReportDetail(w, code)
}

func getReportEvidence(w http.ResponseWriter, r *http.Request) {
	code := chi.URLParam(r, "code")
	htmlReportEvidence(w, code)
}

func getReportSchedule(w http.ResponseWriter, r *http.Request) {
	code := chi.URLParam(r, "code")
	htmlReportSchedule(w, code)
}

func getReportUpdate(w http.ResponseWriter, r *http.Request) {
	code := chi.URLParam(r, "code")
	htmlReportUpdate(w, code)
}

func getRoot(w http.ResponseWriter, r *http.Request) {
	user, err := getAuthenticatedUser(r)
	if err != nil && !errors.Is(err, &NoCredentialsError{}) {
		respondError(w, "Failed to get root", err, http.StatusInternalServerError)
		return
	}
	if user == nil {
		errorCode := r.URL.Query().Get("error")
		htmlSignin(w, errorCode)
		return
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
			htmlOauthPrompt(w, user)
			return
		}
	}
	if err != nil {
		respondError(w, "Failed to render root", err, http.StatusInternalServerError)
	}
}

func getServiceRequest(w http.ResponseWriter, r *http.Request) {
	htmlServiceRequest(w)
}

func getServiceRequestDetail(w http.ResponseWriter, r *http.Request) {
	code := chi.URLParam(r, "code")
	htmlServiceRequestDetail(w, code)
}

func getServiceRequestLocation(w http.ResponseWriter, r *http.Request) {
	htmlServiceRequestLocation(w)
}

func getServiceRequestMosquito(w http.ResponseWriter, r *http.Request) {
	htmlServiceRequestMosquito(w)
}

func getServiceRequestPool(w http.ResponseWriter, r *http.Request) {
	htmlServiceRequestPool(w)
}

func getServiceRequestQuick(w http.ResponseWriter, r *http.Request) {
	htmlServiceRequestQuick(w)
}

func getServiceRequestQuickConfirmation(w http.ResponseWriter, r *http.Request) {
	htmlServiceRequestQuickConfirmation(w)
}

func getServiceRequestUpdates(w http.ResponseWriter, r *http.Request) {
	htmlServiceRequestUpdates(w)
}

func getSettings(w http.ResponseWriter, r *http.Request, u *models.User) {
	htmlSettings(w, r, u)
}

func getSignin(w http.ResponseWriter, r *http.Request) {
	errorCode := r.URL.Query().Get("error")
	htmlSignin(w, errorCode)
}

func getSignup(w http.ResponseWriter, r *http.Request) {
	htmlSignup(w, r.URL.Path)
}

func getSource(w http.ResponseWriter, r *http.Request, u *models.User) {
	globalid := chi.URLParam(r, "globalid")
	if globalid == "" {
		respondError(w, "No globalid provided", nil, http.StatusBadRequest)
		return
	}
	htmlSource(w, r, u, globalid)
}

func getVectorTiles(w http.ResponseWriter, r *http.Request, u *models.User) {
	org_id := chi.URLParam(r, "org_id")
	tileset_id := chi.URLParam(r, "tileset_id")
	zoom := chi.URLParam(r, "zoom")
	x := chi.URLParam(r, "x")
	y := chi.URLParam(r, "y")
	format := chi.URLParam(r, "format")

	log.Info().Str("org_id", org_id).Str("tileset_id", tileset_id).Str("zoom", zoom).Str("x", x).Str("y", y).Str("format", format).Msg("Get vector tiles")

}
func respondError(w http.ResponseWriter, m string, e error, s int) {
	log.Warn().Int("status", s).Err(e).Msg("Responding with an error")
	http.Error(w, m, http.StatusBadRequest)
}

func postSignin(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		respondError(w, "Could not parse form", err, http.StatusBadRequest)
		return
	}

	username := r.FormValue("username")
	password := r.FormValue("password")

	log.Info().Str("username", username).Msg("Signin")

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

	log.Info().Str("username", username).Str("name", name).Str("password", strings.Repeat("*", len(password))).Msg("Signup")

	if terms != "on" {
		log.Warn().Msg("Terms not agreed")
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
