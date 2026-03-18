package text

import (
	"context"
	"fmt"

	"github.com/Gleipnir-Technology/bob"
	"github.com/Gleipnir-Technology/bob/dialect/psql"
	"github.com/Gleipnir-Technology/bob/dialect/psql/im"
	"github.com/Gleipnir-Technology/nidus-sync/db"
	"github.com/Gleipnir-Technology/nidus-sync/db/enums"
	"github.com/Gleipnir-Technology/nidus-sync/db/models"
	"github.com/Gleipnir-Technology/nidus-sync/platform/types"
	"github.com/aarondl/opt/omit"
	"github.com/rs/zerolog/log"
)

func EnsureInDB(ctx context.Context, txn bob.Executor, dst types.E164) (err error) {
	return ensureInDB(ctx, txn, dst.PhoneString())
}
func ensureInDB(ctx context.Context, txn bob.Executor, destination string) (err error) {
	_, err = psql.Insert(
		im.Into("comms.phone", "e164", "is_subscribed", "status"),
		im.Values(
			psql.Arg(destination),
			psql.Arg(false),
			psql.Arg("unconfirmed"),
		),
		im.OnConflict("e164").DoNothing(),
	).Exec(ctx, txn)
	return err
}
func phoneStatus(ctx context.Context, src types.E164) (enums.CommsPhonestatustype, error) {
	phone, err := models.FindCommsPhone(ctx, db.PGInstance.BobDB, src.PhoneString())
	if err != nil {
		return enums.CommsPhonestatustypeUnconfirmed, fmt.Errorf("Failed to determine if '%s' is subscribed: %w", src.PhoneString(), err)
	}
	return phone.Status, nil
}
func setPhoneStatus(ctx context.Context, txn bob.Executor, src types.E164, status enums.CommsPhonestatustype) error {
	phone, err := models.FindCommsPhone(ctx, txn, src.PhoneString())
	if err != nil {
		return fmt.Errorf("Failed to determine if '%s' is subscribed: %w", src, err)
	}
	phone.Update(ctx, txn, &models.CommsPhoneSetter{
		Status: omit.From(status),
	})
	log.Info().Str("src", src.PhoneString()).Str("status", string(status)).Msg("Set number subscribed")
	return nil
}
