package main

import (
	"bytes"
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/Gleipnir-Technology/arcgis-go"
	"github.com/Gleipnir-Technology/nidus-sync/models"
	"github.com/Gleipnir-Technology/nidus-sync/sql"
	"github.com/aarondl/opt/omit"
	"github.com/aarondl/opt/omitnull"
	"github.com/jackc/pgx/v5"
)

var NewOAuthTokenChannel chan struct{}
var CodeVerifier string = "random_secure_string_min_43_chars_long_should_be_stored_in_session"

type ErrorResponse struct {
	Error ErrorResponseContent `json:"error"`
}

type ErrorResponseContent struct {
	Code             int      `json:"code"`
	Error            string   `json:"error"`
	ErrorDescription string   `json:"error_description"`
	Message          string   `json:"message"`
	Details          []string `json:"details"`
}

type OAuthTokenResponse struct {
	AccessToken           string `json:"access_token"`
	ExpiresIn             int    `json:"expires_in"`
	RefreshToken          string `json:"refresh_token"`
	RefreshTokenExpiresIn int    `json:"refresh_token_expires_in"`
	SSL                   bool   `json:"ssl"`
	Username              string `json:"username"`
}

// Build the ArcGIS authorization URL with PKCE
func buildArcGISAuthURL(clientID string, redirectURI string, expiration int) string {
	baseURL := "https://www.arcgis.com/sharing/rest/oauth2/authorize/"

	params := url.Values{}
	params.Add("client_id", clientID)
	params.Add("redirect_uri", redirectURI)
	params.Add("response_type", "code")
	//params.Add("code_challenge", generateCodeChallenge(codeVerifier))
	//params.Add("code_challenge_method", "S256")
	params.Add("expiration", strconv.Itoa(expiration))

	return baseURL + "?" + params.Encode()
}

func futureUTCTimestamp(secondsFromNow int) time.Time {
	return time.Now().UTC().Add(time.Duration(secondsFromNow) * time.Second)
}

// Helper function to generate code challenge from code verifier
func generateCodeChallenge(codeVerifier string) string {
	hash := sha256.Sum256([]byte(codeVerifier))
	return base64.RawURLEncoding.EncodeToString(hash[:])
}

// Generate a random code verifier for PKCE
func generateCodeVerifier() string {
	bytes := make([]byte, 64) // 64 bytes = 512 bits
	rand.Read(bytes)
	return base64.RawURLEncoding.EncodeToString(bytes)
}

// Find out what we can about this user
func updateArcgisUserData(ctx context.Context, user *models.User, access_token string, access_token_expires time.Time, refresh_token string, refresh_token_expires time.Time) {
	client := arcgis.NewArcGIS(
		arcgis.AuthenticatorOAuth{
			AccessToken:         access_token,
			AccessTokenExpires:  access_token_expires,
			RefreshToken:        refresh_token,
			RefreshTokenExpires: refresh_token_expires,
		},
	)
	portal, err := client.PortalsSelf()
	if err != nil {
		slog.Error("Failed to get ArcGIS user data", slog.String("err", err.Error()))
		return
	}
	slog.Info("Got portals data",
		slog.String("Username", portal.User.Username),
		slog.String("user_id", portal.User.ID),
		slog.String("org_id", portal.User.OrgID),
		slog.String("org_name", portal.Name),
		slog.String("license_type_id", portal.User.UserLicenseTypeID))

	_, err = sql.UpdateOauthTokenOrg(portal.User.ID, portal.User.UserLicenseTypeID, refresh_token).Exec(ctx, PGInstance.BobDB)
	if err != nil {
		slog.Error("Failed to update oauth token portal data", slog.String("err", err.Error()))
		return
	}
	var org *models.Organization
	orgs, err := models.Organizations.Query(models.SelectWhere.Organizations.ArcgisName.EQ(portal.Name)).All(ctx, PGInstance.BobDB)
	switch len(orgs) {
	case 0:
		setter := models.OrganizationSetter{
			Name:       omitnull.From(portal.Name),
			ArcgisID:   omitnull.From(portal.User.OrgID),
			ArcgisName: omitnull.From(portal.Name),
		}
		org, err = models.Organizations.Insert(&setter).One(ctx, PGInstance.BobDB)
		if err != nil {
			slog.Error("Failed to create new organization", slog.String("err", err.Error()))
			return
		}
		slog.Info("Created new organization", slog.Int("org_id", int(org.ID)))
	case 1:
		org = orgs[0]
		slog.Info("Organization already exists")
	default:
		slog.Error("Got too many organizations, bailing")
		return

	}
	if err != nil {
		LogErrorTypeInfo(err)
		if errors.Is(err, pgx.ErrNoRows) {
		} else {
			slog.Error("Failed to query for existing org", slog.String("err", err.Error()))
			return
		}
	}
	err = org.AttachUser(ctx, PGInstance.BobDB, user)
	if err != nil {
		slog.Error("Failed to attach user to organization", slog.String("err", err.Error()), slog.Int("user_id", int(user.ID)), slog.Int("org_id", int(org.ID)))
		return
	}

	search, err := client.Search("Fieldseeker")
	if err != nil {
		slog.Error("Failed to get search FieldseekerGIS data", slog.String("err", err.Error()))
		return
	}
	for _, result := range search.Results {
		slog.Info("Got result", slog.String("name", result.Name))
		//if result.Name == "FieldseekerGIS" {
		//slog.Info("Found Fieldseeker", slog.String("url", result.URL))
		//}
	}
}

func handleOauthAccessCode(ctx context.Context, user *models.User, code string) error {
	baseURL := "https://www.arcgis.com/sharing/rest/oauth2/token/"

	//params.Add("code_verifier", "S256")

	form := url.Values{
		"grant_type":   []string{"authorization_code"},
		"code":         []string{code},
		"client_id":    []string{ClientID},
		"redirect_uri": []string{redirectURL()},
	}

	req, err := http.NewRequest("POST", baseURL, strings.NewReader(form.Encode()))
	if err != nil {
		return fmt.Errorf("Failed to create request: %v", err)
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	token, err := handleTokenRequest(ctx, req)
	if err != nil {
		return fmt.Errorf("Failed to exchange authorization code for token: %v", err)
	}
	accessExpires := futureUTCTimestamp(token.ExpiresIn)
	refreshExpires := futureUTCTimestamp(token.RefreshTokenExpiresIn)
	setter := models.OauthTokenSetter{
		AccessToken:         omit.From(token.AccessToken),
		AccessTokenExpires:  omit.From(accessExpires),
		RefreshToken:        omit.From(token.RefreshToken),
		RefreshTokenExpires: omit.From(refreshExpires),
		Username:            omit.From(token.Username),
	}
	err = user.InsertUserOauthTokens(ctx, PGInstance.BobDB, &setter)
	if err != nil {
		return fmt.Errorf("Failed to save token to database: %v", err)
	}
	go updateArcgisUserData(context.Background(), user, token.AccessToken, accessExpires, token.RefreshToken, refreshExpires)
	NewOAuthTokenChannel <- struct{}{}
	return nil
}

func hasFieldseekerConnection(ctx context.Context, user *models.User) (bool, error) {
	result, err := sql.OauthTokenByUserId(user.ID).All(ctx, PGInstance.BobDB)
	if err != nil {
		return false, err
	}
	return len(result) > 0, nil
}
func redirectURL() string {
	return BaseURL + "/arcgis/oauth/callback"
}

// This is a goroutine that is in charge of getting Fieldseeker data and keeping it fresh.
func refreshFieldseekerData(ctx context.Context, newOauthCh <-chan struct{}) {
	for {
		workerCtx, cancel := context.WithCancel(context.Background())
		var wg sync.WaitGroup

		oauths, err := models.OauthTokens.Query().All(ctx, PGInstance.BobDB)
		if err != nil {
			slog.Error("Failed to get oauths", slog.String("err", err.Error()))
			return
		}
		for _, oauth := range oauths {
			wg.Add(1)
			go func() {
				defer wg.Done()
				maintainOAuth(workerCtx, oauth)
			}()
		}
		select {
		case <-ctx.Done():
			slog.Info("Exiting refresh worker...")
			cancel()
			wg.Wait()
			return
		case <-newOauthCh:
			slog.Info("Updating oauth background work")
			cancel()
			wg.Wait()
		}
	}
}

func maintainOAuth(ctx context.Context, oauth *models.OauthToken) {
	refreshDelay := time.Until(oauth.AccessTokenExpires)
	slog.Info("Need to refresh oauth", slog.Int("id", int(oauth.ID)), slog.Float64("seconds", refreshDelay.Seconds()))
	if oauth.AccessTokenExpires.Before(time.Now()) {
		err := refreshOAuth(ctx, oauth)
		if err != nil {
			slog.Error("Failed to refresh token", slog.String("err", err.Error()))
			return
		}
	}
	ticker := time.NewTicker(refreshDelay)
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:

		}
	}

}

func refreshOAuth(ctx context.Context, oauth *models.OauthToken) error {
	baseURL := "https://www.arcgis.com/sharing/rest/oauth2/token/"

	form := url.Values{
		"grant_type":    []string{"refresh_token"},
		"client_id":     []string{ClientID},
		"refresh_token": []string{oauth.RefreshToken},
	}

	req, err := http.NewRequest("POST", baseURL, strings.NewReader(form.Encode()))
	if err != nil {
		return fmt.Errorf("Failed to create request: %v", err)
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	token, err := handleTokenRequest(ctx, req)
	if err != nil {
		return fmt.Errorf("Failed to handle request: %v", err)
	}
	accessExpires := futureUTCTimestamp(token.ExpiresIn)
	refreshExpires := futureUTCTimestamp(token.RefreshTokenExpiresIn)
	setter := models.OauthTokenSetter{
		AccessToken:         omit.From(token.AccessToken),
		AccessTokenExpires:  omit.From(accessExpires),
		RefreshToken:        omit.From(token.RefreshToken),
		RefreshTokenExpires: omit.From(refreshExpires),
		Username:            omit.From(token.Username),
	}
	err = oauth.Update(ctx, PGInstance.BobDB, &setter)
	if err != nil {
		return fmt.Errorf("Failed to update oauth in database: %v", err)
	}
	return nil
}

func handleTokenRequest(ctx context.Context, req *http.Request) (*OAuthTokenResponse, error) {
	client := http.Client{}
	slog.Info("POST", slog.String("url", req.URL.String()))
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("Failed to do request: %v", err)
	}
	defer resp.Body.Close()
	bodyBytes, err := io.ReadAll(resp.Body)
	slog.Info("Token request", slog.Int("status", resp.StatusCode))
	saveResponse(bodyBytes, "token.json")
	if resp.StatusCode >= http.StatusBadRequest {
		if err != nil {
			return nil, fmt.Errorf("Got status code %d and failed to read response body: %v", resp.StatusCode, err)
		}
		bodyString := string(bodyBytes)
		var errorResp map[string]interface{}
		if err := json.Unmarshal(bodyBytes, &errorResp); err == nil {
			return nil, fmt.Errorf("API response JSON error: %d: %v", resp.StatusCode, errorResp)
		}
		return nil, fmt.Errorf("API returned error status %d: %s", resp.StatusCode, bodyString)
	}
	//logResponseHeaders(resp)
	var tokenResponse OAuthTokenResponse
	err = json.Unmarshal(bodyBytes, &tokenResponse)
	if err != nil {
		return nil, fmt.Errorf("Failed to unmarshal JSON: %v", err)
	}
	// Just because we got a 200-level status code doesn't mean it worked. Experience has taught us that
	// we can get errors without anything indicated in the headers or the status code
	if tokenResponse == (OAuthTokenResponse{}) {
		var errorResponse ErrorResponse
		err = json.Unmarshal(bodyBytes, &errorResponse)
		if err != nil {
			return nil, fmt.Errorf("Failed to unmarshal error JSON: %v", err)
		}
		if errorResponse.Error.Code > 0 {
			return nil, errors.New(fmt.Sprintf("API error %d: %s: %s (%s)",
				errorResponse.Error.Code,
				errorResponse.Error.Error,
				errorResponse.Error.ErrorDescription,
				errorResponse.Error.Message,
			))
		}
	}
	slog.Info("Oauth token acquired",
		slog.String("refresh token", tokenResponse.RefreshToken),
		slog.String("access token", tokenResponse.AccessToken),
		slog.Int("access expires", tokenResponse.ExpiresIn),
		slog.Int("refresh expires", tokenResponse.RefreshTokenExpiresIn),
	)
	return &tokenResponse, nil
}

func logResponseHeaders(resp *http.Response) {
	if resp == nil {
		slog.Info("Response is nil")
		return
	}

	slog.Info("HTTP Response headers",
		"status", resp.Status,
		"statusCode", resp.StatusCode)

	for name, values := range resp.Header {
		slog.Info("Header",
			"name", name,
			"values", values)
	}
}

func saveResponse(data []byte, filename string) {
	dest, err := os.Create(filename)
	if err != nil {
		slog.Error("Failed to create file", slog.String("filename", filename), slog.String("err", err.Error()))
		return
	}
	_, err = io.Copy(dest, bytes.NewReader(data))
	if err != nil {
		slog.Error("Failed to write", slog.String("filename", filename), slog.String("err", err.Error()))
		return
	}
	slog.Info("Wrote response", slog.String("filename", filename))
}
