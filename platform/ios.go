package platform

import (
	"context"
	"time"

	"github.com/Gleipnir-Technology/nidus-sync/db"
	"github.com/Gleipnir-Technology/nidus-sync/db/models"
)

func fieldseeker(ctx context.Context, u *models.User, since *time.Time) (fsync FieldseekerRecordsSync, err error) {
	pl, err := u.R.Organization.Pointlocations().All(ctx, db.PGInstance.BobDB)
	if err != nil {
		return
	}
	fsync.MosquitoSources = pl
	return fsync, err
}

func ContentClientIos(ctx context.Context, u *models.User, since *time.Time) (csync ClientSync, err error) {
	fsync, err := fieldseeker(ctx, u, since)
	return ClientSync{
		Fieldseeker: fsync,
	}, err
}
