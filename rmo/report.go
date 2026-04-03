package rmo

import (
	"encoding/json"
	"net/http"
	"strings"

	//"github.com/Gleipnir-Technology/nidus-sync/config"
	"github.com/Gleipnir-Technology/bob"
	"github.com/Gleipnir-Technology/bob/dialect/psql"
	"github.com/Gleipnir-Technology/bob/dialect/psql/sm"
	"github.com/Gleipnir-Technology/nidus-sync/db"
	//"github.com/gorilla/mux"
	"github.com/stephenafamo/scan"
	//"github.com/rs/zerolog/log"
)

type ReportSuggestion struct {
	ID string `json:"id"`
	//Type string `json:"type"`
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
	p := partialSearchParam(partial_report_id)
	ctx := r.Context()
	/*
		rows, err := sql.PublicreportPublicIDSuggestion(p).All(ctx, db.PGInstance.BobDB)
		if err != nil {
			respondError(w, "Failed to query DB: %w", err, http.StatusInternalServerError)
			return
		}
	*/
	type _Row struct {
		Location string `db:"location"`
		PublicID string `db:"public_id"`
	}
	rows, err := bob.All(ctx, db.PGInstance.BobDB, psql.Select(
		sm.Columns("public_id", "location"),
		sm.From("publicreport.report"),
		sm.Where(
			psql.Quote("public_id").Like(psql.Arg(p)),
		),
	), scan.StructMapper[_Row]())

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
			//Type: row.TableName,
			ID: row.PublicID,
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

func partialSearchParam(p string) string {
	result := strings.ReplaceAll(p, "-", "")
	result = strings.ToUpper(result)
	return result + "%"
}
