package api

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/Gleipnir-Technology/nidus-sync/db/models"
	"github.com/go-chi/chi/v5"
)

func getTile(w http.ResponseWriter, r *http.Request, org *models.Organization, user *models.User) {
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
	fmt.Fprintf(w, "%d, %d, %d", x, y, z)
}
