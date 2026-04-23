package resource

import (
	"context"
	"net/http"
	"strconv"

	"github.com/Gleipnir-Technology/nidus-sync/config"
	nhttp "github.com/Gleipnir-Technology/nidus-sync/http"
	"github.com/gorilla/mux"
	"github.com/skip2/go-qrcode"
)

type qrcodeR struct {
	router *router
}

func QRCode(r *router) *qrcodeR {
	return &qrcodeR{
		router: r,
	}
}

func (res *qrcodeR) Mailer(ctx context.Context, w http.ResponseWriter, r *http.Request) *nhttp.ErrorWithStatus {
	vars := mux.Vars(r)
	code := vars["code"]
	if code == "" {
		return nhttp.NewBadRequest("There should always be a id")
	}
	content := config.MakeURLReport("/mailer/%s", code)
	return writeQRCode(w, r, content)
}
func (res *qrcodeR) Marketing(w http.ResponseWriter, r *http.Request) *nhttp.ErrorWithStatus {
	content := "https://nidus.cloud"
	return writeQRCode(w, r, content)
}

func (res *qrcodeR) Report(w http.ResponseWriter, r *http.Request) *nhttp.ErrorWithStatus {
	vars := mux.Vars(r)
	code := vars["code"]
	if code == "" {
		return nhttp.NewBadRequest("There should always be a code")
	}
	content := config.MakeURLNidus("/report/%s", code)
	return writeQRCode(w, r, content)
}
func writeQRCode(w http.ResponseWriter, r *http.Request, content string) *nhttp.ErrorWithStatus {
	// Get optional size parameter (default to 256)
	size := 256
	if sizeStr := r.URL.Query().Get("size"); sizeStr != "" {
		var err error
		size, err = strconv.Atoi(sizeStr)
		if err != nil {
			return nhttp.NewBadRequest("Invalid 'size' parameter, must be an integer")
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
			return nhttp.NewBadRequest("Invalid 'level' parameter, must be L, M, Q, or H")
		}
	}

	// Generate the QR code
	var qr *qrcode.QRCode
	var err error
	qr, err = qrcode.New(content, level)
	if err != nil {
		return nhttp.NewError("Error generating QR code: %w", err)
	}

	// Set the appropriate content type
	w.Header().Set("Content-Type", "image/png")

	// Generate PNG and write directly to the response writer
	png, err := qr.PNG(size)
	if err != nil {
		return nhttp.NewError("Error encoding QR code to PNG: %w", err)
	}

	_, err = w.Write(png)
	if err != nil {
		return nhttp.NewError("Error writing response: %w", err)
	}
	return nil
}
