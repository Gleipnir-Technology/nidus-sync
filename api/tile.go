package api

import (
	"net/http"
	"strconv"

	"github.com/Gleipnir-Technology/nidus-sync/platform"
	"github.com/go-chi/chi/v5"
	"github.com/rs/zerolog/log"
)

func getTile(w http.ResponseWriter, r *http.Request, user platform.User) {
	x_str := chi.URLParam(r, "x")
	y_str := chi.URLParam(r, "y")
	z_str := chi.URLParam(r, "z")

	x, err := strconv.Atoi(x_str)
	if err != nil {
		http.Error(w, "can't parse x as an integer", http.StatusBadRequest)
		return
	}
	y, err := strconv.Atoi(y_str)
	if err != nil {
		http.Error(w, "can't parse x as an integer", http.StatusBadRequest)
		return
	}
	z, err := strconv.Atoi(z_str)
	if err != nil {
		http.Error(w, "can't parse x as an integer", http.StatusBadRequest)
		return
	}
	err = platform.GetTile(r.Context(), w, user.Organization, uint(z), uint(y), uint(x))
	if err != nil {
		log.Error().Err(err).Msg("failed to do tile")
		http.Error(w, "failed to do tile", http.StatusInternalServerError)
		return
	}
}
