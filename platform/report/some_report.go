package report

import (
	"context"
	//"crypto/rand"
	//"fmt"
	//"math/big"
	//"strconv"
	//"strings"
	//"time"

	//"github.com/aarondl/opt/omit"
	//"github.com/aarondl/opt/omitnull"
	"github.com/Gleipnir-Technology/bob"
	//"github.com/Gleipnir-Technology/bob/dialect/psql"
	//"github.com/Gleipnir-Technology/bob/dialect/psql/sm"
	//"github.com/Gleipnir-Technology/bob/dialect/psql/um"
	//"github.com/Gleipnir-Technology/nidus-sync/background"
	//"github.com/Gleipnir-Technology/nidus-sync/db"
	//"github.com/Gleipnir-Technology/nidus-sync/db/models"
	//"github.com/Gleipnir-Technology/nidus-sync/db/sql"
	"github.com/Gleipnir-Technology/nidus-sync/platform/text"
	//"github.com/rs/zerolog/log"
	//"github.com/stephenafamo/scan"
)

type SomeReport interface {
	addNotificationEmail(context.Context, bob.Tx, string) *ErrorWithCode
	addNotificationPhone(context.Context, bob.Tx, text.E164) *ErrorWithCode
	districtID(context.Context) *int32
	updateReporterConsent(context.Context, bob.Tx, bool) *ErrorWithCode
	updateReporterEmail(context.Context, bob.Tx, string) *ErrorWithCode
	updateReporterName(context.Context, bob.Tx, string) *ErrorWithCode
	updateReporterPhone(context.Context, bob.Tx, text.E164) *ErrorWithCode
	PublicReportID() string
	reportID() int32
}
