package email

import (
	"context"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/Gleipnir-Technology/bob/types/pgtypes"
	"github.com/Gleipnir-Technology/nidus-sync/db"
	"github.com/Gleipnir-Technology/nidus-sync/db/enums"
	"github.com/Gleipnir-Technology/nidus-sync/db/models"
	"github.com/aarondl/opt/omit"
	"github.com/aarondl/opt/omitnull"
	"github.com/rs/zerolog/log"
)

func convertToPGData(data map[string]string) pgtypes.HStore {
	result := pgtypes.HStore{}
	for k, v := range data {
		result[k] = sql.Null[string]{V: v, Valid: true}
	}
	return result
}

func convertFromPGData(d pgtypes.HStore) map[string]string {
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

func ensureInDB(ctx context.Context, destination string) (err error) {
	_, err = models.FindCommsEmailContact(ctx, db.PGInstance.BobDB, destination)
	if err != nil {
		// doesn't exist
		if err.Error() == "sql: no rows in result set" {
			public_id := fmt.Sprintf("%x", sha256.Sum256([]byte(destination)))
			_, err = models.CommsEmailContacts.Insert(&models.CommsEmailContactSetter{
				Address:      omit.From(destination),
				Confirmed:    omit.From(false),
				IsSubscribed: omit.From(false),
				PublicID:     omit.From(public_id),
			}).One(ctx, db.PGInstance.BobDB)
			if err != nil {
				return fmt.Errorf("Failed to insert new email: %w", err)
			}
			log.Info().Str("email", destination).Msg("Added email to the comms database")
			return nil
		}
		return fmt.Errorf("Unexpected error searching for contact: %w", err)
	}
	return nil
}

func insertEmailLog(ctx context.Context, data map[string]string, destination string, public_id string, source string, subject string, template_id int32) (err error) {
	data_for_insert := convertToPGData(data)
	var type_ enums.CommsMessagetypeemail
	switch template_id {
	case templateReportNotificationConfirmationID:
		type_ = enums.CommsMessagetypeemailReportNotificationConfirmation
	case templateInitialID:
		type_ = enums.CommsMessagetypeemailInitialContact
	default:
		return fmt.Errorf("Unrecognized template ID %d", template_id)
	}
	_, err = models.CommsEmailLogs.Insert(&models.CommsEmailLogSetter{
		//ID:
		Created:        omit.From(time.Now()),
		DeliveryStatus: omit.From("initial"),
		Destination:    omit.From(destination),
		PublicID:       omit.From(public_id),
		SentAt:         omitnull.FromPtr[time.Time](nil),
		Source:         omit.From(source),
		Subject:        omit.From(subject),
		TemplateID:     omit.From(template_id),
		TemplateData:   omit.From(data_for_insert),
		Type:           omit.From(type_),
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
