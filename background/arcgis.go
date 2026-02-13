package background

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
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/Gleipnir-Technology/arcgis-go"
	"github.com/Gleipnir-Technology/arcgis-go/fieldseeker"
	"github.com/Gleipnir-Technology/arcgis-go/response"
	"github.com/Gleipnir-Technology/bob/dialect/psql"
	"github.com/Gleipnir-Technology/bob/dialect/psql/dm"
	"github.com/Gleipnir-Technology/bob/dialect/psql/sm"
	"github.com/Gleipnir-Technology/bob/dialect/psql/um"
	"github.com/Gleipnir-Technology/nidus-sync/config"
	"github.com/Gleipnir-Technology/nidus-sync/db"
	"github.com/Gleipnir-Technology/nidus-sync/db/models"
	"github.com/Gleipnir-Technology/nidus-sync/db/sql"
	"github.com/Gleipnir-Technology/nidus-sync/debug"
	"github.com/Gleipnir-Technology/nidus-sync/notification"
	"github.com/aarondl/opt/omit"
	"github.com/aarondl/opt/omitnull"
	"github.com/alitto/pond/v2"
	"github.com/jackc/pgx/v5"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

var syncStatusByOrg map[int32]bool

// When the API responds that the token is now invalidated
type InvalidatedTokenError struct{}

func (e InvalidatedTokenError) Error() string { return "The token has been invalidated by the server" }

// When there is no oauth for an organization
type NoOAuthForOrg struct{}

func (e NoOAuthForOrg) Error() string { return "No oauth available for organization" }

var newOAuthTokenChannel chan struct{}
var CodeVerifier string = "random_secure_string_min_43_chars_long_should_be_stored_in_session"

type OAuthTokenResponse struct {
	AccessToken           string `json:"access_token"`
	ExpiresIn             int    `json:"expires_in"`
	RefreshToken          string `json:"refresh_token"`
	RefreshTokenExpiresIn int    `json:"refresh_token_expires_in"`
	SSL                   bool   `json:"ssl"`
	Username              string `json:"username"`
}

func GetOAuthForOrg(ctx context.Context, org *models.Organization) (*models.OauthToken, error) {
	users, err := org.User().All(ctx, db.PGInstance.BobDB)
	if err != nil {
		return nil, fmt.Errorf("Failed to query all users for org: %w", err)
	}
	for _, user := range users {
		oauths, err := user.UserOauthTokens(models.SelectWhere.OauthTokens.InvalidatedAt.IsNull()).All(ctx, db.PGInstance.BobDB)
		if err != nil {
			return nil, fmt.Errorf("Failed to query all oauth tokens for org: %w", err)
		}
		for _, oauth := range oauths {
			return oauth, nil
		}
	}
	return nil, &NoOAuthForOrg{}
}

func HandleOauthAccessCode(ctx context.Context, user *models.User, code string) error {
	baseURL := "https://www.arcgis.com/sharing/rest/oauth2/token/"

	//params.Add("code_verifier", "S256")

	form := url.Values{
		"grant_type":   []string{"authorization_code"},
		"code":         []string{code},
		"client_id":    []string{config.ClientID},
		"redirect_uri": []string{config.ArcGISOauthRedirectURL()},
	}

	req, err := http.NewRequest("POST", baseURL, strings.NewReader(form.Encode()))
	if err != nil {
		return fmt.Errorf("Failed to create request: %w", err)
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	token, err := handleTokenRequest(ctx, req)
	if err != nil {
		return fmt.Errorf("Failed to exchange authorization code for token: %w", err)
	}
	accessExpires := futureUTCTimestamp(token.ExpiresIn)
	refreshExpires := futureUTCTimestamp(token.RefreshTokenExpiresIn)
	setter := models.OauthTokenSetter{
		AccessToken:         omit.From(token.AccessToken),
		AccessTokenExpires:  omit.From(accessExpires),
		ArcgisID:            omitnull.FromPtr[string](nil),
		ArcgisLicenseTypeID: omitnull.FromPtr[string](nil),
		InvalidatedAt:       omitnull.FromPtr[time.Time](nil),
		RefreshToken:        omit.From(token.RefreshToken),
		RefreshTokenExpires: omit.From(refreshExpires),
		UserID:              omit.From(user.ID),
		Username:            omit.From(token.Username),
	}
	oauth, err := models.OauthTokens.Insert(&setter).One(ctx, db.PGInstance.BobDB)
	if err != nil {
		return fmt.Errorf("Failed to save token to database: %w", err)
	}
	go updateArcgisUserData(context.Background(), user, oauth)
	return nil
}

func HasFieldseekerConnection(ctx context.Context, user *models.User) (bool, error) {
	result, err := models.OauthTokens.Query(
		sm.Where(
			models.OauthTokens.Columns.UserID.EQ(psql.Arg(user.ID)),
		),
	).Exists(ctx, db.PGInstance.BobDB)
	if err != nil {
		return false, err
	}
	return result, nil
}

func IsSyncOngoing(org_id int32) bool {
	return syncStatusByOrg[org_id]
}

// This is a goroutine that is in charge of getting Fieldseeker data and keeping it fresh.
func refreshFieldseekerData(background_ctx context.Context, newOauthCh <-chan struct{}) {
	ctx := log.With().Str("component", "arcgis").Logger().Level(zerolog.InfoLevel).WithContext(background_ctx)
	syncStatusByOrg = make(map[int32]bool, 0)
	for {
		workerCtx, cancel := context.WithCancel(context.Background())
		var wg sync.WaitGroup

		oauths, err := models.OauthTokens.Query(models.SelectWhere.OauthTokens.InvalidatedAt.IsNull()).All(ctx, db.PGInstance.BobDB)
		if err != nil {
			log.Error().Err(err).Msg("Failed to get oauths")
			return
		}
		for _, oauth := range oauths {
			wg.Add(1)
			go func() {
				defer wg.Done()
				err := maintainOAuth(workerCtx, oauth)
				if err != nil {
					markTokenFailed(ctx, oauth)
					if errors.Is(err, arcgis.ErrorInvalidRefreshToken) {
						log.Info().Int("oauth_token.id", int(oauth.ID)).Msg("Marked invalid by the server")
					} else {
						debug.LogErrorTypeInfo(err)
						log.Error().Err(err).Msg("Crashed oauth maintenance goroutine")
					}
				}
			}()
		}

		orgs, err := models.Organizations.Query().All(ctx, db.PGInstance.BobDB)
		if err != nil {
			log.Error().Err(err).Msg("Failed to get orgs")
			return
		}
		for _, org := range orgs {
			wg.Add(1)
			go func() {
				defer wg.Done()
				err := periodicallyExportFieldseeker(workerCtx, org)
				if err != nil {
					if errors.Is(err, &NoOAuthForOrg{}) {
						//log.Debug().Int("organization_id", int(org.ID)).Msg("No oauth available for organization, exiting exporter.")
						return
					}
					log.Error().Err(err).Msg("Crashed fieldseeker export goroutine")
				}
			}()
		}

		select {
		case <-ctx.Done():
			log.Info().Msg("Exiting refresh worker...")
			cancel()
			wg.Wait()
			return
		case <-newOauthCh:
			log.Info().Msg("Updating oauth background work")
			cancel()
			wg.Wait()
		}
	}
}

type SyncStats struct {
	Inserts   uint
	Updates   uint
	Unchanged uint
}

func downloadFieldseekerSchema(ctx context.Context, fieldseekerClient *fieldseeker.FieldSeeker, arcgis_id string) {
	layers := fieldseekerClient.Layers()
	log.Debug().Int("len", len(layers)).Msg("Downloading fieldseeker schema")
	for _, layer := range layers {
		err := os.MkdirAll(filepath.Join(config.FieldseekerSchemaDirectory, arcgis_id), os.ModePerm)
		if err != nil {
			log.Error().Err(err).Msg("Failed to create parent directory")
			return
		}
		output, err := os.Create(fmt.Sprintf("%s/%s/%s.json", config.FieldseekerSchemaDirectory, arcgis_id, layer.Name))
		if err != nil {
			log.Error().Err(err).Msg("Failed to open output")
			return
		}
		defer output.Close()
		schema, err := fieldseekerClient.SchemaRaw(ctx, layer.ID)
		if err != nil {
			log.Error().Err(err).Msg("Failed to get schema")
			return
		}
		_, err = output.Write(schema)
		if err != nil {
			log.Error().Err(err).Msg("Failed to write schema file")
			continue
		}
	}
}

func extractURLParts(urlString string) (string, []string, error) {
	parsedURL, err := url.Parse(urlString)
	if err != nil {
		return "", nil, err
	}

	host := parsedURL.Scheme + "://" + parsedURL.Host

	// Split the path and filter empty parts
	var pathParts []string
	for _, part := range strings.Split(parsedURL.Path, "/") {
		if part != "" {
			pathParts = append(pathParts, part)
		}
	}

	return host, pathParts, nil
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
func updateArcgisUserData(ctx context.Context, user *models.User, oauth *models.OauthToken) {
	client, err := arcgis.NewArcGISAuth(
		ctx,
		&arcgis.AuthenticatorOAuth{
			AccessToken:         oauth.AccessToken,
			AccessTokenExpires:  oauth.AccessTokenExpires,
			RefreshToken:        oauth.RefreshToken,
			RefreshTokenExpires: oauth.RefreshTokenExpires,
		},
	)
	if err != nil {
		log.Error().Err(err).Msg("Failed to create ArcGIS client")
		return
	}

	portal, err := updatePortalData(ctx, client, user.ID)
	if err != nil {
		log.Error().Err(err).Msg("Failed to get portal data")
		return
	}

	_, err = models.OauthTokens.Update(
		//um.SetCol(string(models.OauthTokens.Columns.ArcgisID)).ToArg(portal.User.ID),
		//um.SetCol(string(models.OauthTokens.Columns.ArcgisLicenseTypeID)).ToArg(portal.User.UserLicenseTypeID),
		um.SetCol("arcgis_id").ToArg(portal.User.ID),
		um.SetCol("arcgis_license_type_id").ToArg(portal.User.UserLicenseTypeID),
		um.Where(models.OauthTokens.Columns.RefreshToken.EQ(psql.Arg(oauth.RefreshToken))),
	).Exec(ctx, db.PGInstance.BobDB)
	//_, err = sql.UpdateOauthTokenOrg(portal.User.ID, portal.User.UserLicenseTypeID, refresh_token).Exec(ctx, db.PGInstance.BobDB)
	if err != nil {
		log.Error().Err(err).Msg("Failed to update oauth token portal data")
		return
	}
	// At this point we have the arcgis ID. If the ID matches an existing ID, join it with the users organization
	orgs, err := models.Organizations.Query(
		sm.Where(
			psql.Quote("arcgis_id").EQ(psql.Arg(portal.User.OrgID)),
		),
	).All(ctx, db.PGInstance.BobDB)
	if err != nil {
		log.Error().Err(err).Str("arcgis_id", portal.User.OrgID).Msg("Failed to search for organization with arcgis id")
		return
	}
	var org *models.Organization
	switch len(orgs) {
	case 0:
		org = user.R.Organization
		err = org.Update(ctx, db.PGInstance.BobDB, &models.OrganizationSetter{
			ArcgisID:   omitnull.From(portal.User.OrgID),
			ArcgisName: omitnull.From(portal.Name),
		})
		if err != nil {
			log.Error().Err(err).Int32("id", user.R.Organization.ID).Msg("Failed to update organization's arcgis info")
			return
		}
		log.Info().Int32("org_id", org.ID).Str("arcgis_id", portal.User.OrgID).Msg("Updated org arcgis ID")
	case 1:
		org = orgs[0]
		user.Update(ctx, db.PGInstance.BobDB, &models.UserSetter{
			OrganizationID: omit.From(org.ID),
		})
		log.Info().Int32("org_id", org.ID).Int32("user_id", user.ID).Msg("Moved user into organization")
	default:
		log.Warn().Int("orgs", len(orgs)).Msg("Got too many orgs, programmer error.")
		return
	}

	fssync, err := fieldseeker.NewFieldSeekerFromAG(ctx, *client)
	if err != nil {
		log.Error().Err(err).Msg("Failed to create fieldseeker")
	}
	log.Info().Str("url", fssync.ServiceFeature.URL.String()).Msg("Found Fieldseeker")
	setter := models.OrganizationSetter{
		FieldseekerURL: omitnull.From(fssync.ServiceFeature.URL.String()),
	}
	err = org.Update(ctx, db.PGInstance.BobDB, &setter)
	if err != nil {
		log.Error().Err(err).Msg("Failed to create new organization")
		return
	}
	arcgis_id, ok := org.ArcgisID.Get()
	if !ok {
		log.Error().Int("org.id", int(org.ID)).Msg("Cannot get webhooks - ArcGIS ID is null")
	}
	maybeCreateWebhook(ctx, fssync)
	downloadFieldseekerSchema(ctx, fssync, arcgis_id)
	notification.ClearOauth(ctx, user)
	newOAuthTokenChannel <- struct{}{}
}

func NewFieldSeeker(ctx context.Context, oauth *models.OauthToken) (*fieldseeker.FieldSeeker, error) {
	row, err := sql.OrgByOauthId(oauth.ID).One(ctx, db.PGInstance.BobDB)
	if err != nil {
		return nil, fmt.Errorf("Failed to get org ID: %w", err)
	}
	// The URL for fieldseeker should be something like
	// https://foo.arcgis.com/123abc/arcgis/rest/services/FieldSeekerGIS/FeatureServer
	// We need to break it up
	host, pathParts, err := extractURLParts(row.FieldseekerURL.MustGet())
	if err != nil {
		return nil, fmt.Errorf("Failed to break up provided url: %v", err)
	}
	if len(pathParts) < 1 {
		return nil, errors.New("Didn't get enough path parts")
	}
	context := pathParts[0]
	ar, err := arcgis.NewArcGISAuth(
		ctx,
		arcgis.AuthenticatorOAuth{
			AccessToken:         oauth.AccessToken,
			AccessTokenExpires:  oauth.AccessTokenExpires,
			RefreshToken:        oauth.RefreshToken,
			RefreshTokenExpires: oauth.RefreshTokenExpires,
		},
	)
	if err != nil {
		return nil, fmt.Errorf("Failed to create ArcGIS client: %w", err)
	}
	log.Info().Str("context", context).Str("host", host).Msg("Using base fieldseeker URL")
	fssync, err := fieldseeker.NewFieldSeekerFromURL(ctx, *ar, row.FieldseekerURL.MustGet())
	if err != nil {
		if errors.Is(err, arcgis.ErrorInvalidAuthToken) {
			return nil, InvalidatedTokenError{}
		} else if errors.Is(err, arcgis.ErrorInvalidRefreshToken) {
			return nil, InvalidatedTokenError{}
		}
		return nil, fmt.Errorf("Failed to create Fieldseeker client: %w", err)
	}
	return fssync, nil
}
func updatePortalData(ctx context.Context, client *arcgis.ArcGIS, user_id int32) (*response.Portal, error) {
	p, err := client.PortalsSelf(ctx)
	if err != nil {
		return nil, fmt.Errorf("Failed to get ArcGIS user data: %w", err)
	}

	tx, err := db.PGInstance.BobDB.BeginTx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("Failed to create transaction: %w", err)
	}

	_, err = models.ArcgisUserPrivileges.Delete(
		dm.Where(
			models.ArcgisUserPrivileges.Columns.UserID.EQ(psql.Arg(p.User.ID)),
		),
	).Exec(ctx, tx)

	if err != nil {
		tx.Rollback(ctx)
		return nil, fmt.Errorf("Failed to delete previous user privilege data: %w", err)
	}

	_, err = models.ArcgisUsers.Delete(
		dm.Where(
			models.ArcgisUsers.Columns.ID.EQ(psql.Arg(p.User.ID)),
		),
	).Exec(ctx, tx)

	if err != nil {
		tx.Rollback(ctx)
		return nil, fmt.Errorf("Failed to delete previous user data: %w", err)
	}

	setter := models.ArcgisUserSetter{
		Access:            omit.From(p.Access),
		Created:           omit.From(time.Unix(p.User.Created, 0)),
		Email:             omit.From(p.User.Email),
		FullName:          omit.From(p.User.FullName),
		ID:                omit.From(p.User.ID),
		Level:             omit.From(p.User.Level),
		OrgID:             omit.From(p.User.OrgID),
		PublicUserID:      omit.From(user_id),
		Region:            omit.From(p.Region),
		Role:              omit.From(p.User.Role),
		RoleID:            omit.From(p.User.RoleId),
		Username:          omit.From(p.User.Username),
		UserLicenseTypeID: omit.From(p.User.UserLicenseTypeID),
		UserType:          omit.From(p.User.UserType),
	}
	_, err = models.ArcgisUsers.Insert(&setter).One(ctx, tx)
	if err != nil {
		tx.Rollback(ctx)
		return nil, fmt.Errorf("Failed to add arcgis user data: %w", err)
	}
	for _, priv := range p.User.Privileges {
		s := models.ArcgisUserPrivilegeSetter{
			Privilege: omit.From(priv),
			UserID:    omit.From(p.User.ID),
		}
		_, err := models.ArcgisUserPrivileges.Insert(&s).One(ctx, tx)
		if err != nil {
			tx.Rollback(ctx)
			return nil, fmt.Errorf("Failed to add arcgis user privilege data: %w", err)
		}
	}
	err = tx.Commit(ctx)
	log.Info().Str("username", p.User.Username).Str("user_id", p.User.ID).Str("org_id", p.User.OrgID).Str("org_name", p.Name).Str("license_type_id", p.User.UserLicenseTypeID).Msg("Updated portals data")
	return p, nil
}

func maybeCreateWebhook(ctx context.Context, client *fieldseeker.FieldSeeker) {
	webhooks, err := client.WebhookList(ctx)
	if err != nil {
		if errors.Is(err, arcgis.ErrorNotPermitted) {
			log.Info().Msg("This oauth token is not allowed to get webhooks")
			return
		}
		log.Error().Err(err).Msg("Failed to get webhooks")
		return
	}
	if webhooks == nil {
		log.Error().Msg("nil webhooks")
		return
	}
	for _, hook := range *webhooks {
		if hook.Name == "Nidus Sync" {
			log.Info().Msg("Found nidus sync hook")
		} else {
			log.Info().Str("name", hook.Name).Msg("Found webhook")
		}
	}
}

func periodicallyExportFieldseeker(ctx context.Context, org *models.Organization) error {
	pollTicker := time.NewTicker(1)
	for {
		select {
		case <-ctx.Done():
			return nil
		case <-pollTicker.C:
			pollTicker = time.NewTicker(15 * time.Minute)
			oauth, err := GetOAuthForOrg(ctx, org)
			if err != nil {
				return fmt.Errorf("Failed to get oauth for org: %w", err)
			}
			fssync, err := NewFieldSeeker(
				ctx,
				oauth,
			)
			if err != nil {
				return fmt.Errorf("Failed to create fieldseeker client: %w", err)
			}
			logPermissions(ctx, fssync)
			syncStatusByOrg[org.ID] = true
			err = exportFieldseekerData(ctx, fssync, org)
			syncStatusByOrg[org.ID] = false
			if err != nil {
				return fmt.Errorf("Failed to export Fieldseeker data: %w", err)
			}
			log.Info().Msg("Completed exporting data, waiting 15 minutes to go agoin.")
		}
	}
}
func exportFieldseekerData(ctx context.Context, fssync *fieldseeker.FieldSeeker, org *models.Organization) error {
	log.Info().Msg("Update Fieldseeker data")
	var err error
	var stats SyncStats

	pool := pond.NewResultPool[SyncStats](20)
	group := pool.NewGroup()
	var ss SyncStats
	layers := fssync.Layers()
	for _, l := range layers {
		ss, err = exportFieldseekerLayer(ctx, group, org, fssync, l)
		if err != nil {
			return err
		}
		stats.Inserts += ss.Inserts
		stats.Updates += ss.Updates
		stats.Unchanged += ss.Unchanged
	}
	results, err := group.Wait()
	if err != nil {
		return fmt.Errorf("one or more tasks in the work pool failed: %w", err)
	}
	for _, r := range results {
		stats.Inserts += r.Inserts
		stats.Updates += r.Updates
		stats.Unchanged += r.Unchanged
	}

	setter := models.FieldseekerSyncSetter{
		RecordsCreated:   omit.From(int32(stats.Inserts)),
		RecordsUpdated:   omit.From(int32(stats.Updates)),
		RecordsUnchanged: omit.From(int32(stats.Unchanged)),
	}
	err = org.InsertFieldseekerSyncs(ctx, db.PGInstance.BobDB, &setter)
	if err != nil {
		return fmt.Errorf("Failed to insert sync: %w", err)
	}

	updateSummaryTables(ctx, org)
	return nil
}

func logPermissions(ctx context.Context, fssync *fieldseeker.FieldSeeker) {
	/*row, err := sql.OrgByOauthId(oauth.ID).One(ctx, db.PGInstance.BobDB)
	if err != nil {
		log.Error().Err(err).Msg("Failed to get org in log permissions")
		return
	}
	oauth, err := models.FindOauthToken(ctx, db.PGInstance.BobDB, row.ID)
	if err != nil {
		return fmt.Errorf("Failed to update oauth token from database: %w", err)
	}
	*/

	_, err := fssync.AdminInfo(ctx)
	if err != nil {
		if errors.Is(err, arcgis.ErrorNotPermitted) {
			log.Info().Msg("This oauth token is not allowed to query for admin info")
			return
		}
		log.Error().Err(err).Msg("Failed to get admin info log permissions")
		return
	}
	permissions, err := fssync.PermissionList(ctx)
	if err != nil {
		log.Error().Err(err).Msg("Failed to query permissions in log permissions")
		return
	}
	if permissions == nil {
		log.Error().Msg("nil permissions")
		return
	}
	for _, p := range *permissions {
		log.Info().Str("p", p.Principal).Msg("Permission!")
	}
}

func maintainOAuth(ctx context.Context, oauth *models.OauthToken) error {
	for {
		// Refresh from the database
		oauth, err := models.FindOauthToken(ctx, db.PGInstance.BobDB, oauth.ID)
		if err != nil {
			return fmt.Errorf("Failed to update oauth token from database: %w", err)
		}
		var accessTokenDelay time.Duration
		if oauth.AccessTokenExpires.Before(time.Now()) {
			accessTokenDelay = 1
		} else {
			accessTokenDelay = time.Until(oauth.AccessTokenExpires) - (3 * time.Second)
		}
		var refreshTokenDelay time.Duration
		if oauth.RefreshTokenExpires.Before(time.Now()) {
			refreshTokenDelay = 1
		} else {
			refreshTokenDelay = time.Until(oauth.RefreshTokenExpires) - (3 * time.Second)
		}
		log.Info().Int("id", int(oauth.ID)).Float64("seconds", accessTokenDelay.Seconds()).Msg("Need to refresh access token")
		log.Info().Int("id", int(oauth.ID)).Float64("seconds", refreshTokenDelay.Seconds()).Msg("Need to refresh refresh token")
		accessTokenTicker := time.NewTicker(accessTokenDelay)
		refreshTokenTicker := time.NewTicker(refreshTokenDelay)
		select {
		case <-ctx.Done():
			return nil
		case <-accessTokenTicker.C:
			err := refreshAccessToken(ctx, oauth)
			if err != nil {
				return fmt.Errorf("Failed to refresh access token: %w", err)
			}
		case <-refreshTokenTicker.C:
			err := refreshRefreshToken(ctx, oauth)
			if err != nil {
				return fmt.Errorf("Failed to maintain refresh token: %w", err)
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
	err := oauth.Update(ctx, db.PGInstance.BobDB, &oauthSetter)
	if err != nil {
		log.Error().Str("err", err.Error()).Msg("Failed to mark token failed")
	}
	user, err := models.FindUser(ctx, db.PGInstance.BobDB, oauth.UserID)
	if err != nil {
		log.Error().Str("err", err.Error()).Msg("Failed to get oauth user")
		return
	}
	notification.NotifyOauthInvalid(ctx, user)
	log.Info().Int("id", int(oauth.ID)).Msg("Marked oauth token invalid")
}

// Update the access token to keep it fresh and alive
func refreshAccessToken(ctx context.Context, oauth *models.OauthToken) error {
	baseURL := "https://www.arcgis.com/sharing/rest/oauth2/token/"

	form := url.Values{
		"grant_type":    []string{"refresh_token"},
		"client_id":     []string{config.ClientID},
		"refresh_token": []string{oauth.RefreshToken},
	}

	req, err := http.NewRequest("POST", baseURL, strings.NewReader(form.Encode()))
	if err != nil {
		return fmt.Errorf("Failed to create request: %w", err)
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	token, err := handleTokenRequest(ctx, req)
	if err != nil {
		return fmt.Errorf("Failed to handle request: %w", err)
	}
	accessExpires := futureUTCTimestamp(token.ExpiresIn)
	setter := models.OauthTokenSetter{
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
func refreshRefreshToken(ctx context.Context, oauth *models.OauthToken) error {
	baseURL := "https://www.arcgis.com/sharing/rest/oauth2/token/"

	form := url.Values{
		"grant_type":    []string{"exchange_refresh_token"},
		"client_id":     []string{config.ClientID},
		"redirect_uri":  []string{config.ArcGISOauthRedirectURL()},
		"refresh_token": []string{oauth.RefreshToken},
	}

	req, err := http.NewRequest("POST", baseURL, strings.NewReader(form.Encode()))
	if err != nil {
		return fmt.Errorf("Failed to create request: %w", err)
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	token, err := handleTokenRequest(ctx, req)
	if err != nil {
		return fmt.Errorf("Failed to handle request: %w", err)
	}
	refreshExpires := futureUTCTimestamp(token.ExpiresIn)
	setter := models.OauthTokenSetter{
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

func newTimestampedFilename(prefix, suffix string) string {
	timestamp := time.Now().Format("20060102_150405") // YYYYMMDD_HHMMSS format
	return prefix + timestamp + suffix
}

func handleTokenRequest(ctx context.Context, req *http.Request) (*OAuthTokenResponse, error) {
	client := http.Client{}
	log.Info().Str("url", req.URL.String()).Msg("POST")
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("Failed to do request: %w", err)
	}
	defer resp.Body.Close()
	bodyBytes, err := io.ReadAll(resp.Body)
	log.Info().Int("status", resp.StatusCode).Msg("Token request")
	filename := newTimestampedFilename("token", ".json")
	saveResponse(bodyBytes, filename)
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

func logResponseHeaders(resp *http.Response) {
	if resp == nil {
		log.Info().Msg("Response is nil")
		return
	}

	log.Info().Str("status", resp.Status).Int("statusCode", resp.StatusCode).Msg("HTTP Response headers")

	for name, values := range resp.Header {
		log.Info().Str("name", name).Strs("values", values).Msg("Header")
	}
}

func saveResponse(data []byte, filename string) {
	dest, err := os.Create(filename)
	if err != nil {
		log.Error().Str("filename", filename).Str("err", err.Error()).Msg("Failed to create file")
		return
	}
	_, err = io.Copy(dest, bytes.NewReader(data))
	if err != nil {
		log.Error().Str("filename", filename).Str("err", err.Error()).Msg("Failed to write")
		return
	}
	log.Info().Str("filename", filename).Msg("Wrote response")
}

/*
func saveRawQuery(fssync fieldseeker.FieldSeeker, layer arcgis.LayerFeature, query *arcgis.Query, filename string) {
	output, err := os.Create(filename)
	if err != nil {
		log.Error().Str("filename", filename).Msg("Failed to create file")
		return
	}
	qr, err := fssync.DoQueryRaw(
		layer.ID,
		query)
	if err != nil {
		log.Error().Str("err", err.Error()).Msg("Failed to do query")
		return
	}
	_, err = output.Write(qr)
	if err != nil {
		log.Error().Str("err", err.Error()).Msg("Failed to write results")
		return
	}
	log.Info().Str("filename", filename).Msg("Wrote failed query")
}
*/

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
		return inserts, updates, fmt.Errorf("Failed to get existing rows: %w", err)
	}
	// log.Println("Rows from query", len(rows_by_objectid))

	for _, feature := range qr.Features {
		oid := feature.Attributes["OBJECTID"].(float64)
		row := rows_by_objectid[int(oid)]
		// If we have no matching row we'll need to create it
		if len(row) == 0 {

			if err := insertRowFromFeature(ctx, table, sorted_columns, &feature, org_id); err != nil {
				return inserts, updates, fmt.Errorf("Failed to insert row: %w", err)
			}
			inserts += 1
		} else if hasUpdates(row, feature) {
			if err := updateRowFromFeature(ctx, table, sorted_columns, &feature, org_id); err != nil {
				return inserts, updates, fmt.Errorf("Failed to update row: %w", err)
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
	rows, err := db.PGInstance.PGXPool.Query(ctx, query, args)
	if err != nil {
		return result, fmt.Errorf("Failed to query rows: %w", err)
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
		return result, fmt.Errorf("Failed to collect rows: %w", err)
	}
	for _, row := range rowSlice {
		o := row["objectid"]
		objectid, err := strconv.Atoi(o)
		if err != nil {
			return result, fmt.Errorf("Failed to parse objectid %s: %w", o, err)
		}
		result[objectid] = row
	}
	return result, nil
}

func insertRowFromFeature(ctx context.Context, table string, sorted_columns []string, feature *arcgis.Feature, org_id int32) error {
	var options pgx.TxOptions
	transaction, err := db.PGInstance.PGXPool.BeginTx(ctx, options)
	if err != nil {
		return fmt.Errorf("Unable to start transaction")
	}

	err = insertRowFromFeatureFS(ctx, transaction, table, sorted_columns, feature, org_id)
	if err != nil {
		transaction.Rollback(ctx)
		return fmt.Errorf("Unable to insert FS: %w", err)
	}

	err = insertRowFromFeatureHistory(ctx, transaction, table, sorted_columns, feature, org_id, 1)
	if err != nil {
		transaction.Rollback(ctx)
		return fmt.Errorf("Failed to insert history: %w", err)
	}

	err = transaction.Commit(ctx)
	if err != nil {
		return fmt.Errorf("Failed to commit transaction: %w", err)
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
	//args["geometry_x"] = feature.Geometry.X
	//args["geometry_y"] = feature.Geometry.Y
	args["organization_id"] = org_id
	args["updated"] = time.Now()

	_, err := transaction.Exec(ctx, sb.String(), args)
	if err != nil {
		return fmt.Errorf("Failed to insert row into %s: %w", table, err)
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
				log.Error().Msg("Looks like our original value is nil, but our row value is something non-empty with a zero length. Need a programmer to look into this.")
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
				log.Error().Msg(fmt.Sprintf("Failed to convert '%s' to an int to compare against %v for %v", rowdata, featureAsInt, key))
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
				log.Error().Msg(fmt.Sprintf("Failed to convert '%s' to a float64 to compare against %v for %v", rowdata, featureAsFloat, key))
			}
			if rowAsFloat != featureAsFloat {
				return true
			} else {
				continue
			}
		}
		log.Error().Str("key", key).Str("rowdata", rowdata).Msg("we've hit a point where we can't tell if we have an update or not, need a programmer to look at the above")
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
	if err := db.PGInstance.PGXPool.QueryRow(ctx, sb.String(), args).Scan(&version); err != nil {
		return fmt.Errorf("Failed to query for version: %w", err)
	}

	var options pgx.TxOptions
	transaction, err := db.PGInstance.PGXPool.BeginTx(ctx, options)
	if err != nil {
		return fmt.Errorf("Unable to start transaction")
	}

	err = insertRowFromFeatureHistory(ctx, transaction, table, sorted_columns, feature, org_id, version+1)
	if err != nil {
		transaction.Rollback(ctx)
		return fmt.Errorf("Failed to insert history: %w", err)
	}
	err = updateRowFromFeatureFS(ctx, transaction, table, sorted_columns, feature)
	if err != nil {
		transaction.Rollback(ctx)
		return fmt.Errorf("Failed to update row from feature: %w", err)
	}

	err = transaction.Commit(ctx)
	if err != nil {
		return fmt.Errorf("Failed to commit transaction: %w", err)
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
		return fmt.Errorf("Failed to insert history row into %s: %w", table, err)
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
	//args["geometry_x"] = feature.Geometry.X
	//args["geometry_y"] = feature.Geometry.Y
	args["updated"] = time.Now()

	_, err := transaction.Exec(ctx, sb.String(), args)
	if err != nil {
		return fmt.Errorf("Failed to update row into %s: %w", table, err)
	}
	return nil
}

func exportFieldseekerLayer(ctx context.Context, group pond.ResultTaskGroup[SyncStats], org *models.Organization, fssync *fieldseeker.FieldSeeker, layer response.Layer) (SyncStats, error) {
	var stats SyncStats
	count, err := fssync.QueryCount(ctx, layer.ID)
	if err != nil {
		return stats, fmt.Errorf("Failed to get counts for layer %s (%d): %w", layer.Name, layer.ID, err)
	}
	if count.Count == 0 {
		log.Info().Str("name", layer.Name).Uint("layer_id", layer.ID).Int32("org_id", org.ID).Msg("No records to download")
		return stats, nil
	}
	max_records, err := fssync.MaxRecordCount(ctx)
	if err != nil {
		return stats, fmt.Errorf("Failed to get max records: %w", err)
	}
	l, err := fieldseeker.NameToLayerType(layer.Name)
	if err != nil {
		return stats, fmt.Errorf("Failed to get layer for '%s': %w", layer.Name, err)
	}
	log.Info().Str("name", layer.Name).Uint("layer_id", layer.ID).Int32("org_id", org.ID).Int("count", count.Count).Uint("iterations", uint(count.Count)/uint(max_records)).Msg("Queuing jobs for layer")
	for offset := uint(0); offset < uint(count.Count); offset += uint(max_records) {
		group.SubmitErr(func() (SyncStats, error) {
			var ss SyncStats
			var name string
			var inserts, unchanged, updates uint
			var err error
			switch l {
			case fieldseeker.LayerAerialSpraySession:
				name = "AerialSpraySession"
				rows, err := fssync.AerialSpraySession(ctx, offset)
				if err != nil {
					return SyncStats{}, fmt.Errorf("Failed to query %s: %w", name, err)
				}
				inserts, updates, err = db.SaveOrUpdateAerialSpraySession(ctx, org, rows)
				if err != nil {
					return SyncStats{}, fmt.Errorf("Failed to update %s: %w", name, err)
				}
				unchanged = uint(len(rows)) - inserts - updates
			case fieldseeker.LayerAerialSprayLine:
				name = "LayerAerialSprayLine"
				rows, err := fssync.AerialSprayLine(ctx, offset)
				if err != nil {
					return SyncStats{}, fmt.Errorf("Failed to query %s: %w", name, err)
				}
				inserts, updates, err = db.SaveOrUpdateAerialSprayLine(ctx, org, rows)
				if err != nil {
					return SyncStats{}, fmt.Errorf("Failed to update %s: %w", name, err)
				}
				unchanged = uint(len(rows)) - inserts - updates
			case fieldseeker.LayerBarrierSpray:
				name = "LayerBarrierSpray"
				rows, err := fssync.BarrierSpray(ctx, offset)
				if err != nil {
					return SyncStats{}, fmt.Errorf("Failed to query %s: %w", name, err)
				}
				inserts, updates, err = db.SaveOrUpdateBarrierSpray(ctx, org, rows)
				if err != nil {
					return SyncStats{}, fmt.Errorf("Failed to update %s: %w", name, err)
				}
				unchanged = uint(len(rows)) - inserts - updates
			case fieldseeker.LayerBarrierSprayRoute:
				name = "LayerBarrierSprayRoute"
				rows, err := fssync.BarrierSprayRoute(ctx, offset)
				if err != nil {
					return SyncStats{}, fmt.Errorf("Failed to query %s: %w", name, err)
				}
				inserts, updates, err = db.SaveOrUpdateBarrierSprayRoute(ctx, org, rows)
				if err != nil {
					return SyncStats{}, fmt.Errorf("Failed to update %s: %w", name, err)
				}
				unchanged = uint(len(rows)) - inserts - updates
			case fieldseeker.LayerContainerRelate:
				name = "LayerContainerRelate"
				rows, err := fssync.ContainerRelate(ctx, offset)
				if err != nil {
					return SyncStats{}, fmt.Errorf("Failed to query %s: %w", name, err)
				}
				inserts, updates, err = db.SaveOrUpdateContainerRelate(ctx, org, rows)
				if err != nil {
					return SyncStats{}, fmt.Errorf("Failed to update %s: %w", name, err)
				}
				unchanged = uint(len(rows)) - inserts - updates
			case fieldseeker.LayerFieldScoutingLog:
				name = "LayerFieldScoutingLog"
				rows, err := fssync.FieldScoutingLog(ctx, offset)
				if err != nil {
					return SyncStats{}, fmt.Errorf("Failed to query %s: %w", name, err)
				}
				inserts, updates, err = db.SaveOrUpdateFieldScoutingLog(ctx, org, rows)
				if err != nil {
					return SyncStats{}, fmt.Errorf("Failed to update %s: %w", name, err)
				}
				unchanged = uint(len(rows)) - inserts - updates
			case fieldseeker.LayerHabitatRelate:
				name = "LayerHabitatRelate"
				rows, err := fssync.HabitatRelate(ctx, offset)
				if err != nil {
					return SyncStats{}, fmt.Errorf("Failed to query %s: %w", name, err)
				}
				inserts, updates, err = db.SaveOrUpdateHabitatRelate(ctx, org, rows)
				if err != nil {
					return SyncStats{}, fmt.Errorf("Failed to update %s: %w", name, err)
				}
				unchanged = uint(len(rows)) - inserts - updates
			case fieldseeker.LayerInspectionSample:
				name = "LayerInspectionSample"
				rows, err := fssync.InspectionSample(ctx, offset)
				if err != nil {
					return SyncStats{}, fmt.Errorf("Failed to query %s: %w", name, err)
				}
				inserts, updates, err = db.SaveOrUpdateInspectionSample(ctx, org, rows)
				if err != nil {
					return SyncStats{}, fmt.Errorf("Failed to update %s: %w", name, err)
				}
				unchanged = uint(len(rows)) - inserts - updates
			case fieldseeker.LayerInspectionSampleDetail:
				name = "LayerInspectionSampleDetail"
				rows, err := fssync.InspectionSampleDetail(ctx, offset)
				if err != nil {
					return SyncStats{}, fmt.Errorf("Failed to query %s: %w", name, err)
				}
				inserts, updates, err = db.SaveOrUpdateInspectionSampleDetail(ctx, org, rows)
				if err != nil {
					return SyncStats{}, fmt.Errorf("Failed to update %s: %w", name, err)
				}
				unchanged = uint(len(rows)) - inserts - updates
			case fieldseeker.LayerLandingCount:
				name = "LayerLandingCount"
				rows, err := fssync.LandingCount(ctx, offset)
				if err != nil {
					return SyncStats{}, fmt.Errorf("Failed to query %s: %w", name, err)
				}
				inserts, updates, err = db.SaveOrUpdateLandingCount(ctx, org, rows)
				if err != nil {
					return SyncStats{}, fmt.Errorf("Failed to update %s: %w", name, err)
				}
				unchanged = uint(len(rows)) - inserts - updates
			case fieldseeker.LayerLandingCountLocation:
				name = "LayerLandingCountLocation"
				rows, err := fssync.LandingCountLocation(ctx, offset)
				if err != nil {
					return SyncStats{}, fmt.Errorf("Failed to query %s: %w", name, err)
				}
				inserts, updates, err = db.SaveOrUpdateLandingCountLocation(ctx, org, rows)
				if err != nil {
					return SyncStats{}, fmt.Errorf("Failed to update %s: %w", name, err)
				}
				unchanged = uint(len(rows)) - inserts - updates
			case fieldseeker.LayerLineLocation:
				name = "LayerLineLocation"
				rows, err := fssync.LineLocation(ctx, offset)
				if err != nil {
					return SyncStats{}, fmt.Errorf("Failed to query %s: %w", name, err)
				}
				inserts, updates, err = db.SaveOrUpdateLineLocation(ctx, org, rows)
				if err != nil {
					return SyncStats{}, fmt.Errorf("Failed to update %s: %w", name, err)
				}
				unchanged = uint(len(rows)) - inserts - updates
			case fieldseeker.LayerLocationTracking:
				name = "LayerLocationTracking"
				rows, err := fssync.LocationTracking(ctx, offset)
				if err != nil {
					return SyncStats{}, fmt.Errorf("Failed to query %s: %w", name, err)
				}
				inserts, updates, err = db.SaveOrUpdateLocationTracking(ctx, org, rows)
				if err != nil {
					return SyncStats{}, fmt.Errorf("Failed to update %s: %w", name, err)
				}
				unchanged = uint(len(rows)) - inserts - updates
			case fieldseeker.LayerMosquitoInspection:
				name = "LayerMosquitoInspection"
				rows, err := fssync.MosquitoInspection(ctx, offset)
				if err != nil {
					return SyncStats{}, fmt.Errorf("Failed to query %s: %w", name, err)
				}
				inserts, updates, err = db.SaveOrUpdateMosquitoInspection(ctx, org, rows)
				if err != nil {
					return SyncStats{}, fmt.Errorf("Failed to update %s: %w", name, err)
				}
				unchanged = uint(len(rows)) - inserts - updates
			case fieldseeker.LayerOfflineMapAreas:
				name = "LayerOfflineMapAreas"
				rows, err := fssync.OfflineMapAreas(ctx, offset)
				if err != nil {
					return SyncStats{}, fmt.Errorf("Failed to query %s: %w", name, err)
				}
				inserts, updates, err = db.SaveOrUpdateOfflineMapAreas(ctx, org, rows)
				if err != nil {
					return SyncStats{}, fmt.Errorf("Failed to update %s: %w", name, err)
				}
				unchanged = uint(len(rows)) - inserts - updates
			case fieldseeker.LayerProposedTreatmentArea:
				name = "LayerProposedTreatmentArea"
				rows, err := fssync.ProposedTreatmentArea(ctx, offset)
				if err != nil {
					return SyncStats{}, fmt.Errorf("Failed to query %s: %w", name, err)
				}
				inserts, updates, err = db.SaveOrUpdateProposedTreatmentArea(ctx, org, rows)
				if err != nil {
					return SyncStats{}, fmt.Errorf("Failed to update %s: %w", name, err)
				}
				unchanged = uint(len(rows)) - inserts - updates
			case fieldseeker.LayerPointLocation:
				name = "LayerPointLocation"
				rows, err := fssync.PointLocation(ctx, offset)
				if err != nil {
					return SyncStats{}, fmt.Errorf("Failed to query %s: %w", name, err)
				}
				inserts, updates, err = db.SaveOrUpdatePointLocation(ctx, org, rows)
				if err != nil {
					return SyncStats{}, fmt.Errorf("Failed to update %s: %w", name, err)
				}
				unchanged = uint(len(rows)) - inserts - updates
			case fieldseeker.LayerPolygonLocation:
				name = "LayerPolygonLocation"
				rows, err := fssync.PolygonLocation(ctx, offset)
				if err != nil {
					return SyncStats{}, fmt.Errorf("Failed to query %s: %w", name, err)
				}
				inserts, updates, err = db.SaveOrUpdatePolygonLocation(ctx, org, rows)
				if err != nil {
					return SyncStats{}, fmt.Errorf("Failed to update %s: %w", name, err)
				}
				unchanged = uint(len(rows)) - inserts - updates
			case fieldseeker.LayerPoolDetail:
				name = "LayerPoolDetail"
				rows, err := fssync.PoolDetail(ctx, offset)
				if err != nil {
					return SyncStats{}, fmt.Errorf("Failed to query %s: %w", name, err)
				}
				inserts, updates, err = db.SaveOrUpdatePoolDetail(ctx, org, rows)
				if err != nil {
					return SyncStats{}, fmt.Errorf("Failed to update %s: %w", name, err)
				}
				unchanged = uint(len(rows)) - inserts - updates
			case fieldseeker.LayerPool:
				name = "LayerPool"
				rows, err := fssync.Pool(ctx, offset)
				if err != nil {
					return SyncStats{}, fmt.Errorf("Failed to query %s: %w", name, err)
				}
				inserts, updates, err = db.SaveOrUpdatePool(ctx, org, rows)
				if err != nil {
					return SyncStats{}, fmt.Errorf("Failed to update %s: %w", name, err)
				}
				unchanged = uint(len(rows)) - inserts - updates
			case fieldseeker.LayerPoolBuffer:
				name = "LayerPoolBuffer"
				rows, err := fssync.PoolBuffer(ctx, offset)
				if err != nil {
					return SyncStats{}, fmt.Errorf("Failed to query %s: %w", name, err)
				}
				inserts, updates, err = db.SaveOrUpdatePoolBuffer(ctx, org, rows)
				if err != nil {
					return SyncStats{}, fmt.Errorf("Failed to update %s: %w", name, err)
				}
				unchanged = uint(len(rows)) - inserts - updates
			case fieldseeker.LayerQALarvCount:
				name = "LayerQALarvCount"
				rows, err := fssync.QALarvCount(ctx, offset)
				if err != nil {
					return SyncStats{}, fmt.Errorf("Failed to query %s: %w", name, err)
				}
				inserts, updates, err = db.SaveOrUpdateQALarvCount(ctx, org, rows)
				if err != nil {
					return SyncStats{}, fmt.Errorf("Failed to update %s: %w", name, err)
				}
				unchanged = uint(len(rows)) - inserts - updates
			case fieldseeker.LayerQAMosquitoInspection:
				name = "LayerQAMosquitoInspection"
				rows, err := fssync.QAMosquitoInspection(ctx, offset)
				if err != nil {
					return SyncStats{}, fmt.Errorf("Failed to query %s: %w", name, err)
				}
				inserts, updates, err = db.SaveOrUpdateQAMosquitoInspection(ctx, org, rows)
				if err != nil {
					return SyncStats{}, fmt.Errorf("Failed to update %s: %w", name, err)
				}
				unchanged = uint(len(rows)) - inserts - updates
			case fieldseeker.LayerQAProductObservation:
				name = "LayerQAProductObservation"
				rows, err := fssync.QAProductObservation(ctx, offset)
				if err != nil {
					return SyncStats{}, fmt.Errorf("Failed to query %s: %w", name, err)
				}
				inserts, updates, err = db.SaveOrUpdateQAProductObservation(ctx, org, rows)
				if err != nil {
					return SyncStats{}, fmt.Errorf("Failed to update %s: %w", name, err)
				}
				unchanged = uint(len(rows)) - inserts - updates
			case fieldseeker.LayerRestrictedArea:
				name = "LayerRestrictedArea"
				rows, err := fssync.RestrictedArea(ctx, offset)
				if err != nil {
					return SyncStats{}, fmt.Errorf("Failed to query %s: %w", name, err)
				}
				inserts, updates, err = db.SaveOrUpdateRestrictedArea(ctx, org, rows)
				if err != nil {
					return SyncStats{}, fmt.Errorf("Failed to update %s: %w", name, err)
				}
				unchanged = uint(len(rows)) - inserts - updates
			case fieldseeker.LayerRodentInspection:
				name = "LayerRodentInspection"
				rows, err := fssync.RodentInspection(ctx, offset)
				if err != nil {
					return SyncStats{}, fmt.Errorf("Failed to query %s: %w", name, err)
				}
				inserts, updates, err = db.SaveOrUpdateRodentInspection(ctx, org, rows)
				if err != nil {
					return SyncStats{}, fmt.Errorf("Failed to update %s: %w", name, err)
				}
				unchanged = uint(len(rows)) - inserts - updates
			case fieldseeker.LayerRodentLocation:
				name = "LayerRodentLocation"
				rows, err := fssync.RodentLocation(ctx, offset)
				if err != nil {
					return SyncStats{}, fmt.Errorf("Failed to query %s: %w", name, err)
				}
				inserts, updates, err = db.SaveOrUpdateRodentLocation(ctx, org, rows)
				if err != nil {
					return SyncStats{}, fmt.Errorf("Failed to update %s: %w", name, err)
				}
				unchanged = uint(len(rows)) - inserts - updates
			case fieldseeker.LayerSampleCollection:
				name = "LayerSampleCollection"
				rows, err := fssync.SampleCollection(ctx, offset)
				if err != nil {
					return SyncStats{}, fmt.Errorf("Failed to query %s: %w", name, err)
				}
				inserts, updates, err = db.SaveOrUpdateSampleCollection(ctx, org, rows)
				if err != nil {
					return SyncStats{}, fmt.Errorf("Failed to update %s: %w", name, err)
				}
				unchanged = uint(len(rows)) - inserts - updates
			case fieldseeker.LayerSampleLocation:
				name = "LayerSampleLocation"
				rows, err := fssync.SampleLocation(ctx, offset)
				if err != nil {
					return SyncStats{}, fmt.Errorf("Failed to query %s: %w", name, err)
				}
				inserts, updates, err = db.SaveOrUpdateSampleLocation(ctx, org, rows)
				if err != nil {
					return SyncStats{}, fmt.Errorf("Failed to update %s: %w", name, err)
				}
				unchanged = uint(len(rows)) - inserts - updates
			case fieldseeker.LayerServiceRequest:
				name = "LayerServiceRequest"
				rows, err := fssync.ServiceRequest(ctx, offset)
				if err != nil {
					return SyncStats{}, fmt.Errorf("Failed to query %s: %w", name, err)
				}
				inserts, updates, err = db.SaveOrUpdateServiceRequest(ctx, org, rows)
				if err != nil {
					return SyncStats{}, fmt.Errorf("Failed to update %s: %w", name, err)
				}
				unchanged = uint(len(rows)) - inserts - updates
			case fieldseeker.LayerSpeciesAbundance:
				name = "LayerSpeciesAbundance"
				rows, err := fssync.SpeciesAbundance(ctx, offset)
				if err != nil {
					return SyncStats{}, fmt.Errorf("Failed to query %s: %w", name, err)
				}
				inserts, updates, err = db.SaveOrUpdateSpeciesAbundance(ctx, org, rows)
				if err != nil {
					return SyncStats{}, fmt.Errorf("Failed to update %s: %w", name, err)
				}
				unchanged = uint(len(rows)) - inserts - updates
			case fieldseeker.LayerStormDrain:
				name = "LayerStormDrain"
				rows, err := fssync.StormDrain(ctx, offset)
				if err != nil {
					return SyncStats{}, fmt.Errorf("Failed to query %s: %w", name, err)
				}
				inserts, updates, err = db.SaveOrUpdateStormDrain(ctx, org, rows)
				if err != nil {
					return SyncStats{}, fmt.Errorf("Failed to update %s: %w", name, err)
				}
				unchanged = uint(len(rows)) - inserts - updates
			case fieldseeker.LayerTracklog:
				name = "LayerTracklog"
				rows, err := fssync.Tracklog(ctx, offset)
				if err != nil {
					return SyncStats{}, fmt.Errorf("Failed to query %s: %w", name, err)
				}
				inserts, updates, err = db.SaveOrUpdateTracklog(ctx, org, rows)
				if err != nil {
					return SyncStats{}, fmt.Errorf("Failed to update %s: %w", name, err)
				}
				unchanged = uint(len(rows)) - inserts - updates
			case fieldseeker.LayerTrapLocation:
				name = "LayerTrapLocation"
				rows, err := fssync.TrapLocation(ctx, offset)
				if err != nil {
					return SyncStats{}, fmt.Errorf("Failed to query %s: %w", name, err)
				}
				inserts, updates, err = db.SaveOrUpdateTrapLocation(ctx, org, rows)
				if err != nil {
					return SyncStats{}, fmt.Errorf("Failed to update %s: %w", name, err)
				}
				unchanged = uint(len(rows)) - inserts - updates
			case fieldseeker.LayerTrapData:
				name = "LayerTrapData"
				rows, err := fssync.TrapData(ctx, offset)
				if err != nil {
					return SyncStats{}, fmt.Errorf("Failed to query %s: %w", name, err)
				}
				inserts, updates, err = db.SaveOrUpdateTrapData(ctx, org, rows)
				if err != nil {
					return SyncStats{}, fmt.Errorf("Failed to update %s: %w", name, err)
				}
				unchanged = uint(len(rows)) - inserts - updates
			case fieldseeker.LayerTimeCard:
				name = "LayerTimeCard"
				rows, err := fssync.TimeCard(ctx, offset)
				if err != nil {
					return SyncStats{}, fmt.Errorf("Failed to query %s: %w", name, err)
				}
				inserts, updates, err = db.SaveOrUpdateTimeCard(ctx, org, rows)
				if err != nil {
					return SyncStats{}, fmt.Errorf("Failed to update %s: %w", name, err)
				}
				unchanged = uint(len(rows)) - inserts - updates
			case fieldseeker.LayerTreatment:
				name = "LayerTreatment"
				rows, err := fssync.Treatment(ctx, offset)
				if err != nil {
					return SyncStats{}, fmt.Errorf("Failed to query %s: %w", name, err)
				}
				inserts, updates, err = db.SaveOrUpdateTreatment(ctx, org, rows)
				if err != nil {
					return SyncStats{}, fmt.Errorf("Failed to update %s: %w", name, err)
				}
				unchanged = uint(len(rows)) - inserts - updates
			case fieldseeker.LayerTreatmentArea:
				name = "LayerTreatmentArea"
				rows, err := fssync.TreatmentArea(ctx, offset)
				if err != nil {
					return SyncStats{}, fmt.Errorf("Failed to query %s: %w", name, err)
				}
				inserts, updates, err = db.SaveOrUpdateTreatmentArea(ctx, org, rows)
				if err != nil {
					return SyncStats{}, fmt.Errorf("Failed to update %s: %w", name, err)
				}
				unchanged = uint(len(rows)) - inserts - updates
			case fieldseeker.LayerULVSprayRoute:
				name = "LayerULVSprayRoute"
				rows, err := fssync.ULVSprayRoute(ctx, offset)
				if err != nil {
					return SyncStats{}, fmt.Errorf("Failed to query %s: %w", name, err)
				}
				inserts, updates, err = db.SaveOrUpdateULVSprayRoute(ctx, org, rows)
				if err != nil {
					return SyncStats{}, fmt.Errorf("Failed to update %s: %w", name, err)
				}
				unchanged = uint(len(rows)) - inserts - updates
			case fieldseeker.LayerZones:
				name = "LayerZones"
				rows, err := fssync.Zones(ctx, offset)
				if err != nil {
					return SyncStats{}, fmt.Errorf("Failed to query %s: %w", name, err)
				}
				inserts, updates, err = db.SaveOrUpdateZones(ctx, org, rows)
				if err != nil {
					return SyncStats{}, fmt.Errorf("Failed to update %s: %w", name, err)
				}
				unchanged = uint(len(rows)) - inserts - updates
			case fieldseeker.LayerZones2:
				name = "LayerZones2"
				rows, err := fssync.Zones2(ctx, offset)
				if err != nil {
					return SyncStats{}, fmt.Errorf("Failed to query %s: %w", name, err)
				}
				inserts, updates, err = db.SaveOrUpdateZones2(ctx, org, rows)
				if err != nil {
					return SyncStats{}, fmt.Errorf("Failed to update %s: %w", name, err)
				}
				unchanged = uint(len(rows)) - inserts - updates
			default:
				return ss, errors.New("Unrecognized layer")
			}
			ss.Inserts = inserts
			ss.Updates = updates
			ss.Unchanged = unchanged
			return ss, err
		})
	}
	//log.Info().Uint("inserts", stats.Inserts).Uint("updates", stats.Updates).Uint("no change", stats.Unchanged).Str("layer", layer.Name).Msg("Finished layer")
	return stats, nil
}
