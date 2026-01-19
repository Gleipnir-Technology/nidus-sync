package publicreport

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"net/http"
	"strings"

	"github.com/Gleipnir-Technology/nidus-sync/comms"
	"github.com/go-chi/chi/v5"
)

// GenerateReportID creates a 12-character random string using only unambiguous
// capital letters and numbers
func GenerateReportID() (string, error) {
	// Define character set (no O/0, I/l/1, 2/Z to avoid confusion)
	const charset = "ABCDEFGHJKLMNPQRSTUVWXY3456789"
	const length = 12

	var builder strings.Builder
	builder.Grow(length)

	// Use crypto/rand for secure randomness
	for i := 0; i < length; i++ {
		// Generate a random index within our charset
		n, err := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		if err != nil {
			return "", fmt.Errorf("failed to generate random number: %w", err)
		}

		// Add the randomly selected character to our ID
		builder.WriteByte(charset[n.Int64()])
	}

	return builder.String(), nil
}

func getEmailReportSubscriptionConfirmation(w http.ResponseWriter, r *http.Request) {
	report_id := chi.URLParam(r, "report_id")
	comms.RenderEmailReportConfirmation(w, report_id)
}
