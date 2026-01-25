package text

import (
	"context"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/Gleipnir-Technology/nidus-sync/db"
	"github.com/Gleipnir-Technology/nidus-sync/db/enums"
	"github.com/Gleipnir-Technology/nidus-sync/db/models"
	"github.com/aarondl/opt/omit"
	"github.com/rs/zerolog/log"
	"github.com/stephenafamo/bob/types/pgtypes"
)

func convertToPGData(data map[string]string) pgtypes.HStore {
	result := pgtypes.HStore{}
	for k, v := range data {
		result[k] = sql.Null[string]{V: v, Valid: true}
	}
	return result
}

func ensureInDB(ctx context.Context, destination string) (err error) {
	_, err = models.FindCommsPhone(ctx, db.PGInstance.BobDB, destination)
	if err != nil {
		// doesn't exist
		if err.Error() == "sql: no rows in result set" {
			_, err = models.CommsPhones.Insert(&models.CommsPhoneSetter{
				E164:         omit.From(destination),
				IsSubscribed: omit.From(false),
			}).One(ctx, db.PGInstance.BobDB)
			if err != nil {
				return fmt.Errorf("Failed to insert new phone contact: %w", err)
			}
			log.Info().Str("phone", destination).Msg("Added text to the comms database")
			return nil
		}
		return fmt.Errorf("Unexpected error searching for phone contact: %w", err)
	}
	return nil
}

func insertTextLog(ctx context.Context, content string, destination string, source string, origin enums.CommsTextorigin) (err error) {
	_, err = models.CommsTextLogs.Insert(&models.CommsTextLogSetter{
		//ID:
		Content:     omit.From(content),
		Created:     omit.From(time.Now()),
		Destination: omit.From(destination),
		Origin:      omit.From(origin),
		Source:      omit.From(source),
	}).One(ctx, db.PGInstance.BobDB)

	return err
}
func generatePublicId(t enums.CommsMessagetypeemail, m map[string]string) string {
	if m == nil || len(m) == 0 {
		// Return hash of empty string for empty maps
		emptyHash := sha256.Sum256([]byte(""))
		return hex.EncodeToString(emptyHash[:])
	}

	// Get and sort keys for deterministic ordering
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	// Build a string with all key-value pairs
	var sb strings.Builder
	// Add type first
	sb.WriteString(fmt.Sprintf("type:%s,", t))
	for _, k := range keys {
		sb.WriteString(k)
		sb.WriteString(":") // Separator between key and value
		sb.WriteString(m[k])
		sb.WriteString(",") // Separator between pairs
	}

	// Compute SHA-256 hash
	hasher := sha256.New()
	hasher.Write([]byte(sb.String()))
	hashBytes := hasher.Sum(nil)

	// Convert to hex string and return
	return hex.EncodeToString(hashBytes)
}
