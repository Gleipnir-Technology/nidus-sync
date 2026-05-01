package arcgis

import (
	"context"

	"github.com/Gleipnir-Technology/nidus-sync/db"
	"github.com/Gleipnir-Technology/nidus-sync/db/gen/nidus-sync/arcgis/model"
	"github.com/Gleipnir-Technology/nidus-sync/db/gen/nidus-sync/arcgis/table"
	"github.com/go-jet/jet/v2/postgres"
)

func OAuthTokenInsert(ctx context.Context, m *model.OAuthToken) (*model.OAuthToken, error) {
	statement := table.OAuthToken.INSERT(table.OAuthToken.MutableColumns).
		MODEL(m)
	return db.ExecuteOne[model.OAuthToken](ctx, statement)
}
func OAuthTokenInvalidate(ctx context.Context, id int64) error {
	statement := table.OAuthToken.UPDATE().
		SET(table.OAuthToken.InvalidatedAt.SET(postgres.LOCALTIMESTAMP())).
		WHERE(table.OAuthToken.ID.EQ(postgres.Int(id)))
	return db.ExecuteNone(ctx, statement)
}
func OAuthTokensValid(ctx context.Context) ([]*model.OAuthToken, error) {
	statement := table.OAuthToken.SELECT(table.OAuthToken.AllColumns).
		FROM(table.OAuthToken).
		WHERE(table.OAuthToken.InvalidatedAt.IS_NULL())
	return db.ExecuteMany[model.OAuthToken](ctx, statement)
}
func OAuthTokenFromID(ctx context.Context, id int64) (*model.OAuthToken, error) {
	statement := table.OAuthToken.SELECT(
		table.OAuthToken.AllColumns,
	).FROM(table.OAuthToken).
		WHERE(table.OAuthToken.ID.EQ(postgres.Int(id)))
	return db.ExecuteOne[model.OAuthToken](ctx, statement)
}
func OAuthTokenForUser(ctx context.Context, user_id int64) (*model.OAuthToken, error) {
	statement := table.OAuthToken.SELECT(table.OAuthToken.AllColumns).
		FROM(table.OAuthToken).
		WHERE(table.OAuthToken.InvalidatedAt.IS_NULL().AND(
			table.OAuthToken.UserID.EQ(postgres.Int(user_id)),
		)).
		ORDER_BY(table.OAuthToken.Created.DESC()).
		LIMIT(1)
	return db.ExecuteOne[model.OAuthToken](ctx, statement)
}
func OAuthTokensForUser(ctx context.Context, user_id int64) ([]*model.OAuthToken, error) {
	statement := table.OAuthToken.SELECT(table.OAuthToken.AllColumns).
		FROM(table.OAuthToken).
		WHERE(table.OAuthToken.InvalidatedAt.IS_NULL().AND(
			table.OAuthToken.UserID.EQ(postgres.Int(user_id)),
		))
	return db.ExecuteMany[model.OAuthToken](ctx, statement)
}
func OAuthTokenForUserExists(ctx context.Context, user_id int64) (*bool, error) {
	statement := table.OAuthToken.SELECT(postgres.Bool(true)).
		FROM(table.OAuthToken).
		WHERE(table.OAuthToken.UserID.EQ(postgres.Int(user_id))).
		LIMIT(1)
	return db.ExecuteOne[bool](ctx, statement)
}
func OAuthTokenUpdateAccessToken(ctx context.Context, oauth_id int64, updates *model.OAuthToken) error {
	statement := table.OAuthToken.UPDATE(
		table.OAuthToken.AccessToken,
		table.OAuthToken.AccessTokenExpires,
		table.OAuthToken.Username,
	).MODEL(updates).
		WHERE(table.OAuthToken.ID.EQ(postgres.Int(oauth_id)))
	return db.ExecuteNone(ctx, statement)
}
func OAuthTokenUpdateRefreshToken(ctx context.Context, oauth_id int64, updates *model.OAuthToken) error {
	statement := table.OAuthToken.UPDATE(
		table.OAuthToken.RefreshToken,
		table.OAuthToken.RefreshTokenExpires,
		table.OAuthToken.Username,
	).MODEL(updates).
		WHERE(table.OAuthToken.ID.EQ(postgres.Int(oauth_id)))
	return db.ExecuteNone(ctx, statement)

}
func OAuthTokenUpdateLicense(ctx context.Context, refresh_token string, updates *model.OAuthToken) error {
	statement := table.OAuthToken.UPDATE(
		table.OAuthToken.ArcgisID,
		table.OAuthToken.ArcgisLicenseTypeID,
	).MODEL(updates).
		WHERE(table.OAuthToken.RefreshToken.EQ(postgres.String(refresh_token)))
	return db.ExecuteNone(ctx, statement)
}
