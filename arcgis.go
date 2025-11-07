package main

import (
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
	"strconv"
	"strings"
	"time"

	"github.com/Gleipnir-Technology/arcgis-go"
	"github.com/Gleipnir-Technology/nidus-sync/models"
	"github.com/Gleipnir-Technology/nidus-sync/sql"
	"github.com/aarondl/opt/omit"
	"github.com/aarondl/opt/omitnull"
	"github.com/jackc/pgx/v5"
)

var CodeVerifier string = "random_secure_string_min_43_chars_long_should_be_stored_in_session"

type OAuthTokenResponse struct {
	AccessToken  string `json:"access_token"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
	Username     string `json:"username"`
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
func updateArcgisUserData(ctx context.Context, user *models.User, access_token string, expires time.Time, refresh_token string) {
	client := arcgis.NewArcGIS(
		arcgis.AuthenticatorOAuth{
			AccessToken:  access_token,
			Expires:      expires,
			RefreshToken: refresh_token,
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
	client := http.Client{}
	slog.Info("POST", slog.String("url", baseURL))
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("Failed to do request: %v", err)
	}
	defer resp.Body.Close()
	bodyBytes, err := io.ReadAll(resp.Body)
	slog.Info("Response", slog.Int("status", resp.StatusCode))
	if resp.StatusCode >= http.StatusBadRequest {
		if err != nil {
			return fmt.Errorf("Got status code %d and failed to read response body: %v", resp.StatusCode, err)
		}
		bodyString := string(bodyBytes)
		var errorResp map[string]interface{}
		if err := json.Unmarshal(bodyBytes, &errorResp); err == nil {
			return fmt.Errorf("API response JSON error: %d: %v", resp.StatusCode, errorResp)
		}
		return fmt.Errorf("API returned error status %d: %s", resp.StatusCode, bodyString)
	}
	var tokenResponse OAuthTokenResponse
	err = json.Unmarshal(bodyBytes, &tokenResponse)
	if err != nil {
		return fmt.Errorf("Failed to unmarshal JSON: %v", err)
	}
	slog.Info("Oauth token acquired",
		slog.String("refresh token", tokenResponse.RefreshToken),
		slog.String("access token", tokenResponse.AccessToken),
		slog.Int("expires", tokenResponse.ExpiresIn),
	)

	expires := futureUTCTimestamp(tokenResponse.ExpiresIn)
	setter := models.OauthTokenSetter{
		AccessToken:  omit.From(tokenResponse.AccessToken),
		Expires:      omit.From(expires),
		RefreshToken: omit.From(tokenResponse.RefreshToken),
		Username:     omit.From(tokenResponse.Username),
	}
	err = user.InsertUserOauthTokens(ctx, PGInstance.BobDB, &setter)
	if err != nil {
		return fmt.Errorf("Failed to save token to database: %v", err)
	}
	go updateArcgisUserData(context.Background(), user, tokenResponse.AccessToken, expires, tokenResponse.RefreshToken)
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
func refreshFieldseekerData(ctx context.Context, newOauthCh <-chan int) {
	for {
		select {
		case <-ctx.Done():
			slog.Info("Exiting refresh worker")
			return
		case id := <-newOauthCh:
			slog.Info("Adding oauth to background work", slog.Int("oauth id", id))
		}
	}
}
