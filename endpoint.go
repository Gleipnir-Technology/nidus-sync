package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/Gleipnir-Technology/nidus-sync/auth"
	"github.com/Gleipnir-Technology/nidus-sync/background"
	"github.com/Gleipnir-Technology/nidus-sync/config"
	"github.com/Gleipnir-Technology/nidus-sync/db/models"
	"github.com/Gleipnir-Technology/nidus-sync/htmlpage"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"github.com/skip2/go-qrcode"
)

func getArcgisOauthBegin(w http.ResponseWriter, r *http.Request) {
	authURL := config.BuildArcGISAuthURL(config.ClientID)
	http.Redirect(w, r, authURL, http.StatusFound)
}

func getArcgisOauthCallback(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")
	log.Info().Str("code", code).Msg("Handling oauth callback")
	if code == "" {
		respondError(w, "Access code is empty", nil, http.StatusBadRequest)
		return
	}
	user, err := auth.GetAuthenticatedUser(r)
	if err != nil {
		respondError(w, "You're not currently authenticated, which really shouldn't happen.", err, http.StatusUnauthorized)
		return
	}
	err = background.HandleOauthAccessCode(r.Context(), user, code)
	if err != nil {
		respondError(w, "Failed to handle access code", err, http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, config.MakeURLSync("/"), http.StatusFound)
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
	htmlpage.Cell(r.Context(), w, user, cell)
}

func getFavicon(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "image/x-icon")

	http.ServeFile(w, r, "static/favicon.ico")
}

func getOAuthRefresh(w http.ResponseWriter, r *http.Request) {
	user, err := auth.GetAuthenticatedUser(r)
	if err != nil {
		http.Redirect(w, r, "/?next=/oauth/refresh", http.StatusFound)
		return
	}
	htmlpage.OauthPrompt(w, user)
}

func getQRCodeReport(w http.ResponseWriter, r *http.Request) {
	code := chi.URLParam(r, "code")
	if code == "" {
		respondError(w, "There should always be a code", nil, http.StatusBadRequest)
	}
	content := config.MakeURLSync("/report/" + code)
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

func getRoot(w http.ResponseWriter, r *http.Request) {
	user, err := auth.GetAuthenticatedUser(r)
	if err != nil {
		// No credentials or user not found: go to login
		if errors.Is(err, &auth.NoCredentialsError{}) || errors.Is(err, &auth.NoUserError{}) {
			http.Redirect(w, r, "/signin", http.StatusFound)
			return
		} else {
			respondError(w, "Failed to get root", err, http.StatusInternalServerError)
			return
		}
	}
	if user == nil {
		errorCode := r.URL.Query().Get("error")
		htmlpage.Signin(w, errorCode)
		return
	} else {
		has, err := background.HasFieldseekerConnection(r.Context(), user)
		if err != nil {
			respondError(w, "Failed to check for ArcGIS connection", err, http.StatusInternalServerError)
			return
		}
		if has {
			htmlpage.Dashboard(r.Context(), w, user)
			return
		} else {
			htmlpage.OauthPrompt(w, user)
			return
		}
	}
	if err != nil {
		respondError(w, "Failed to render root", err, http.StatusInternalServerError)
	}
}

func getSettings(w http.ResponseWriter, r *http.Request, u *models.User) {
	htmlpage.Settings(w, r, u)
}

func getSignin(w http.ResponseWriter, r *http.Request) {
	errorCode := r.URL.Query().Get("error")
	htmlpage.Signin(w, errorCode)
}

func getSignup(w http.ResponseWriter, r *http.Request) {
	htmlpage.Signup(w, r.URL.Path)
}

func getSource(w http.ResponseWriter, r *http.Request, u *models.User) {
	globalid_s := chi.URLParam(r, "globalid")
	if globalid_s == "" {
		respondError(w, "No globalid provided", nil, http.StatusBadRequest)
		return
	}
	globalid, err := uuid.Parse(globalid_s)
	if err != nil {
		respondError(w, "globalid is not a UUID", nil, http.StatusBadRequest)
		return
	}
	htmlpage.Source(w, r, u, globalid)
}

func postSMS(w http.ResponseWriter, r *http.Request) {
	// Log all request headers
	for name, values := range r.Header {
		for _, value := range values {
			log.Info().Str("name", name).Str("value", value).Msg("header")
		}
	}

	// Read the request body
	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		//return nil, fmt.Errorf("failed to read request body: %w", err)
		respondError(w, "Failed to read request body", err, http.StatusInternalServerError)
		return
	}
	log.Info().Str("body", string(bodyBytes)).Msg("body")
	// Close the original body
	defer r.Body.Close()

	// Parse JSON into webhook struct
	var body SMSWebhookBody
	if err := json.Unmarshal(bodyBytes, &body); err != nil {
		respondError(w, "Failed to parse JSON", err, http.StatusBadRequest)
		return
	}

	if err := handleSMSMessage(&body.Data); err != nil {
		log.Error().Err(err).Msg("Failed to handle SMS Message")
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("ok"))
}
func getSMS(w http.ResponseWriter, r *http.Request) {
	org := chi.URLParam(r, "org")

	to := r.URL.Query().Get("error")
	from := r.URL.Query().Get("error")
	message := r.URL.Query().Get("error")
	files := r.URL.Query().Get("error")
	id := r.URL.Query().Get("error")
	date := r.URL.Query().Get("error")

	log.Info().Str("org", org).Str("to", to).Str("from", from).Str("message", message).Str("files", files).Str("id", id).Str("date", date).Msg("Got SMS Message")
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-type", "text/plain")
	// Signifies to Voip.ms that the callback worked.
	fmt.Fprintf(w, "ok")
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

// Respond with an error that is visible to the user
func respondError(w http.ResponseWriter, m string, e error, s int) {
	log.Warn().Int("status", s).Err(e).Str("user message", m).Msg("Responding with an error")
	http.Error(w, m, s)
}

func postSignin(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		respondError(w, "Could not parse form", err, http.StatusBadRequest)
		return
	}

	username := r.FormValue("username")
	password := r.FormValue("password")

	log.Info().Str("username", username).Msg("Signin")

	_, err := auth.SigninUser(r, username, password)
	if err != nil {
		if errors.Is(err, auth.InvalidCredentials{}) {
			http.Redirect(w, r, "/signin?error=invalid-credentials", http.StatusFound)
			return
		}
		if errors.Is(err, auth.InvalidUsername{}) {
			http.Redirect(w, r, "/signin?error=invalid-credentials", http.StatusFound)
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

	user, err := auth.SignupUser(r.Context(), username, name, password)
	if err != nil {
		respondError(w, "Failed to signup user", err, http.StatusInternalServerError)
		return
	}

	auth.AddUserSession(r, user)

	http.Redirect(w, r, "/", http.StatusFound)
}

func renderMock(templateName string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		code := chi.URLParam(r, "code")
		if code == "" {
			code = "abc-123"
		}
		htmlpage.Mock(templateName, w, code)
	}
}
