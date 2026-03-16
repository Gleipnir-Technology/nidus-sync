package platform

import (
	"context"
	"fmt"
	"net/url"
	"time"

	"github.com/Gleipnir-Technology/bob/dialect/psql/sm"
	"github.com/Gleipnir-Technology/nidus-sync/config"
	"github.com/Gleipnir-Technology/nidus-sync/db"
	"github.com/Gleipnir-Technology/nidus-sync/db/models"
	"github.com/Gleipnir-Technology/nidus-sync/platform/oauth"
	"github.com/aarondl/opt/omit"
	"github.com/aarondl/opt/omitnull"
)

// When there is no oauth for an organization
type NoOAuthForOrg struct{}

func (e NoOAuthForOrg) Error() string { return "No oauth available for organization" }

func GetOAuthForOrg(ctx context.Context, org Organization) (*models.ArcgisOauthToken, error) {
	result, err := oauth.GetOAuthForOrg(ctx, org.model)
	if result == nil && err == nil {
		return nil, &NoOAuthForOrg{}
	}
	return result, err
}

func GetOAuthForUser(ctx context.Context, user User) (*models.ArcgisOauthToken, error) {
	oauth, err := user.model.UserOauthTokens(
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

func HandleOauthAccessCode(ctx context.Context, user User, code string) error {
	form := url.Values{
		"grant_type":   []string{"authorization_code"},
		"code":         []string{code},
		"redirect_uri": []string{config.ArcGISOauthRedirectURL()},
	}

	token, err := oauth.DoTokenRequest(ctx, form)
	if err != nil {
		return fmt.Errorf("Failed to exchange authorization code for token: %w", err)
	}
	accessExpires := oauth.FutureUTCTimestamp(token.ExpiresIn)
	refreshExpires := oauth.FutureUTCTimestamp(token.RefreshTokenExpiresIn)
	setter := models.ArcgisOauthTokenSetter{
		AccessToken:        omit.From(token.AccessToken),
		AccessTokenExpires: omit.From(accessExpires),
		//ArcgisAccountID:     omit.From(
		ArcgisID:            omitnull.FromPtr[string](nil),
		ArcgisLicenseTypeID: omitnull.FromPtr[string](nil),
		Created:             omit.From(time.Now()),
		InvalidatedAt:       omitnull.FromPtr[time.Time](nil),
		RefreshToken:        omit.From(token.RefreshToken),
		RefreshTokenExpires: omit.From(refreshExpires),
		UserID:              omit.From(int32(user.ID)),
		Username:            omit.From(token.Username),
	}
	oauth, err := models.ArcgisOauthTokens.Insert(&setter).One(ctx, db.PGInstance.BobDB)
	if err != nil {
		return fmt.Errorf("Failed to save token to database: %w", err)
	}
	go updateArcgisUserData(context.Background(), user.model, oauth)
	return nil
}
