package auth

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/Gleipnir-Technology/nidus-sync/db"
	"github.com/Gleipnir-Technology/nidus-sync/db/enums"
	"github.com/Gleipnir-Technology/nidus-sync/db/models"
	"github.com/Gleipnir-Technology/nidus-sync/db/sql"
	"github.com/Gleipnir-Technology/nidus-sync/debug"
	"github.com/aarondl/opt/omit"
	"github.com/aarondl/opt/omitnull"
	"github.com/rs/zerolog/log"
	"github.com/stephenafamo/bob/dialect/psql"
	"github.com/stephenafamo/bob/dialect/psql/sm"
	"golang.org/x/crypto/bcrypt"
)

type NoCredentialsError struct{}

func (e NoCredentialsError) Error() string { return "No credentials were present in the request" }

type NoUserError struct{}

func (e NoUserError) Error() string { return "That user does not exist" }

type InvalidCredentials struct{}

func (e InvalidCredentials) Error() string { return "No username with that password exists" }

type InvalidUsername struct{}

func (e InvalidUsername) Error() string { return "That username doesn't exist" }

type AuthenticatedHandler func(http.ResponseWriter, *http.Request, *models.User)
type EnsureAuth struct {
	handler AuthenticatedHandler
}

func AddUserSession(r *http.Request, user *models.User) {
	id := strconv.Itoa(int(user.ID))
	sessionManager.Put(r.Context(), "user_id", id)
	sessionManager.Put(r.Context(), "username", user.Username)
	log.Info().Str("username", user.Username).Str("user_id", id).Msg("Created new user session")
}

func GetAuthenticatedUser(r *http.Request) (*models.User, error) {
	//user_id := sessionManager.GetInt(r.Context(), "user_id")
	user_id_str := sessionManager.GetString(r.Context(), "user_id")
	if user_id_str != "" {
		user_id, err := strconv.Atoi(user_id_str)
		if err != nil {
			return nil, fmt.Errorf("Failed to convert user_id to int: %w", err)
		}
		username := sessionManager.GetString(r.Context(), "username")
		log.Info().Int("user_id", user_id).Str("username", username).Msg("Current session info")
		if user_id > 0 && username != "" {
			return findUser(r.Context(), user_id)
		}
	}
	// If we can't get the user from the session try to get from auth headers
	username, password, ok := r.BasicAuth()
	if !ok {
		return nil, &NoCredentialsError{}
	}
	user, err := validateUser(r.Context(), username, password)
	if err != nil {
		return nil, err
	}
	AddUserSession(r, user)
	return user, nil
}

func NewEnsureAuth(handlerToWrap AuthenticatedHandler) *EnsureAuth {
	return &EnsureAuth{handlerToWrap}
}

func (ea *EnsureAuth) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// If this is an API request respond with a more machine-readable error state
	accept := r.Header.Values("Accept")
	offers := []string{"application/json", "text/html"}

	content_type := NegotiateContent(accept, offers)
	user, err := GetAuthenticatedUser(r)
	if err != nil || user == nil {
		var msg []byte
		// Separate return codes for different authentication failures
		if _, ok := err.(*NoCredentialsError); ok {
			fmt.Println("No credentials present and no session")
			w.Header().Set("WWW-Authenticate-Error", "no-credentials")
			msg = []byte("Please provide credentials.\n")
		} else if _, ok := err.(*NoUserError); ok {
			w.Header().Set("WWW-Authenticate-Error", "invalid-credentials")
			msg = []byte("Invalid credentials provided.\n")
		} else if _, ok := err.(*InvalidCredentials); ok {
			w.Header().Set("WWW-Authenticate-Error", "invalid-credentials")
			msg = []byte("Invalid credentials provided.\n")
		}

		if content_type == "text/html" {
			http.Redirect(w, r, "/signin?next="+r.URL.Path, http.StatusSeeOther)
			return
		}
		w.Header().Set("WWW-Authenticate", `Basic realm="Nidus Sync"`)
		w.WriteHeader(401)
		w.Write(msg)
		return
	}

	ea.handler(w, r, user)
}
func SigninUser(r *http.Request, username string, password string) (*models.User, error) {
	user, err := validateUser(r.Context(), username, password)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New("No matching user")
	}
	AddUserSession(r, user)
	return user, nil
}

func SignupUser(ctx context.Context, username string, name string, password string) (*models.User, error) {
	passwordHash, err := hashPassword(password)
	if err != nil {
		return nil, fmt.Errorf("Cannot signup user, failed to create hashed password: %w", err)
	}
	o_setter := models.OrganizationSetter{
		Name:           omit.From(fmt.Sprintf("%s's organization", username)),
		ArcgisID:       omitnull.From(""),
		ArcgisName:     omitnull.From(""),
		FieldseekerURL: omitnull.From(""),
	}
	o, err := models.Organizations.Insert(&o_setter).One(ctx, db.PGInstance.BobDB)
	if err != nil {
		return nil, fmt.Errorf("Failed to create organization: %w", err)
	}
	log.Info().Int32("id", o.ID).Msg("Created organization")
	u_setter := models.UserSetter{
		DisplayName:      omit.From(name),
		OrganizationID:   omit.From(o.ID),
		PasswordHash:     omit.From(passwordHash),
		PasswordHashType: omit.From(enums.HashtypeBcrypt14),
		Username:         omit.From(username),
	}
	u, err := models.Users.Insert(&u_setter).One(ctx, db.PGInstance.BobDB)
	if err != nil {
		return nil, fmt.Errorf("Failed to create user: %w", err)
	}
	log.Info().Int32("id", u.ID).Str("username", u.Username).Msg("Created user")

	return u, nil
}

// Helper function to translate strings into solid error types for operating on
func findUser(ctx context.Context, user_id int) (*models.User, error) {
	//user, err := models.FindUser(ctx, db.PGInstance.BobDB, int32(user_id))
	user, err := models.Users.Query(
		models.Preload.User.Organization(),
		sm.Where(models.Users.Columns.ID.EQ(psql.Arg(user_id))),
	).One(ctx, db.PGInstance.BobDB)
	if err != nil {
		if err.Error() == "No such user" || err.Error() == "sql: no rows in result set" {
			return nil, &NoUserError{}
		} else {
			debug.LogErrorTypeInfo(err)
			log.Error().Err(err).Msg("Unrecognized error. This should be updated in the findUser code")
			return nil, err
		}
	}
	log.Info().Int32("user_id", user.ID).Int32("org_id", user.OrganizationID).Msg("Found user")
	return user, err
}

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func validatePassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func validateUser(ctx context.Context, username string, password string) (*models.User, error) {
	passwordHash, err := hashPassword(password)
	if err != nil {
		return nil, fmt.Errorf("Failed to hash password: %w", err)
	}
	log.Info().Str("username", username).Str("password", password).Str("hash", passwordHash).Msg("Validating user")
	result, err := sql.UserByUsername(username).All(ctx, db.PGInstance.BobDB)
	if err != nil {
		return nil, fmt.Errorf("Failed to query for user: %w", err)
	}
	switch len(result) {
	case 0:
		return nil, InvalidUsername{}
	case 1:
		row := result[0]
		if !validatePassword(password, row.PasswordHash) {
			return nil, InvalidCredentials{}
		}
		user := models.User{
			ID:                        row.ID,
			ArcgisAccessToken:         row.ArcgisAccessToken,
			ArcgisLicense:             row.ArcgisLicense,
			ArcgisRefreshToken:        row.ArcgisRefreshToken,
			ArcgisRefreshTokenExpires: row.ArcgisRefreshTokenExpires,
			ArcgisRole:                row.ArcgisRole,
			DisplayName:               row.DisplayName,
			Email:                     row.Email,
			OrganizationID:            row.OrganizationID,
			Username:                  row.Username,
		}
		return &user, nil
	default:
		return nil, errors.New("More than one matching row, this should be impossible.")

	}
}
