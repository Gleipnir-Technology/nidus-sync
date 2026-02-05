package rmo

import (
	"encoding/json"
	"net/http"
	"strings"

	//"github.com/Gleipnir-Technology/nidus-sync/config"
	"github.com/Gleipnir-Technology/nidus-sync/db"
	"github.com/Gleipnir-Technology/nidus-sync/db/sql"
	//"github.com/go-chi/chi/v5"
	//"github.com/rs/zerolog/log"
)

type ReportSuggestion struct {
	ID   string `json:"id"`
	Type string `json:"type"`
	//Location string
}
type ReportSuggestionResponse struct {
	Reports []ReportSuggestion `json:"reports"`
}

func getReportSuggestion(w http.ResponseWriter, r *http.Request) {
	partial_report_id := r.FormValue("r")
	if partial_report_id == "" {
		respondError(w, "You need at least a bit of an 'r'", nil, http.StatusBadRequest)
		return
	}
	p := strings.ToUpper(partial_report_id) + "%"
	ctx := r.Context()
	rows, err := sql.PublicreportPublicIDSuggestion(p).All(ctx, db.PGInstance.BobDB)
	if err != nil {
		respondError(w, "Failed to query DB: %w", err, http.StatusInternalServerError)
		return
	}
	var result ReportSuggestionResponse
	for _, row := range rows {
		/*
			value, err := row.Location.Value()
			if err != nil {
				log.Warn().Err(err).Msg("Failed to get value")
				continue
			}
			value_str, ok := value.(string)
			if !ok {
				log.Warn().Msg("Failed to get location as string")
				continue
			}
			log.Debug().Str("location", value_str).Msg("Looking at row")
		*/
		result.Reports = append(result.Reports, ReportSuggestion{
			Type: row.TableName,
			ID:   row.PublicID,
			//Location: "",
		})
	}
	jsonBody, err := json.Marshal(result)
	if err != nil {
		respondError(w, "Failed to marshal JSON: %w", err, http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonBody)
}
