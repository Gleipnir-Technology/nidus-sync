package main

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/aarondl/opt/omit"
	"github.com/Gleipnir-Technology/nidus-sync/models"
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
	log.Printf("POST %s", baseURL)
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("Failed to do request: %v", err)
	}
	defer resp.Body.Close()
	bodyBytes, err := io.ReadAll(resp.Body)
	log.Printf("Response %d", resp.StatusCode)
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
	log.Printf("Refresh token '%s'", tokenResponse.RefreshToken)

	setter := models.OauthTokenSetter{
		AccessToken: omit.From(tokenResponse.AccessToken),
		Expires: omit.From(futureUTCTimestamp(tokenResponse.ExpiresIn)),
		RefreshToken: omit.From(tokenResponse.RefreshToken),
		Username: omit.From(tokenResponse.Username),
	}
	err = user.InsertUserOauthTokens(ctx, PGInstance.BobDB, &setter)
	if err != nil {
		return fmt.Errorf("Failed to save token to database: %v", err)
	}
	return nil
}

func redirectURL() string {
	return BaseURL + "/arcgis/oauth/callback"
}
