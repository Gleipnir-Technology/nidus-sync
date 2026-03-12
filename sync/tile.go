package sync

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/Gleipnir-Technology/nidus-sync/platform"
)

func getTileGPS(w http.ResponseWriter, r *http.Request, u platform.User) {
	ctx := r.Context()
	if err := r.ParseForm(); err != nil {
		respondError(w, "Could not parse form", err, http.StatusBadRequest)
		return
	}
	lat_s := r.FormValue("lat")
	lng_s := r.FormValue("lng")
	level_s := r.FormValue("level")
	if lat_s == "" || lng_s == "" || level_s == "" {
		respondError(w, "you must specify lat, lng, and level", nil, http.StatusBadRequest)
		return
	}

	level, err := strconv.Atoi(level_s)
	if err != nil {
		respondError(w, "couldn't parse level", err, http.StatusBadRequest)
		return
	}
	lat, err := strconv.ParseFloat(lat_s, 10)
	if err != nil {
		respondError(w, "couldn't parse lat", err, http.StatusBadRequest)
		return
	}
	lng, err := strconv.ParseFloat(lng_s, 10)
	if err != nil {
		respondError(w, "couldn't parse lng", err, http.StatusBadRequest)
		return
	}
	img, err := platform.ImageAtPoint(ctx, u.Organization, uint(level), lat, lng)
	if err != nil {
		respondError(w, "image at point", err, http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "image/png")
	w.Header().Set("Content-Length", fmt.Sprintf("%d", len(img.Content)))
	_, err = io.Copy(w, bytes.NewBuffer(img.Content))
	if err != nil {
		respondError(w, "copy bytes", err, http.StatusInternalServerError)
		return
	}
}
