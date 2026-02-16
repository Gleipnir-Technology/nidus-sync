package db

import (
	"database/sql"
	"github.com/Gleipnir-Technology/bob/types/pgtypes"
	"github.com/rs/zerolog/log"
)

func ConvertFromPGData(d pgtypes.HStore) map[string]string {
	result := make(map[string]string, 0)
	for k, v := range d {
		value, err := v.Value()
		if err != nil {
			log.Warn().Err(err).Str("key", k).Msg("Failed to convert from HSTORE")
			continue
		}
		value_str, ok := value.(string)
		if !ok {
			log.Warn().Msg("Failed to convert to string")
		}
		result[k] = value_str
	}
	return result
}

func ConvertToPGData(data map[string]string) pgtypes.HStore {
	result := pgtypes.HStore{}
	for k, v := range data {
		result[k] = sql.Null[string]{V: v, Valid: true}
	}
	return result
}
