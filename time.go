package main

import (
	"time"

	"github.com/rs/zerolog/log"
)

// calculateCadenceVariance takes a slice of time.Time instances and returns
// the average time span between consecutive instances and how much each
// instance deviates from its expected position in an evenly-spaced sequence
// timestamps should be in descending order (most recent time in index 0)
func calculateCadenceVariance(timestamps []time.Time) (averageSpan time.Duration, deviations []time.Duration) {
	if len(timestamps) <= 1 {
		return 0, []time.Duration{}
	}

	// Calculate total span between first and last timestamp
	totalSpan := timestamps[0].Sub(timestamps[len(timestamps)-1])

	// Calculate average span
	averageSpan = totalSpan / time.Duration(len(timestamps)-1)
	if averageSpan < 0 {
		log.Error().Int64("average", int64(averageSpan)).Msg("Negative average")
		return 0, []time.Duration{}
	}

	// Calculate deviations
	deviations = make([]time.Duration, len(timestamps))

	// The first timestamp is our reference point (deviation = 0)
	deviations[0] = 0

	// Calculate expected timestamps based on perfect intervals
	for i := 1; i < len(timestamps); i++ {
		expectedTime := timestamps[0].Add(-averageSpan * time.Duration(i))
		deviations[i] = timestamps[i].Sub(expectedTime)
	}

	log.Info().Int64("average", int64(averageSpan)).Int("number deviations", len(deviations)).Msg("Calculated cadence")
	return averageSpan, deviations
}
