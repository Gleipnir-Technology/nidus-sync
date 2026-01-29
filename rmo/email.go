package rmo

import (
	"fmt"
	"net/http"

	//"github.com/Gleipnir-Technology/nidus-sync/comms/email"
	"github.com/go-chi/chi/v5"
)

func getEmailByCode(w http.ResponseWriter, r *http.Request) {
	code := chi.URLParam(r, "code")
	if code == "" {
		http.Error(w, "You must specify a code", http.StatusBadRequest)
		return
	}
	fmt.Fprintf(w, "Pretend email contet for %s", code)
}
