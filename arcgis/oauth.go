package arcgis

import (
	"context"

	"github.com/Gleipnir-Technology/bob/dialect/psql/sm"
	"github.com/Gleipnir-Technology/nidus-sync/db"
	"github.com/Gleipnir-Technology/nidus-sync/db/models"
)

func GetOAuthForUser(ctx context.Context, user *models.User) (*models.ArcgisOauthToken, error) {
	oauth, err := user.UserOauthTokens(
		sm.OrderBy("created").Desc(),
	).One(ctx, db.PGInstance.BobDB)
	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			return nil, nil
		}
		return nil, err
	}
	return oauth, nil
}
