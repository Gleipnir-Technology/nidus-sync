package platform

import (
	"context"
	"fmt"
	"time"

	"github.com/Gleipnir-Technology/nidus-sync/db"
	"github.com/Gleipnir-Technology/nidus-sync/db/models"
	"github.com/Gleipnir-Technology/nidus-sync/platform/event"
	"github.com/aarondl/opt/omit"
	"github.com/aarondl/opt/omitnull"
	"github.com/rs/zerolog/log"
)

func ImpersonationCreate(ctx context.Context, user User, target int) (*models.LogImpersonation, error) {
	if !user.HasRoot() {
		return nil, fmt.Errorf("user %d is not root, and therefore can't impersonate user %d", user.ID, target)
	}
	setter := models.LogImpersonationSetter{
		BeginAt: omit.From(time.Now()),
		EndAt:   omitnull.FromPtr[time.Time](nil),
		//ID: ,
		ImpersonatorID: omit.From(int32(user.ID)),
		TargetID:       omit.From(int32(target)),
	}
	log, err := models.LogImpersonations.Insert(&setter).One(ctx, db.PGInstance.BobDB)
	if err != nil {
		return nil, fmt.Errorf("insert log: %w", err)
	}
	event.UpdatedUser(event.TypeSession, user.model.ID, "")
	event.UpdatedUser(event.TypeSession, int32(target), "")
	return log, nil
}
func ImpersonationEnd(ctx context.Context, user User, impersonator_id int32) error {
	l, err := models.LogImpersonations.Query(
		models.SelectWhere.LogImpersonations.EndAt.IsNull(),
		models.SelectWhere.LogImpersonations.ImpersonatorID.EQ(impersonator_id),
	).One(ctx, db.PGInstance.BobDB)
	if err != nil {
		return fmt.Errorf("query impersonations: %w", err)
	}
	err = l.Update(ctx, db.PGInstance.BobDB, &models.LogImpersonationSetter{
		EndAt: omitnull.From(time.Now()),
	})
	if err != nil {
		return fmt.Errorf("update impersonation log: %w", err)
	}
	log.Info().Int32("impersonator", l.ImpersonatorID).Int32("target", l.TargetID).Msg("Stopped impersonating")
	event.UpdatedUser(event.TypeSession, user.model.ID, "")
	event.UpdatedUser(event.TypeSession, impersonator_id, "")
	return nil
}
