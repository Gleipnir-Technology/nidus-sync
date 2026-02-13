package arcgis

import (
	"context"

	"github.com/Gleipnir-Technology/bob/dialect/psql/sm"
	"github.com/Gleipnir-Technology/nidus-sync/db"
	"github.com/Gleipnir-Technology/nidus-sync/db/models"
)

func GetOAuthForUser(ctx context.Context, user *models.User) (*models.OauthToken, error) {
	return user.UserOauthTokens(
		sm.OrderBy("created").Desc(),
	).One(ctx, db.PGInstance.BobDB)
}
