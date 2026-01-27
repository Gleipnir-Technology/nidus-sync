package text

import (
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"fmt"
	"sort"
	"strings"

	"github.com/Gleipnir-Technology/bob/types/pgtypes"
	"github.com/Gleipnir-Technology/nidus-sync/db/enums"
)

func convertToPGData(data map[string]string) pgtypes.HStore {
	result := pgtypes.HStore{}
	for k, v := range data {
		result[k] = sql.Null[string]{V: v, Valid: true}
	}
	return result
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
