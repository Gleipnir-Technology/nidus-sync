package platform

import (
	"context"
	"fmt"
	"net/url"
	"time"

	"github.com/Gleipnir-Technology/nidus-sync/config"
	"github.com/Gleipnir-Technology/nidus-sync/db/gen/nidus-sync/arcgis/model"
	queryarcgis "github.com/Gleipnir-Technology/nidus-sync/db/query/arcgis"
	"github.com/Gleipnir-Technology/nidus-sync/platform/oauth"
)

// When there is no oauth for an organization
type NoOAuthForOrg struct{}

func (e NoOAuthForOrg) Error() string { return "No oauth available for organization" }

func GetOAuthForOrg(ctx context.Context, org Organization) (*model.OAuthToken, error) {
	result, err := oauth.GetOAuthForOrg(ctx, org.model)
	if result == nil && err == nil {
		return nil, &NoOAuthForOrg{}
	}
	return result, err
}

func GetOAuthForUser(ctx context.Context, user User) (*model.OAuthToken, error) {
	oauth, err := queryarcgis.OAuthTokenForUser(ctx, int64(user.ID))
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
	setter := model.OAuthToken{
		AccessToken:         token.AccessToken,
		AccessTokenExpires:  accessExpires,
		ArcgisAccountID:     nil,
		ArcgisID:            nil,
		ArcgisLicenseTypeID: nil,
		Created:             time.Now(),
		InvalidatedAt:       nil,
		RefreshToken:        token.RefreshToken,
		RefreshTokenExpires: refreshExpires,
		UserID:              int32(user.ID),
		Username:            token.Username,
	}
	oauth, err := queryarcgis.OAuthTokenInsert(ctx, &setter)
	if err != nil {
		return fmt.Errorf("Failed to save token to database: %w", err)
	}
	go updateArcgisUserData(context.Background(), user.model, oauth)
	return nil
}
