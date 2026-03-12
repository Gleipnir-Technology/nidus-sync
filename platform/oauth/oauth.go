package oauth

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/Gleipnir-Technology/arcgis-go"
	"github.com/Gleipnir-Technology/nidus-sync/config"
	"github.com/Gleipnir-Technology/nidus-sync/db"
	"github.com/Gleipnir-Technology/nidus-sync/db/models"
	"github.com/aarondl/opt/omit"
	"github.com/rs/zerolog/log"
)

// When the API responds that the token is now invalidated
type InvalidatedTokenError struct{}

func (e InvalidatedTokenError) Error() string { return "The token has been invalidated by the server" }

type OAuthTokenResponse struct {
	AccessToken           string `json:"access_token"`
	ExpiresIn             int    `json:"expires_in"`
	RefreshToken          string `json:"refresh_token"`
	RefreshTokenExpiresIn int    `json:"refresh_token_expires_in"`
	SSL                   bool   `json:"ssl"`
	Username              string `json:"username"`
}

func DoTokenRequest(ctx context.Context, form url.Values) (*OAuthTokenResponse, error) {
	form.Set("client_id", config.ClientID)

	baseURL := "https://www.arcgis.com/sharing/rest/oauth2/token/"
	req, err := http.NewRequest("POST", baseURL, strings.NewReader(form.Encode()))
	if err != nil {
		return nil, fmt.Errorf("Failed to create request: %w", err)
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	client := http.Client{}
	log.Info().Str("url", req.URL.String()).Msg("POST")
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("Failed to do request: %w", err)
	}
	defer resp.Body.Close()
	bodyBytes, err := io.ReadAll(resp.Body)
	log.Info().Int("status", resp.StatusCode).Msg("Token request")
	if resp.StatusCode >= http.StatusBadRequest {
		if err != nil {
			return nil, fmt.Errorf("Got status code %d and failed to read response body: %w", resp.StatusCode, err)
		}
		bodyString := string(bodyBytes)
		var errorResp arcgis.ErrorResponse
		if err := json.Unmarshal(bodyBytes, &errorResp); err == nil {
			if errorResp.Error.Code == 498 && errorResp.Error.Description == "invalidated refresh_token" {
				return nil, InvalidatedTokenError{}
			}
			return nil, fmt.Errorf("API response JSON error: %d: %d %s", resp.StatusCode, errorResp.Error.Code, errorResp.Error.Description)
		}
		return nil, fmt.Errorf("API returned error status %d: %s", resp.StatusCode, bodyString)
	}
	//logResponseHeaders(resp)
	var tokenResponse OAuthTokenResponse
	err = json.Unmarshal(bodyBytes, &tokenResponse)
	if err != nil {
		return nil, fmt.Errorf("Failed to unmarshal JSON: %w", err)
	}
	// Just because we got a 200-level status code doesn't mean it worked. Experience has taught us that
	// we can get errors without anything indicated in the headers or the status code
	if tokenResponse == (OAuthTokenResponse{}) {
		var errorResponse arcgis.ErrorResponse
		err = json.Unmarshal(bodyBytes, &errorResponse)
		if err != nil {
			return nil, fmt.Errorf("Failed to unmarshal error JSON: %w", err)
		}
		if errorResponse.Error.Code > 0 {
			return nil, errorResponse.AsError(ctx)
		}
	}
	log.Info().Str("refresh token", tokenResponse.RefreshToken).Str("access token", tokenResponse.AccessToken).Int("access expires", tokenResponse.ExpiresIn).Int("refresh expires", tokenResponse.RefreshTokenExpiresIn).Msg("Oauth token acquired")
	return &tokenResponse, nil
}

func FutureUTCTimestamp(secondsFromNow int) time.Time {
	return time.Now().UTC().Add(time.Duration(secondsFromNow) * time.Second)
}

func GetOAuthForOrg(ctx context.Context, org *models.Organization) (*models.ArcgisOauthToken, error) {
	users, err := org.User().All(ctx, db.PGInstance.BobDB)
	if err != nil {
		return nil, fmt.Errorf("Failed to query all users for org: %w", err)
	}
	for _, user := range users {
		oauths, err := user.UserOauthTokens(models.SelectWhere.ArcgisOauthTokens.InvalidatedAt.IsNull()).All(ctx, db.PGInstance.BobDB)
		if err != nil {
			return nil, fmt.Errorf("Failed to query all oauth tokens for org: %w", err)
		}
		for _, oauth := range oauths {
			return oauth, nil
		}
	}
	return nil, nil
}

// Update the access token to keep it fresh and alive
func RefreshAccessToken(ctx context.Context, oauth *models.ArcgisOauthToken) error {
	form := url.Values{
		"grant_type":    []string{"refresh_token"},
		"client_id":     []string{config.ClientID},
		"refresh_token": []string{oauth.RefreshToken},
	}
	token, err := DoTokenRequest(ctx, form)
	if err != nil {
		return fmt.Errorf("Failed to handle request: %w", err)
	}
	accessExpires := FutureUTCTimestamp(token.ExpiresIn)
	setter := models.ArcgisOauthTokenSetter{
		AccessToken:        omit.From(token.AccessToken),
		AccessTokenExpires: omit.From(accessExpires),
		Username:           omit.From(token.Username),
	}
	err = oauth.Update(ctx, db.PGInstance.BobDB, &setter)
	if err != nil {
		return fmt.Errorf("Failed to update oauth in database: %w", err)
	}
	log.Info().Int("oauth token id", int(oauth.ID)).Msg("Updated oauth token")
	return nil
}

// Update the refresh token to keep it fresh and alive
func RefreshRefreshToken(ctx context.Context, oauth *models.ArcgisOauthToken) error {

	form := url.Values{
		"grant_type":    []string{"exchange_refresh_token"},
		"redirect_uri":  []string{config.ArcGISOauthRedirectURL()},
		"refresh_token": []string{oauth.RefreshToken},
	}

	token, err := DoTokenRequest(ctx, form)
	if err != nil {
		return fmt.Errorf("Failed to handle request: %w", err)
	}
	refreshExpires := FutureUTCTimestamp(token.ExpiresIn)
	setter := models.ArcgisOauthTokenSetter{
		RefreshToken:        omit.From(token.RefreshToken),
		RefreshTokenExpires: omit.From(refreshExpires),
		Username:            omit.From(token.Username),
	}
	err = oauth.Update(ctx, db.PGInstance.BobDB, &setter)
	if err != nil {
		return fmt.Errorf("Failed to update oauth in database: %w", err)
	}
	log.Info().Int("oauth token id", int(oauth.ID)).Msg("Updated oauth token")
	return nil
}
