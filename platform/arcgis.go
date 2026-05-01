package platform

import (
	"bytes"
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
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
	"github.com/Gleipnir-Technology/bob"
	"github.com/Gleipnir-Technology/bob/dialect/psql"
	"github.com/Gleipnir-Technology/bob/dialect/psql/dialect"
	"github.com/Gleipnir-Technology/bob/dialect/psql/dm"
	"github.com/Gleipnir-Technology/bob/dialect/psql/im"
	"github.com/Gleipnir-Technology/nidus-sync/config"
	"github.com/Gleipnir-Technology/nidus-sync/db"
	"github.com/Gleipnir-Technology/nidus-sync/db/enums"
	"github.com/Gleipnir-Technology/nidus-sync/db/gen/nidus-sync/arcgis/model"
	"github.com/Gleipnir-Technology/nidus-sync/db/models"
	queryarcgis "github.com/Gleipnir-Technology/nidus-sync/db/query/arcgis"
	"github.com/Gleipnir-Technology/nidus-sync/db/sql"
	"github.com/Gleipnir-Technology/nidus-sync/db/types"
	"github.com/Gleipnir-Technology/nidus-sync/debug"
	"github.com/Gleipnir-Technology/nidus-sync/h3utils"
	"github.com/Gleipnir-Technology/nidus-sync/platform/oauth"
	"github.com/aarondl/opt/omit"
	"github.com/aarondl/opt/omitnull"
	"github.com/alitto/pond/v2"
	"github.com/jackc/pgx/v5"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/uber/h3-go/v4"
)

var syncStatusByOrg map[int32]bool

var CodeVerifier string = "random_secure_string_min_43_chars_long_should_be_stored_in_session"

func HasFieldseekerConnection(ctx context.Context, user_id int32) (bool, error) {
	result, err := queryarcgis.OAuthTokenForUserExists(ctx, int64(user_id))
	if err != nil {
		return false, err
	}
	return *result, nil
}

func IsSyncOngoing(org_id int32) bool {
	return syncStatusByOrg[org_id]
}
func getOAuthForOrg(ctx context.Context, org *models.Organization) (*model.OAuthToken, error) {
	users, err := org.User().All(ctx, db.PGInstance.BobDB)
	if err != nil {
		return nil, fmt.Errorf("Failed to query all users for org: %w", err)
	}
	for _, user := range users {
		oauths, err := queryarcgis.OAuthTokensForUser(ctx, int64(user.ID))
		if err != nil {
			return nil, fmt.Errorf("Failed to query all oauth tokens for org: %w", err)
		}
		for _, oauth := range oauths {
			return oauth, nil
		}
	}
	return nil, nil
}

// This is a goroutine that is in charge of getting Fieldseeker data and keeping it fresh.
func refreshFieldseekerData(background_ctx context.Context, newOauthCh <-chan struct{}) {
	ctx := log.With().Str("component", "arcgis").Logger().Level(zerolog.InfoLevel).WithContext(background_ctx)
	syncStatusByOrg = make(map[int32]bool, 0)
	for {
		workerCtx, cancel := context.WithCancel(context.Background())
		var wg sync.WaitGroup

		oauths, err := queryarcgis.OAuthTokensValid(ctx)
		if err != nil {
			log.Error().Err(err).Msg("Failed to get oauths")
			return
		}
		if len(oauths) == 0 {
			log.Info().Msg("No oauths to maintain")
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
		if len(orgs) == 0 {
			log.Info().Msg("No orgs to maintain")
		}
		for _, org := range orgs {
			wg.Add(1)
			go func() {
				defer wg.Done()
				err := periodicallyExportFieldseeker(workerCtx, org)
				if err != nil {
					log.Error().Err(err).Msg("Crashed fieldseeker export goroutine")
				}
			}()
		}

		select {
		case <-ctx.Done():
			log.Debug().Msg("Exiting arcgis refresh worker...")
			cancel()
			wg.Wait()
			log.Debug().Msg("arcgis refresh worker exited.")
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
	layers, err := fieldseekerClient.Layers(ctx)
	if err != nil {
		log.Error().Err(err).Msg("Failed to get layers")
		return
	}
	log.Debug().Int("len", len(layers)).Msg("Downloading fieldseeker schema")
	for i, layer := range layers {
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
		schema, err := fieldseekerClient.SchemaRaw(ctx, uint(i))
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
func updateArcgisUserData(ctx context.Context, user *models.User, oauth *model.OAuthToken) {
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

	txn, err := db.PGInstance.BobDB.BeginTx(ctx, nil)
	if err != nil {
		log.Error().Err(err).Msg("Create transaction")
		return
	}
	defer txn.Rollback(ctx)

	account, ag_user, err := updateArcgisAccount(ctx, txn, client, user)
	if err != nil {
		log.Error().Err(err).Msg("Failed to get portal data")
		return
	}

	err = updateServiceData(ctx, txn, client, user, account)
	if err != nil {
		log.Error().Err(err).Msg("Failed to get service data")
		return
	}

	model := model.OAuthToken{
		ArcgisID:            &ag_user.ID,
		ArcgisLicenseTypeID: &ag_user.UserLicenseTypeID,
	}
	err = queryarcgis.OAuthTokenUpdateLicense(ctx, oauth.RefreshToken, &model)
	if err != nil {
		log.Error().Err(err).Msg("Failed to update oauth token portal data")
		return
	}
	org := user.R.Organization
	if org.ArcgisAccountID.IsNull() {
		err = org.Update(ctx, txn, &models.OrganizationSetter{
			ArcgisAccountID: omitnull.From(ag_user.OrgID),
		})
		if err != nil {
			log.Error().Err(err).Int32("id", user.R.Organization.ID).Msg("Failed to update organization's arcgis info")
			return
		}
		log.Info().Int32("org_id", org.ID).Str("arcgis_id", ag_user.OrgID).Msg("Updated org arcgis ID")
	}

	fssync, err := fieldseeker.NewFieldSeekerFromAG(ctx, *client)
	if err != nil {
		log.Error().Err(err).Msg("Failed to create fieldseeker")
		return
	}
	log.Info().Str("url", fssync.ServiceFeature.URL.String()).Msg("Found Fieldseeker")

	// Ensure the fieldseeker service is saved on the account
	// Why yes, we do get 'ArcGIS' and 'arcgis' from the API, why do you ask?
	url_corrected := strings.Replace(fssync.ServiceFeature.URL.String(), "/arcgis/", "/ArcGIS/", 1)
	service_account, err := queryarcgis.ServiceFeatureFromURL(ctx, url_corrected)
	if err != nil {
		log.Error().Err(err).Str("url", fssync.ServiceFeature.URL.String()).Str("url_corrected", url_corrected).Msg("no fieldseeker service to link, it should have been created before")
		return
	}
	setter := models.OrganizationSetter{
		FieldseekerServiceFeatureItemID: omitnull.From(service_account.ItemID),
	}
	err = org.Update(ctx, db.PGInstance.BobDB, &setter)
	if err != nil {
		log.Error().Err(err).Msg("Failed to create new organization")
		return
	}
	maybeCreateWebhook(ctx, fssync)
	downloadFieldseekerSchema(ctx, fssync, account.ID)
	//notification.ClearOauth(ctx, user)
	newOAuthTokenChannel <- struct{}{}
}

func newFieldSeeker(ctx context.Context, oa *model.OAuthToken) (*fieldseeker.FieldSeeker, error) {
	if oa == nil {
		return nil, fmt.Errorf("no oath token")
	}
	row, err := sql.OrgByOauthId(oa.ID).One(ctx, db.PGInstance.BobDB)
	if err != nil {
		return nil, fmt.Errorf("Failed to get org ID from oauth %d: %w", oa.ID, err)
	}
	// The URL for fieldseeker should be something like
	// https://foo.arcgis.com/123abc/arcgis/rest/services/FieldSeekerGIS/FeatureServer
	// We need to break it up
	host, pathParts, err := extractURLParts(row.FieldseekerURL)
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
			AccessToken:         oa.AccessToken,
			AccessTokenExpires:  oa.AccessTokenExpires,
			RefreshToken:        oa.RefreshToken,
			RefreshTokenExpires: oa.RefreshTokenExpires,
		},
	)
	if err != nil {
		if errors.Is(err, arcgis.ErrorInvalidAuthToken) {
			return nil, oauth.InvalidatedTokenError{}
		} else if errors.Is(err, arcgis.ErrorInvalidRefreshToken) {
			return nil, oauth.InvalidatedTokenError{}
		}
		return nil, fmt.Errorf("Failed to create ArcGIS client: %w", err)
	}
	log.Info().Str("context", context).Str("host", host).Msg("Using base fieldseeker URL")
	fssync, err := fieldseeker.NewFieldSeekerFromURL(ctx, *ar, row.FieldseekerURL)
	if err != nil {
		return nil, fmt.Errorf("Failed to create Fieldseeker client: %w", err)
	}
	return fssync, nil
}
func updateArcgisAccount(ctx context.Context, txn bob.Tx, client *arcgis.ArcGIS, user *models.User) (*model.Account, *model.User, error) {
	p, err := client.PortalsSelf(ctx)
	if err != nil {
		return nil, nil, fmt.Errorf("Failed to get ArcGIS user data: %w", err)
	}

	// Ensure that an arcgis account exists to attach to
	account, err := ensureArcgisAccount(ctx, txn, p, user)
	ag_user, err := queryarcgis.UserFromID(ctx, p.User.ID)
	if err != nil {
		log.Warn().Err(err).Msg("need arcgis user account?")
		if err.Error() == "sql: no rows in result set" {
			setter := model.User{
				Access:            p.Access,
				Created:           time.Unix(p.User.Created, 0),
				Email:             p.User.Email,
				FullName:          p.User.FullName,
				ID:                p.User.ID,
				Level:             p.User.Level,
				OrgID:             p.User.OrgID,
				PublicUserID:      user.ID,
				Region:            p.Region,
				Role:              p.User.Role,
				RoleID:            p.User.RoleId,
				Username:          p.User.Username,
				UserLicenseTypeID: p.User.UserLicenseTypeID,
				UserType:          p.User.UserType,
			}
			ag_user, err = queryarcgis.UserInsert(ctx, txn, &setter)
			if err != nil {
				return nil, nil, fmt.Errorf("Failed to add arcgis user data: %w", err)
			}
		} else {
			return nil, nil, fmt.Errorf("Failed to find arcgis user: %w", err)
		}
	}

	err = queryarcgis.UserPrivilegesDeleteByUserID(ctx, txn, p.User.ID)

	if err != nil {
		return nil, nil, fmt.Errorf("Failed to delete previous user privilege data: %w", err)
	}

	for _, priv := range p.User.Privileges {
		s := model.UserPrivilege{
			Privilege: priv,
			UserID:    p.User.ID,
		}
		err := queryarcgis.UserPrivilegeInsert(ctx, txn, &s)
		if err != nil {
			return nil, nil, fmt.Errorf("Failed to add arcgis user privilege data: %w", err)
		}
	}
	log.Info().Str("username", p.User.Username).Str("user_id", p.User.ID).Str("org_id", p.User.OrgID).Str("org_name", p.Name).Str("license_type_id", p.User.UserLicenseTypeID).Msg("Updated portals data")
	return account, ag_user, nil
}
func updateServiceData(ctx context.Context, txn bob.Tx, client *arcgis.ArcGIS, user *models.User, account *model.Account) error {
	service_maps, err := client.MapServices(ctx)
	if err != nil {
		return fmt.Errorf("list map services: %w", err)
	}
	for _, sm := range service_maps {
		log.Info().Str("account-id", account.ID).Str("arcgis-id", sm.ID).Str("name", sm.Name).Str("title", sm.Title).Str("url", sm.URL.String()).Msg("inserting map service")
		_, err := queryarcgis.ServiceMapFromID(ctx, sm.ID)
		if err != nil {
			if err.Error() == "sql: no rows in result set" {
				setter := model.ServiceMap{
					AccountID: account.ID,
					ArcgisID:  sm.ID,
					Name:      sm.Name,
					Title:     sm.Title,
					URL:       sm.URL.String(),
				}
				err := queryarcgis.ServiceMapInsert(ctx, txn, &setter)
				if err != nil {
					return fmt.Errorf("save map service: %w", err)
				}
				_, err = models.TileServices.Insert(&models.TileServiceSetter{
					Name:     omit.From(sm.Name),
					ArcgisID: omitnull.From(sm.ID),
				}).One(ctx, txn)
				if err != nil {
					return fmt.Errorf("save tile service: %w", err)
				}
			} else {
				return err
			}
		}
	}

	services, err := client.Services(ctx)
	for _, service := range services {
		err := ensureServiceFeature(ctx, txn, client, user, account, service)
		if err != nil {
			return fmt.Errorf("ensure service feature: %w", err)
		}
	}
	return nil
}
func ensureServiceFeature(ctx context.Context, txn bob.Tx, client *arcgis.ArcGIS, user *models.User, account *model.Account, service *arcgis.ServiceFeature) error {
	_, err := queryarcgis.ServiceFeatureFromURL(ctx, service.URL.String())
	if err == nil {
		return nil
	}
	if err.Error() != "sql: no rows in result set" {
		return err
	}
	metadata, err := service.PopulateMetadata(ctx)
	if err != nil {
		return fmt.Errorf("populate metadata: %w", err)
	}

	setter := model.ServiceFeature{
		AccountID: &account.ID,
		Extent: types.Box2D{
			XMax: 180,
			YMax: 90,
			XMin: -180,
			YMin: -90,
		},
		ItemID:           metadata.ServiceItemId,
		SpatialReference: int32(*metadata.SpatialReference.LatestWKID),
		URL:              service.URL.String(),
	}
	return queryarcgis.ServiceFeatureInsert(ctx, txn, &setter)
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
			oa, err := getOAuthForOrg(ctx, org)
			if err != nil {
				return fmt.Errorf("Failed to get oauth for org: %w", err)
			}
			if oa == nil {
				//log.Debug().Int32("org.id", org.ID).Msg("No oauth for org")
				continue
			}
			fssync, err := newFieldSeeker(ctx, oa)
			if err != nil {
				if errors.Is(err, &oauth.InvalidatedTokenError{}) {
					log.Info().Int32("org", org.ID).Msg("oauth token for org is invalid, waiting for refresh")
					continue
				}
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
	layers, err := fssync.Layers(ctx)
	if err != nil {
		return fmt.Errorf("get layers: %w", err)
	}
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
		log.Warn().Err(err).Msg("Failed to get admin info during log permissions")
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

func maintainOAuth(ctx context.Context, aot *model.OAuthToken) error {
	for {
		// Refresh from the database
		oa, err := queryarcgis.OAuthTokenFromID(ctx, int64(aot.ID))
		if err != nil {
			return fmt.Errorf("Failed to update oauth token from database: %w", err)
		}
		var accessTokenDelay time.Duration
		if oa.AccessTokenExpires.Before(time.Now()) || time.Until(oa.AccessTokenExpires) < (3*time.Second) {
			accessTokenDelay = time.Second
		} else {
			accessTokenDelay = time.Until(oa.AccessTokenExpires) - (3 * time.Second)
		}
		var refreshTokenDelay time.Duration
		if oa.RefreshTokenExpires.Before(time.Now()) || time.Until(oa.RefreshTokenExpires) < (3*time.Second) {
			refreshTokenDelay = time.Second
		} else {
			refreshTokenDelay = time.Until(oa.RefreshTokenExpires) - (3 * time.Second)
		}
		log.Info().Int("id", int(oa.ID)).Float64("seconds", accessTokenDelay.Seconds()).Msg("Need to refresh access token")
		log.Info().Int("id", int(oa.ID)).Float64("seconds", refreshTokenDelay.Seconds()).Msg("Need to refresh refresh token")
		accessTokenTicker := time.NewTicker(accessTokenDelay)
		refreshTokenTicker := time.NewTicker(refreshTokenDelay)
		select {
		case <-ctx.Done():
			return nil
		case <-accessTokenTicker.C:
			err := oauth.RefreshAccessToken(ctx, oa)
			if err != nil {
				return fmt.Errorf("Failed to refresh access token: %w", err)
			}
		case <-refreshTokenTicker.C:
			err := oauth.RefreshRefreshToken(ctx, oa)
			if err != nil {
				return fmt.Errorf("Failed to maintain refresh token: %w", err)
			}
		}
	}

}

// Mark that a given oauth token has failed. This includes a notification to
// the user.
func markTokenFailed(ctx context.Context, oauth *model.OAuthToken) {
	err := queryarcgis.OAuthTokenInvalidate(ctx, int64(oauth.ID))
	if err != nil {
		log.Error().Str("err", err.Error()).Msg("Failed to mark token failed")
	}
	/*
		user, err := models.FindUser(ctx, db.PGInstance.BobDB, oauth.UserID)
		if err != nil {
			log.Error().Str("err", err.Error()).Msg("Failed to get oauth user")
			return
		}
		notification.NotifyOauthInvalid(ctx, user)
	*/
	log.Info().Int("id", int(oauth.ID)).Msg("Marked oauth token invalid")
}

func newTimestampedFilename(prefix, suffix string) string {
	timestamp := time.Now().Format("20060102_150405") // YYYYMMDD_HHMMSS format
	return prefix + timestamp + suffix
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

func saveOrUpdateDBRecords(ctx context.Context, table string, qr *response.QueryResult, org_id int32) (int, int, error) {
	inserts, updates := 0, 0
	sorted_columns := make([]string, 0, len(qr.Fields))
	for _, f := range qr.Fields {
		sorted_columns = append(sorted_columns, *f.Name)
	}
	sort.Strings(sorted_columns)

	objectids := make([]int, 0)
	for _, l := range qr.Features {
		attr := l.Attributes["OBJECTID"]
		attr_s := attr.String()
		oid, err := strconv.Atoi(attr_s)
		if err != nil {
			log.Warn().Str("attr_s", attr_s).Msg("failed to convert")
			continue
		}
		objectids = append(objectids, oid)
	}

	rows_by_objectid, err := rowmapViaQuery(ctx, table, sorted_columns, objectids)
	if err != nil {
		return inserts, updates, fmt.Errorf("Failed to get existing rows: %w", err)
	}
	// log.Println("Rows from query", len(rows_by_objectid))

	for _, feature := range qr.Features {
		attr := feature.Attributes["OBJECTID"]
		attr_s := attr.String()
		oid, err := strconv.Atoi(attr_s)
		if err != nil {
			log.Warn().Str("attr_s", attr_s).Msg("failed to convert")
			continue
		}
		row := rows_by_objectid[oid]
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
	copy(columnNames, sorted_columns)
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

func insertRowFromFeature(ctx context.Context, table string, sorted_columns []string, feature *response.Feature, org_id int32) error {
	txn, err := db.PGInstance.BobDB.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("Unable to start transaction")
	}
	defer txn.Rollback(ctx)

	err = insertRowFromFeatureFS(ctx, txn, table, sorted_columns, feature, org_id)
	if err != nil {
		return fmt.Errorf("Unable to insert FS: %w", err)
	}

	err = insertRowFromFeatureHistory(ctx, txn, table, sorted_columns, feature, org_id, 1)
	if err != nil {
		return fmt.Errorf("Failed to insert history: %w", err)
	}

	txn.Commit(ctx)
	if err != nil {
		return fmt.Errorf("Failed to commit transaction: %w", err)
	}
	return nil
}

func insertRowFromFeatureFS(ctx context.Context, txn bob.Tx, table string, sorted_columns []string, feature *response.Feature, org_id int32) error {
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

	_, err := txn.ExecContext(ctx, sb.String(), args)
	if err != nil {
		return fmt.Errorf("Failed to insert row into %s: %w", table, err)
	}
	return nil
}
func hasUpdates(row map[string]string, feature response.Feature) bool {
	return false
	/*
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
			if featureAsString, ok := value.(response.TextValue); ok {
				if featureAsString.String() != rowdata {
					return true
				}
				continue
			} else if featureAsInt, ok := value.(response.Int32Value); ok {
				// Previously had a nil value, now we have a real value
				if rowdata == "" {
					return true
				}
				rowAsInt, err := strconv.Atoi(rowdata)
				if err != nil {
					log.Error().Msg(fmt.Sprintf("Failed to convert '%s' to an int to compare against %v for %v", rowdata, featureAsInt, key))
				}
				if rowAsInt != featureAsInt.V {
					return true
				} else {
					continue
				}
			} else if featureAsFloat, ok := value.(Float64Value); ok {
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
	*/
}
func updateRowFromFeature(ctx context.Context, table string, sorted_columns []string, feature *response.Feature, org_id int32) error {
	return nil
	/*
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

		txn, err := db.PGInstance.BobDB.BeginTx(ctx, nil)
		if err != nil {
			return fmt.Errorf("Unable to start transaction")
		}
		defer txn.Rollback(ctx)

		err = insertRowFromFeatureHistory(ctx, txn, table, sorted_columns, feature, org_id, version+1)
		if err != nil {
			return fmt.Errorf("Failed to insert history: %w", err)
		}
		err = updateRowFromFeatureFS(ctx, txn, table, sorted_columns, feature)
		if err != nil {
			return fmt.Errorf("Failed to update row from feature: %w", err)
		}

		txn.Commit(ctx)
		return nil
	*/
}
func insertRowFromFeatureHistory(ctx context.Context, transaction bob.Tx, table string, sorted_columns []string, feature *response.Feature, org_id int32, version int) error {
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
	if _, err := transaction.ExecContext(ctx, sb.String(), args); err != nil {
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
	return "History_" + table[3:]
}

func updateRowFromFeatureFS(ctx context.Context, transaction bob.Tx, table string, sorted_columns []string, feature *response.Feature) error {
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

	_, err := transaction.ExecContext(ctx, sb.String(), args)
	if err != nil {
		return fmt.Errorf("Failed to update row into %s: %w", table, err)
	}
	return nil
}

func exportFieldseekerLayer(ctx context.Context, group pond.ResultTaskGroup[SyncStats], org *models.Organization, fssync *fieldseeker.FieldSeeker, layer response.Layer) (SyncStats, error) {
	var stats SyncStats
	return stats, nil
	/*
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
	*/
}

func ensureArcgisAccount(ctx context.Context, txn bob.Tx, portal *response.Portal, user *models.User) (*model.Account, error) {
	account, err := queryarcgis.AccountFromID(ctx, portal.User.OrgID)
	if err != nil {
		log.Warn().Err(err).Msg("need arcgis account?")
		if err.Error() == "sql: no rows in result set" {
			setter := model.Account{
				ID:             portal.User.OrgID,
				Name:           portal.Name,
				OrganizationID: user.OrganizationID,
				URLFeatures:    nil,
				URLInsights:    nil,
				URLGeometry:    nil,
				URLNotebooks:   nil,
				URLTiles:       nil,
			}
			account, err = queryarcgis.AccountInsert(ctx, txn, &setter)
			if err != nil {
				return nil, fmt.Errorf("create arcgis account: %w", err)
			}
		} else {
			return nil, fmt.Errorf("find arcgis account: %w", err)
		}
	}
	return account, nil
}
func updateSummaryTables(ctx context.Context, org *models.Organization) {
	updateSummaryMosquitoSource(ctx, org)
	updateSummaryServiceRequest(ctx, org)
	updateSummaryTrap(ctx, org)
}

func aggregateAtResolution(ctx context.Context, resolution int, org_id int32, type_ enums.H3aggregationtype, cells []h3.Cell) error {
	var err error
	log.Debug().Int("resolution", resolution).Str("type", string(type_)).Msg("Working summary layer")
	cellToCount := make(map[h3.Cell]int, 0)
	for _, cell := range cells {
		scaled, err := cell.Parent(resolution)
		if err != nil {
			log.Error().Err(err).Int("resolution", resolution).Msg("Failed to get cell's parent at resolution")
			continue
		}
		cellToCount[scaled] = cellToCount[scaled] + 1
	}

	_, err = models.H3Aggregations.Delete(
		dm.Where(
			psql.And(
				models.H3Aggregations.Columns.OrganizationID.EQ(psql.Arg(org_id)),
				models.H3Aggregations.Columns.Resolution.EQ(psql.Arg(resolution)),
				models.H3Aggregations.Columns.Type.EQ(psql.Arg(type_)),
			),
		),
	).Exec(ctx, db.PGInstance.BobDB)
	if err != nil {
		return fmt.Errorf("Failed to clear previous aggregation: %w", err)
	}
	var to_insert = make([]bob.Mod[*dialect.InsertQuery], 0)
	to_insert = append(to_insert, im.Into("h3_aggregation", "cell", "resolution", "count_", "type_", "organization_id", "geometry"))
	for cell, count := range cellToCount {
		polygon, err := h3utils.CellToPostgisGeometry(cell)
		if err != nil {
			log.Error().Err(err).Msg("Failed to get PostGIS geometry")
			continue
		}
		// log.Info().Str("polygon", polygon).Msg("Going to insert")
		to_insert = append(to_insert, im.Values(psql.Arg(cell.String(), resolution, count, type_, org_id), psql.F("st_geomfromtext", psql.S(polygon), 4326)))
	}
	to_insert = append(to_insert, im.OnConflict("cell, organization_id, type_").DoUpdate(
		im.SetCol("count_").To(psql.Raw("EXCLUDED.count_")),
	))
	//log.Info().Str("sql", insertQueryToString(psql.Insert(to_insert...))).Msg("Updating...")
	_, err = psql.Insert(to_insert...).Exec(ctx, db.PGInstance.BobDB)
	if err != nil {
		return fmt.Errorf("Failed to add h3 aggregation: %w", err)
	}
	return nil
}

func updateSummaryMosquitoSource(ctx context.Context, org *models.Organization) {
	point_locations, err := org.Pointlocations().All(ctx, db.PGInstance.BobDB)
	if err != nil {
		log.Error().Err(err).Msg("Failed to get all point locations")
		return
	}
	if len(point_locations) == 0 {
		log.Info().Int("org_id", int(org.ID)).Msg("No updates to perform")
		return
	}

	cells := make([]h3.Cell, 0)
	for _, p := range point_locations {
		if p.H3cell.IsNull() {
			continue
		}
		cell, err := h3utils.ToCell(p.H3cell.MustGet())
		if err != nil {
			log.Error().Err(err).Msg("Failed to get geometry point")
			continue
		}
		cells = append(cells, cell)
	}

	for i := range 16 {
		err = aggregateAtResolution(ctx, i, org.ID, enums.H3aggregationtypeMosquitosource, cells)
		if err != nil {
			log.Error().Err(err).Int("resolution", i).Msg("Failed to aggregate mosquito source")
		}
	}
}

func updateSummaryServiceRequest(ctx context.Context, org *models.Organization) {
	service_requests, err := org.Servicerequests().All(ctx, db.PGInstance.BobDB)
	if err != nil {
		log.Error().Err(err).Msg("Failed to get all service requests")
		return
	}
	if len(service_requests) == 0 {
		log.Info().Int("org_id", int(org.ID)).Msg("No updates to perform")
		return
	}

	cells := make([]h3.Cell, 0)
	for _, p := range service_requests {
		if p.H3cell.IsNull() {
			continue
		}
		cell, err := h3utils.ToCell(p.H3cell.MustGet())
		if err != nil {
			log.Error().Err(err).Msg("Failed to get geometry point")
			continue
		}
		cells = append(cells, cell)
	}
	for i := range 16 {
		err = aggregateAtResolution(ctx, i, org.ID, enums.H3aggregationtypeServicerequest, cells)
		if err != nil {
			log.Error().Err(err).Int("resolution", i).Msg("Failed to aggregate service request")
		}
	}
}

func updateSummaryTrap(ctx context.Context, org *models.Organization) {
	traps, err := org.Traplocations().All(ctx, db.PGInstance.BobDB)
	if err != nil {
		log.Error().Err(err).Msg("Failed to get all trap locations")
		return
	}
	if len(traps) == 0 {
		log.Info().Int("org_id", int(org.ID)).Msg("No updates to perform")
		return
	}

	cells := make([]h3.Cell, 0)
	for _, t := range traps {
		if t.H3cell.IsNull() {
			continue
		}
		cell, err := h3utils.ToCell(t.H3cell.MustGet())
		if err != nil {
			log.Error().Err(err).Msg("Failed to get geometry point")
			continue
		}
		cells = append(cells, cell)
	}
	for i := range 16 {
		err = aggregateAtResolution(ctx, i, org.ID, enums.H3aggregationtypeTrap, cells)
		if err != nil {
			log.Error().Err(err).Int("resolution", i).Msg("Failed to aggregate trap")
		}
	}
}
