package sync

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/Gleipnir-Technology/nidus-sync/config"
	"github.com/Gleipnir-Technology/nidus-sync/html"
	"github.com/go-chi/chi/v5"
	//"github.com/rs/zerolog/log"
	"github.com/skip2/go-qrcode"
)

// Unauthenticated pages
/*
	admin                           = buildTemplate("admin", "base")
	dataEntry                       = buildTemplate("data-entry", "base")
	dataEntryGood                   = buildTemplate("data-entry-good", "base")
	dataEntryBad                    = buildTemplate("data-entry-bad", "base")
	dispatch                        = buildTemplate("dispatch", "base")
	dispatchResults                 = buildTemplate("dispatch-results", "base")
	mockRoot                        = buildTemplate("mock-root", "base")
	reportPage                      = buildTemplate("report", "base")
	reportConfirmation              = buildTemplate("report-confirmation", "base")
	reportContribute                = buildTemplate("report-contribute", "base")
	reportDetail                    = buildTemplate("report-detail", "base")
	reportEvidence                  = buildTemplate("report-evidence", "base")
	reportSchedule                  = buildTemplate("report-schedule", "base")
	reportUpdate                    = buildTemplate("report-update", "base")
	serviceRequest                  = buildTemplate("service-request", "base")
	serviceRequestDetail            = buildTemplate("service-request-detail", "base")
	serviceRequestLocation          = buildTemplate("service-request-location", "base")
	serviceRequestMosquito          = buildTemplate("service-request-mosquito", "base")
	serviceRequestPool              = buildTemplate("service-request-pool", "base")
	serviceRequestQuick             = buildTemplate("service-request-quick", "base")
	serviceRequestQuickConfirmation = buildTemplate("service-request-quick-confirmation", "base")
	serviceRequestUpdates           = buildTemplate("service-request-updates", "base")
	settingRoot                     = buildTemplate("setting-mock", "base")
	settingIntegration              = buildTemplate("setting-integration", "base")
	settingPesticide                = buildTemplate("setting-pesticide", "base")
	settingPesticideAdd             = buildTemplate("setting-pesticide-add", "base")
	settingUsers                    = buildTemplate("setting-user", "base")
	settingUsersAdd                 = buildTemplate("setting-user-add", "base")
*/

func getQRCodeReport(w http.ResponseWriter, r *http.Request) {
	code := chi.URLParam(r, "code")
	if code == "" {
		respondError(w, "There should always be a code", nil, http.StatusBadRequest)
	}
	content := config.MakeURLNidus("/report/%s", code)
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

func mock(t string, w http.ResponseWriter, code string) {
	data := ContentMock{
		DistrictName: "Delta MVCD",
		URLs: ContentMockURLs{
			Dispatch:            "/mock/dispatch",
			DispatchResults:     "/mock/dispatch-results",
			ReportConfirmation:  fmt.Sprintf("/mock/report/%s/confirm", code),
			ReportDetail:        fmt.Sprintf("/mock/report/%s", code),
			ReportContribute:    fmt.Sprintf("/mock/report/%s/contribute", code),
			ReportEvidence:      fmt.Sprintf("/mock/report/%s/evidence", code),
			ReportSchedule:      fmt.Sprintf("/mock/report/%s/schedule", code),
			ReportUpdate:        fmt.Sprintf("/mock/report/%s/update", code),
			Root:                "/mock",
			Setting:             "/mock/setting",
			SettingIntegration:  "/mock/setting/integration",
			SettingPesticide:    "/mock/setting/pesticide",
			SettingPesticideAdd: "/mock/setting/pesticide/add",
			SettingUser:         "/mock/setting/user",
			SettingUserAdd:      "/mock/setting/user/add",
		},
	}
	html.RenderOrError(w, t, data)
}

func renderMock(templateName string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		code := chi.URLParam(r, "code")
		if code == "" {
			code = "abc-123"
		}
		mock(templateName, w, code)
	}
}
