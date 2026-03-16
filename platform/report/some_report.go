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
	"github.com/Gleipnir-Technology/nidus-sync/platform/types"
	//"github.com/rs/zerolog/log"
	//"github.com/stephenafamo/scan"
)

type SomeReport interface {
	addNotificationEmail(context.Context, bob.Executor, string) *ErrorWithCode
	addNotificationPhone(context.Context, bob.Executor, types.E164) *ErrorWithCode
	districtID(context.Context) *int32
	updateReporterConsent(context.Context, bob.Executor, bool) *ErrorWithCode
	updateReporterEmail(context.Context, bob.Executor, string) *ErrorWithCode
	updateReporterName(context.Context, bob.Executor, string) *ErrorWithCode
	updateReporterPhone(context.Context, bob.Executor, types.E164) *ErrorWithCode
	PublicReportID() string
	reportID() int32
}
