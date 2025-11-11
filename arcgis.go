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
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/Gleipnir-Technology/arcgis-go"
	"github.com/Gleipnir-Technology/arcgis-go/fieldseeker"
	enums "github.com/Gleipnir-Technology/nidus-sync/enums"
	"github.com/Gleipnir-Technology/nidus-sync/models"
	"github.com/Gleipnir-Technology/nidus-sync/sql"
	"github.com/aarondl/opt/omit"
	"github.com/aarondl/opt/omitnull"
	"github.com/alitto/pond/v2"
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
		if result.Name == "FieldSeekerGIS" {
			slog.Info("Found Fieldseeker", slog.String("url", result.URL))
			setter := models.OrganizationSetter{
				FieldseekerURL: omitnull.From(result.URL),
			}
			err = org.Update(ctx, PGInstance.BobDB, &setter)
			if err != nil {
				slog.Error("Failed to create new organization", slog.String("err", err.Error()))
				return
			}
		}
	}
	NewOAuthTokenChannel <- struct{}{}
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
				err := maintainOAuth(workerCtx, oauth)
				if err != nil {
					slog.Error("Crashed oauth maintenance goroutine", slog.String("err", err.Error()))
				}
			}()
		}

		orgs, err := models.Organizations.Query().All(ctx, PGInstance.BobDB)
		if err != nil {
			slog.Error("Failed to get orgs", slog.String("err", err.Error()))
			return
		}
		for _, org := range orgs {
			wg.Add(1)
			go func() {
				defer wg.Done()
				err := periodicallyExportFieldseeker(workerCtx, org)
				if err != nil {
					slog.Error("Crashed fieldseeker export goroutine", slog.String("err", err.Error()))
				}
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

type SyncStats struct {
	Inserts   int
	Updates   int
	Unchanged int
}

func downloadAllRecords(ctx context.Context, fssync *fieldseeker.FieldSeeker, layer arcgis.LayerFeature, org_id int32) (SyncStats, error) {
	var stats SyncStats
	count, err := fssync.QueryCount(layer.ID)
	if err != nil {
		return stats, fmt.Errorf("Failed to get counts for layer %s (%d): %v", layer.Name, layer.ID, err)
	}
	slog.Info("Starting on layer", slog.String("name", layer.Name), slog.Int("id", layer.ID))
	if count.Count == 0 {
		return stats, nil
	}
	pool := pond.NewResultPool[SyncStats](20)
	group := pool.NewGroup()
	maxRecords := fssync.MaxRecordCount()
	for offset := 0; offset < count.Count; offset += maxRecords {
		group.SubmitErr(func() (SyncStats, error) {
			query := arcgis.NewQuery()
			query.ResultRecordCount = maxRecords
			query.ResultOffset = offset
			query.SpatialReference = "4326"
			query.OutFields = "*"
			query.Where = "1=1"
			qr, err := fssync.DoQuery(
				layer.ID,
				query)
			if err != nil {
				return SyncStats{}, fmt.Errorf("Failed to get layer %s (%d) at offset %d: %v", layer.Name, layer.ID, offset, err)
			}
			i, u, err := saveOrUpdateDBRecords(ctx, "FS_"+layer.Name, qr, org_id)
			if err != nil {
				filename := fmt.Sprintf("failure-%s-%d-%d.json", layer.Name, layer.ID, offset)
				saveRawQuery(fssync, layer, query, filename)
				slog.Error("Faield to save DB records", slog.String("err", err.Error()))
				return SyncStats{}, fmt.Errorf("Failed to save records: %v", err)
			}
			return SyncStats{
				Inserts:   i,
				Updates:   u,
				Unchanged: len(qr.Features) - u - i,
			}, nil
		})
	}
	results, err := group.Wait()
	if err != nil {
		return stats, fmt.Errorf("one or more tasks in the work pool failed: %v", err)
	}
	for _, r := range results {
		stats.Inserts += r.Inserts
		stats.Updates += r.Updates
		stats.Unchanged += r.Unchanged
	}
	slog.Info("Finished layer", slog.Int("inserts", stats.Inserts), slog.Int("updates", stats.Updates), slog.Int("no change", stats.Unchanged))
	return stats, nil
}

func getOAuthForOrg(ctx context.Context, org *models.Organization) (*models.OauthToken, error) {
	users, err := org.User().All(ctx, PGInstance.BobDB)
	if err != nil {
		return nil, fmt.Errorf("Failed to query all users for org: %v", err)
	}
	for _, user := range users {
		oauths, err := user.UserOauthTokens().All(ctx, PGInstance.BobDB)
		if err != nil {
			return nil, fmt.Errorf("Failed to query all oauth tokens for org: %v", err)
		}
		for _, oauth := range oauths {
			return oauth, nil
		}
	}
	return nil, errors.New("No oauth tokens found")
}

func periodicallyExportFieldseeker(ctx context.Context, org *models.Organization) error {
	pollTicker := time.NewTicker(1)
	for {
		select {
		case <-ctx.Done():
			return nil
		case <-pollTicker.C:
			oauth, err := getOAuthForOrg(ctx, org)
			if err != nil {
				return fmt.Errorf("Failed to get oauth for org: %v", err)
			}
			err = exportFieldseekerData(ctx, org, oauth)
			if err != nil {
				return fmt.Errorf("Failed to export Fieldseeker data: %v", err)
			}
			slog.Info("Completed exporting data, waiting 15 minutes to go agoin.")
			pollTicker = time.NewTicker(15 * time.Minute)
		}
	}
}
func exportFieldseekerData(ctx context.Context, org *models.Organization, oauth *models.OauthToken) error {
	slog.Info("Update Fieldseeker data")
	ar := arcgis.NewArcGIS(
		arcgis.AuthenticatorOAuth{
			AccessToken:         oauth.AccessToken,
			AccessTokenExpires:  oauth.AccessTokenExpires,
			RefreshToken:        oauth.RefreshToken,
			RefreshTokenExpires: oauth.RefreshTokenExpires,
		},
	)
	row, err := sql.OrgByOauthId(oauth.ID).One(ctx, PGInstance.BobDB)
	if err != nil {
		return fmt.Errorf("Failed to get org ID: %v", err)
	}
	fssync, err := fieldseeker.NewFieldSeeker(
		ar,
		row.FieldseekerURL.MustGet(),
	)
	if err != nil {
		return fmt.Errorf("Failed to create fssync: %v", err)
	}
	var stats SyncStats
	for _, layer := range fssync.FeatureServerLayers() {
		ss, err := downloadAllRecords(ctx, fssync, layer, row.OrganizationID)
		if err != nil {
			return fmt.Errorf("Failed to get layer %s: %v", layer, err)
		}
		stats.Inserts += ss.Inserts
		stats.Updates += ss.Updates
		stats.Unchanged += ss.Unchanged
	}

	setter := models.FieldseekerSyncSetter{
		RecordsCreated:   omit.From(int32(stats.Inserts)),
		RecordsUpdated:   omit.From(int32(stats.Updates)),
		RecordsUnchanged: omit.From(int32(stats.Unchanged)),
	}
	err = org.InsertFieldseekerSyncs(ctx, PGInstance.BobDB, &setter)
	if err != nil {
		return fmt.Errorf("Failed to insert sync: %v", err)
	}

	return nil
}

func maintainOAuth(ctx context.Context, oauth *models.OauthToken) error {
	refreshDelay := time.Until(oauth.AccessTokenExpires)
	slog.Info("Need to refresh oauth", slog.Int("id", int(oauth.ID)), slog.Float64("seconds", refreshDelay.Seconds()))
	if oauth.AccessTokenExpires.Before(time.Now()) {
		err := refreshOAuth(ctx, oauth)
		if err != nil {
			markTokenFailed(ctx, oauth)
			return fmt.Errorf("Failed to refresh token: %v", err)
		}
		refreshDelay = time.Until(oauth.AccessTokenExpires)
	}
	refreshTicker := time.NewTicker(refreshDelay)
	for {
		select {
		case <-ctx.Done():
			return nil
		case <-refreshTicker.C:
			err := refreshOAuth(ctx, oauth)
			if err != nil {
				return fmt.Errorf("Failed to refresh token: %v", err)
			}
		}
	}

}

// Mark that a given oauth token has failed. This includes a notification to
// the user.
func markTokenFailed(ctx context.Context, oauth *models.OauthToken) {
	oauthSetter := models.OauthTokenSetter{
		InvalidatedAt: omitnull.From(time.Now()),
	}
	err := oauth.Update(ctx, PGInstance.BobDB, &oauthSetter)
	if err != nil {
		slog.Error("Failed to mark token failed", slog.String("err", err.Error()))
	}
	user, err := models.FindUser(ctx, PGInstance.BobDB, oauth.UserID)
	if err != nil {
		slog.Error("Failed to get oauth user", slog.String("err", err.Error()))
		return
	}
	notificationSetter := models.NotificationSetter{
		Message: omitnull.From("Oauth token invalidated"),
		Link:    omitnull.From("/oauth/refresh"),
		Type:    omitnull.From(enums.NotificationtypeOauthTokenInvalidated),
	}
	err = user.InsertUserNotifications(ctx, PGInstance.BobDB, &notificationSetter)
	if err != nil {
		slog.Error("Failed to get oauth user", slog.String("err", err.Error()))
		return
	}
	slog.Info("Marked oauth token invalid", slog.Int("id", int(oauth.ID)))
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
	setter := models.OauthTokenSetter{
		AccessToken:        omit.From(token.AccessToken),
		AccessTokenExpires: omit.From(accessExpires),
		Username:           omit.From(token.Username),
	}
	err = oauth.Update(ctx, PGInstance.BobDB, &setter)
	if err != nil {
		return fmt.Errorf("Failed to update oauth in database: %v", err)
	}
	slog.Info("Updated oauth token", slog.Int("oauth token id", int(oauth.ID)))
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

func saveRawQuery(fssync *fieldseeker.FieldSeeker, layer arcgis.LayerFeature, query *arcgis.Query, filename string) {
	output, err := os.Create(filename)
	if err != nil {
		slog.Error("Failed to create file", slog.String("filename", filename))
		return
	}
	qr, err := fssync.DoQueryRaw(
		layer.ID,
		query)
	if err != nil {
		slog.Error("Failed to do query", slog.String("err", err.Error()))
		return
	}
	_, err = output.Write(qr)
	if err != nil {
		slog.Error("Failed to write results", slog.String("err", err.Error()))
		return
	}
	slog.Info("Wrote failed query", slog.String("filename", filename))
}

func saveOrUpdateDBRecords(ctx context.Context, table string, qr *arcgis.QueryResult, org_id int32) (int, int, error) {
	inserts, updates := 0, 0
	sorted_columns := make([]string, 0, len(qr.Fields))
	for _, f := range qr.Fields {
		sorted_columns = append(sorted_columns, f.Name)
	}
	sort.Strings(sorted_columns)

	objectids := make([]int, 0)
	for _, l := range qr.Features {
		oid := l.Attributes["OBJECTID"].(float64)
		objectids = append(objectids, int(oid))
	}

	rows_by_objectid, err := rowmapViaQuery(ctx, table, sorted_columns, objectids)
	if err != nil {
		return inserts, updates, fmt.Errorf("Failed to get existing rows: %v", err)
	}
	// log.Println("Rows from query", len(rows_by_objectid))

	for _, feature := range qr.Features {
		oid := feature.Attributes["OBJECTID"].(float64)
		row := rows_by_objectid[int(oid)]
		// If we have no matching row we'll need to create it
		if len(row) == 0 {

			if err := insertRowFromFeature(ctx, table, sorted_columns, &feature, org_id); err != nil {
				return inserts, updates, fmt.Errorf("Failed to insert row: %v", err)
			}
			inserts += 1
		} else if hasUpdates(row, feature) {
			if err := updateRowFromFeature(ctx, table, sorted_columns, &feature, org_id); err != nil {
				return inserts, updates, fmt.Errorf("Failed to update row: %v", err)
			}
			updates += 1
		}
	}
	return inserts, updates, nil
}

// Produces a map of OBJECTID to a 'row' which is in turn a map of column names to their values as strings
func rowmapViaQuery(ctx context.Context, table string, sorted_columns []string, objectids []int) (map[int]map[string]string, error) {
	result := make(map[int]map[string]string)

	query := selectAllFromQueryResult(table, sorted_columns)

	args := pgx.NamedArgs{
		"objectids": objectids,
	}
	rows, err := PGInstance.PGXPool.Query(ctx, query, args)
	if err != nil {
		return result, fmt.Errorf("Failed to query rows: %v", err)
	}
	defer rows.Close()

	// +2 for geometry x and geometry x
	columnNames := make([]string, len(sorted_columns)+2)
	for i, c := range sorted_columns {
		columnNames[i] = c
	}
	columnNames[len(sorted_columns)] = "geometry_x"
	columnNames[len(sorted_columns)+1] = "geometry_y"

	rowSlice, err := pgx.CollectRows(rows, func(row pgx.CollectableRow) (map[string]string, error) {
		fieldDescriptions := row.FieldDescriptions()
		values := make([]interface{}, len(fieldDescriptions))
		valuePtrs := make([]interface{}, len(fieldDescriptions))

		for i := range values {
			valuePtrs[i] = &values[i]
		}

		if err := row.Scan(valuePtrs...); err != nil {
			return nil, err
		}

		result := make(map[string]string)
		for i, fd := range fieldDescriptions {
			if values[i] != nil {
				result[fd.Name] = fmt.Sprintf("%v", values[i])
				//log.Printf("col %v type %T val %v", fd.Name, values[i], values[i])
			} else {
				result[fd.Name] = ""
			}
		}

		return result, nil
	})
	if err != nil {
		return result, fmt.Errorf("Failed to collect rows: %v", err)
	}
	for _, row := range rowSlice {
		o := row["objectid"]
		objectid, err := strconv.Atoi(o)
		if err != nil {
			return result, fmt.Errorf("Failed to parse objectid %s: %v", o, err)
		}
		result[objectid] = row
	}
	return result, nil
}

func insertRowFromFeature(ctx context.Context, table string, sorted_columns []string, feature *arcgis.Feature, org_id int32) error {
	var options pgx.TxOptions
	transaction, err := PGInstance.PGXPool.BeginTx(ctx, options)
	if err != nil {
		return fmt.Errorf("Unable to start transaction")
	}

	err = insertRowFromFeatureFS(ctx, transaction, table, sorted_columns, feature, org_id)
	if err != nil {
		transaction.Rollback(ctx)
		return fmt.Errorf("Unable to insert FS: %v", err)
	}

	err = insertRowFromFeatureHistory(ctx, transaction, table, sorted_columns, feature, org_id, 1)
	if err != nil {
		transaction.Rollback(ctx)
		return fmt.Errorf("Failed to insert history: %v", err)
	}

	err = transaction.Commit(ctx)
	if err != nil {
		return fmt.Errorf("Failed to commit transaction: %v", err)
	}
	return nil
}

func insertRowFromFeatureFS(ctx context.Context, transaction pgx.Tx, table string, sorted_columns []string, feature *arcgis.Feature, org_id int32) error {
	// Create the query to produce the main row
	var sb strings.Builder
	sb.WriteString("INSERT INTO ")
	sb.WriteString(table)
	sb.WriteString(" (")
	for _, field := range sorted_columns {
		sb.WriteString(field)
		sb.WriteString(",")
	}
	// Specially add the geometry values since they aren't in the fields
	sb.WriteString("geometry_x,geometry_y,organization_id,updated")
	sb.WriteString(")\nVALUES (")
	for _, field := range sorted_columns {
		sb.WriteString("@")
		sb.WriteString(field)
		sb.WriteString(",")
	}
	// Specially add the geometry values since they aren't in the fields
	sb.WriteString("@geometry_x,@geometry_y,@organization_id,@updated)")

	args := pgx.NamedArgs{}
	for k, v := range feature.Attributes {
		args[k] = v
	}
	// specially add geometry since it isn't in the list of attributes
	args["geometry_x"] = feature.Geometry.X
	args["geometry_y"] = feature.Geometry.Y
	args["organization_id"] = org_id
	args["updated"] = time.Now()

	_, err := transaction.Exec(ctx, sb.String(), args)
	if err != nil {
		return fmt.Errorf("Failed to insert row into %s: %v", table, err)
	}
	return nil
}
func hasUpdates(row map[string]string, feature arcgis.Feature) bool {
	for key, value := range feature.Attributes {
		rowdata := row[strings.ToLower(key)]
		// We'll accept any 'nil' as represented by the empty string in the database
		if value == nil {
			if rowdata == "" {
				continue
			} else if len(rowdata) > 0 {
				return true
			} else {
				slog.Error("Looks like our original value is nil, but our row value is something non-empty with a zero length. Need a programmer to look into this.")
			}
		}
		// check strings first, their simplest
		if featureAsString, ok := value.(string); ok {
			if featureAsString != rowdata {
				return true
			}
			continue
		} else if featureAsInt, ok := value.(int); ok {
			// Previously had a nil value, now we have a real value
			if rowdata == "" {
				return true
			}
			rowAsInt, err := strconv.Atoi(rowdata)
			if err != nil {
				slog.Error(fmt.Sprintf("Failed to convert '%s' to an int to compare against %v for %v", rowdata, featureAsInt, key))
			}
			if rowAsInt != featureAsInt {
				return true
			} else {
				continue
			}
		} else if featureAsFloat, ok := value.(float64); ok {
			// Previously had a nil value, now we have a real value
			if rowdata == "" {
				return true
			}
			rowAsFloat, err := strconv.ParseFloat(rowdata, 64)
			if err != nil {
				slog.Error(fmt.Sprintf("Failed to convert '%s' to a float64 to compare against %v for %v", rowdata, featureAsFloat, key))
			}
			if rowAsFloat != featureAsFloat {
				return true
			} else {
				continue
			}
		}
		slog.Info(fmt.Sprintf("key: %s\tvalue: %v (type %T)\trow: %s\n", key, value, value, rowdata))
		slog.Error("we've hit a point where we can't tell if we have an update or not, need a programmer to look at the above")
	}
	return false
}
func updateRowFromFeature(ctx context.Context, table string, sorted_columns []string, feature *arcgis.Feature, org_id int32) error {
	// Get the current highest version for the row in question
	history_table := toHistoryTable(table)
	var sb strings.Builder
	sb.WriteString("SELECT MAX(version) FROM ")
	sb.WriteString(history_table)
	sb.WriteString(" WHERE OBJECTID=@objectid")

	args := pgx.NamedArgs{}
	o := feature.Attributes["OBJECTID"].(float64)
	args["objectid"] = int(o)

	var version int
	if err := PGInstance.PGXPool.QueryRow(ctx, sb.String(), args).Scan(&version); err != nil {
		return fmt.Errorf("Failed to query for version: %v", err)
	}

	var options pgx.TxOptions
	transaction, err := PGInstance.PGXPool.BeginTx(ctx, options)
	if err != nil {
		return fmt.Errorf("Unable to start transaction")
	}

	err = insertRowFromFeatureHistory(ctx, transaction, table, sorted_columns, feature, org_id, version+1)
	if err != nil {
		transaction.Rollback(ctx)
		return fmt.Errorf("Failed to insert history: %v", err)
	}
	err = updateRowFromFeatureFS(ctx, transaction, table, sorted_columns, feature)
	if err != nil {
		transaction.Rollback(ctx)
		return fmt.Errorf("Failed to update row from feature: %v", err)
	}

	err = transaction.Commit(ctx)
	if err != nil {
		return fmt.Errorf("Failed to commit transaction: %v", err)
	}
	return nil
}
func insertRowFromFeatureHistory(ctx context.Context, transaction pgx.Tx, table string, sorted_columns []string, feature *arcgis.Feature, org_id int32, version int) error {
	history_table := toHistoryTable(table)
	var sb strings.Builder
	sb.WriteString("INSERT INTO ")
	sb.WriteString(history_table)
	sb.WriteString(" (")
	for _, field := range sorted_columns {
		sb.WriteString(field)
		sb.WriteString(",")
	}
	// Specially add the geometry values since they aren't in the fields
	sb.WriteString("created,geometry_x,geometry_y,organization_id,version")
	sb.WriteString(")\nVALUES (")
	for _, field := range sorted_columns {
		sb.WriteString("@")
		sb.WriteString(field)
		sb.WriteString(",")
	}
	// Specially add the geometry values since they aren't in the fields
	sb.WriteString("@created,@geometry_x,@geometry_y,@organization_id,@version)")
	args := pgx.NamedArgs{}
	for k, v := range feature.Attributes {
		args[k] = v
	}
	args["created"] = time.Now()
	args["organization_id"] = org_id
	args["version"] = version
	if _, err := transaction.Exec(ctx, sb.String(), args); err != nil {
		return fmt.Errorf("Failed to insert history row into %s: %v", table, err)
	}
	return nil
}
func selectAllFromQueryResult(table string, sorted_columns []string) string {
	var sb strings.Builder
	sb.WriteString("SELECT * FROM ")
	sb.WriteString(table)
	sb.WriteString(" WHERE OBJECTID=ANY(@objectids)")
	return sb.String()
}
func toHistoryTable(table string) string {
	return "History_" + table[3:len(table)]
}

func updateRowFromFeatureFS(ctx context.Context, transaction pgx.Tx, table string, sorted_columns []string, feature *arcgis.Feature) error {
	// Create the query to produce the main row
	var sb strings.Builder
	sb.WriteString("UPDATE ")
	sb.WriteString(table)
	sb.WriteString(" SET ")
	for _, field := range sorted_columns {
		// OBJECTID is special as our primary key, so skip it
		if field == "OBJECTID" {
			continue
		}
		sb.WriteString(field)
		sb.WriteString("=@")
		sb.WriteString(field)
		sb.WriteString(",")
	}
	// Specially add the geometry values since they aren't in the fields
	sb.WriteString("geometry_x=@geometry_x,geometry_y=@geometry_y,updated=@updated WHERE OBJECTID=@OBJECTID")

	args := pgx.NamedArgs{}
	for k, v := range feature.Attributes {
		args[k] = v
	}
	// specially add geometry since it isn't in the list of attributes
	args["geometry_x"] = feature.Geometry.X
	args["geometry_y"] = feature.Geometry.Y
	args["updated"] = time.Now()

	_, err := transaction.Exec(ctx, sb.String(), args)
	if err != nil {
		return fmt.Errorf("Failed to update row into %s: %v", table, err)
	}
	return nil
}
