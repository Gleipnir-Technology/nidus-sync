package api

import (
	"net/http"
	"strconv"

	"github.com/Gleipnir-Technology/nidus-sync/platform"
	"github.com/gorilla/mux"
	"github.com/rs/zerolog/log"
)

func getTile(w http.ResponseWriter, r *http.Request, user platform.User) {
	vars := mux.Vars(r)
	x_str := vars["x"]
	y_str := vars["y"]
	z_str := vars["z"]

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
