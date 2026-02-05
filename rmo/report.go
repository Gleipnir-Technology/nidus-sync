package rmo

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	//"github.com/Gleipnir-Technology/nidus-sync/config"
	"github.com/Gleipnir-Technology/nidus-sync/db"
	"github.com/Gleipnir-Technology/nidus-sync/db/enums"
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

type LatLngForm struct {
	Latitude      *float64
	Longitude     *float64
	MapZoom       float32
	AccuracyValue float32
	AccuracyType  enums.PublicreportAccuracytype
}

func parseLatLng(r *http.Request) (LatLngForm, error) {
	result := LatLngForm{
		AccuracyType:  enums.PublicreportAccuracytypeNone,
		AccuracyValue: 0.0,
		Latitude:      nil,
		Longitude:     nil,
		MapZoom:       0.0,
	}
	latitude_str := r.FormValue("latitude")
	longitude_str := r.FormValue("longitude")
	latlng_accuracy_type_str := r.PostFormValue("latlng-accuracy-type")
	latlng_accuracy_value_str := r.PostFormValue("latlng-accuracy-value")
	map_zoom_str := r.PostFormValue("map-zoom")

	var err error
	if latlng_accuracy_type_str != "" {
		err := result.AccuracyType.Scan(latlng_accuracy_type_str)
		if err != nil {
			return result, fmt.Errorf("Failed to parse accuracy type '%s': %w", latlng_accuracy_type_str, err)
		}
	}
	if latlng_accuracy_value_str != "" {
		var t float64
		t, err = strconv.ParseFloat(latlng_accuracy_value_str, 32)
		if err != nil {
			return result, fmt.Errorf("Failed to parse latlng_accuracy_value '%s': %w", latlng_accuracy_value_str, err)
		}
		result.AccuracyValue = float32(t)
	}

	if latitude_str != "" {
		var t float64
		t, err = strconv.ParseFloat(latitude_str, 64)
		if err != nil {
			return result, fmt.Errorf("Failed to parse latitude '%s': %w", latitude_str, err)
		}
		result.Latitude = &t
	}
	if longitude_str != "" {
		var t float64
		t, err := strconv.ParseFloat(longitude_str, 64)
		if err != nil {
			return result, fmt.Errorf("Failed to parse longitude '%s': %w", longitude_str, err)
		}
		result.Longitude = &t
	}

	if map_zoom_str != "" {
		var t float64
		t, err = strconv.ParseFloat(map_zoom_str, 32)
		if err != nil {
			return result, fmt.Errorf("Failed to parse map_zoom_str '%s': %w", map_zoom_str, err)
		} else {
			result.MapZoom = float32(t)
		}
	}
	return result, nil
}
