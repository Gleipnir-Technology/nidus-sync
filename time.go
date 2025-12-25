package main

import (
	"sort"
	"time"
	//"github.com/rs/zerolog/log"
)

// TreatmentModel represents the calculated model for a year's treatments
type TreatmentModel struct {
	Year           int
	SeasonStart    time.Time
	SeasonEnd      time.Time
	Interval       time.Duration
	ActualDates    []time.Time
	PredictedDates []time.Time
	Errors         []time.Duration
}

func modelTreatment(treatments []Treatment) []TreatmentModel {
	treatment_times := make([]time.Time, 0)
	for _, treatment := range treatments {
		if treatment.Date != nil {
			treatment_times = append(treatment_times, *treatment.Date)
		}
	}
	models := calculateTreatmentModels(treatment_times)
	/*models_by_year := make(map[int]TreatmentModel)
	for _, m := range models {
		models_by_year[m.Year] = m
	}*/
	offset := 0
	for _, model := range models {
		for _, e := range model.Errors {
			treatments[offset].CadenceDelta = e
			offset = offset + 1
		}
	}
	/*
		for i, treatment := range treatments {
			model
			treatment.CadenceDelta = deltas[i]
			treatments[i] = treatment
		}
	*/
	/*cadence, deltas := calculateCadenceVariance(treatment_times)
	for i, treatment := range treatments {
		if i >= len(deltas) {
			break
		}
		treatment.CadenceDelta = deltas[i]
		treatments[i] = treatment
	}*/
	return models
}

// calculateTreatmentModels segments treatments by year and calculates a model for each year
func calculateTreatmentModels(treatments []time.Time) []TreatmentModel {
	// Group treatments by year
	yearMap := make(map[int][]time.Time)
	for _, t := range treatments {
		year := t.Year()
		yearMap[year] = append(yearMap[year], t)
	}

	// Calculate a model for each year
	var models []TreatmentModel
	for year, dates := range yearMap {
		// Sort dates within the year
		sort.Slice(dates, func(i, j int) bool {
			return dates[i].Before(dates[j])
		})

		model := calculateYearModel(year, dates)
		models = append(models, model)
	}

	// Sort models by year
	sort.Slice(models, func(i, j int) bool {
		return models[i].Year < models[j].Year
	})

	return models
}

// calculateYearModel creates a model for a specific year using linear regression
func calculateYearModel(year int, dates []time.Time) TreatmentModel {
	n := len(dates)
	if n < 2 {
		// Not enough data for a model
		return TreatmentModel{
			Year:        year,
			ActualDates: dates,
		}
	}

	// Convert dates to numeric values (seconds since epoch)
	var x []float64
	var y []float64
	for i, date := range dates {
		x = append(x, float64(i))
		y = append(y, float64(date.Unix()))
	}

	// Calculate linear regression
	slope, intercept := linearRegression(x, y)

	// Convert back to time.Time and time.Duration
	startTime := time.Unix(int64(intercept), 0)
	intervalSeconds := int64(slope)
	interval := time.Duration(intervalSeconds) * time.Second

	// Calculate end of season
	endTime := time.Unix(int64(intercept+slope*float64(n-1)), 0)

	// Generate predicted dates and calculate errors
	var predictedDates []time.Time
	var errors []time.Duration

	for i := 0; i < n; i++ {
		predicted := time.Unix(int64(intercept+slope*float64(i)), 0)
		predictedDates = append(predictedDates, predicted)

		// Calculate error
		actualTime := dates[i]
		error := actualTime.Sub(predicted)
		errors = append(errors, error)
	}

	return TreatmentModel{
		Year:           year,
		SeasonStart:    startTime,
		SeasonEnd:      endTime,
		Interval:       interval,
		ActualDates:    dates,
		PredictedDates: predictedDates,
		Errors:         errors,
	}
}

// linearRegression calculates the slope and intercept of the best-fit line
func linearRegression(x, y []float64) (float64, float64) {
	n := float64(len(x))
	if n < 2 {
		return 0, 0
	}

	var sumX, sumY, sumXY, sumX2 float64
	for i := 0; i < len(x); i++ {
		sumX += x[i]
		sumY += y[i]
		sumXY += x[i] * y[i]
		sumX2 += x[i] * x[i]
	}

	// Calculate slope
	slope := (n*sumXY - sumX*sumY) / (n*sumX2 - sumX*sumX)

	// Calculate intercept
	intercept := (sumY - slope*sumX) / n

	return slope, intercept
}
