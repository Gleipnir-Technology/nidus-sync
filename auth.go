package main

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/Gleipnir-Technology/nidus-sync/enums"
	"github.com/Gleipnir-Technology/nidus-sync/models"
	"github.com/Gleipnir-Technology/nidus-sync/sql"
	"github.com/aarondl/opt/omit"
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

func NewEnsureAuth(handlerToWrap AuthenticatedHandler) *EnsureAuth {
	return &EnsureAuth{handlerToWrap}
}

func (ea *EnsureAuth) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// If this is an API request respond with a more machine-readable error state
	accept := r.Header.Values("Accept")
	offers := []string{"application/json", "text/html"}

	content_type := NegotiateContent(accept, offers)
	user, err := getAuthenticatedUser(r)
	if err != nil {
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
func addUserSession(r *http.Request, user *models.User) {
	id := strconv.Itoa(int(user.ID))
	sessionManager.Put(r.Context(), "user_id", id)
	sessionManager.Put(r.Context(), "username", user.Username)
	slog.Info("Created new user session",
		slog.String("username", user.Username),
		slog.String("user_id", id))
}

// Helper function to translate strings into solid error types for operating on
func findUser(ctx context.Context, user_id int) (*models.User, error) {
	user, err := models.FindUser(ctx, PGInstance.BobDB, int32(user_id))
	if err != nil {
		if err.Error() == "No such user" {
			return nil, &NoUserError{}
		} else {
			LogErrorTypeInfo(err)
			slog.Error("Unrecognized error. This should be updated in the findUser code", slog.String("err", err.Error()))
			return nil, err
		}
	}
	return user, err
}

func getAuthenticatedUser(r *http.Request) (*models.User, error) {
	//user_id := sessionManager.GetInt(r.Context(), "user_id")
	user_id_str := sessionManager.GetString(r.Context(), "user_id")
	if user_id_str != "" {
		user_id, err := strconv.Atoi(user_id_str)
		if err != nil {
			return nil, fmt.Errorf("Failed to convert user_id to int: %w", err)
		}
		username := sessionManager.GetString(r.Context(), "username")
		slog.Info("Current session info",
			slog.Int("user_id", user_id),
			slog.String("username", username))
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
	addUserSession(r, user)
	return user, nil
}

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func signinUser(r *http.Request, username string, password string) (*models.User, error) {
	user, err := validateUser(r.Context(), username, password)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New("No matching user")
	}
	addUserSession(r, user)
	return user, nil
}

func signupUser(username string, name string, password string) (*models.User, error) {
	passwordHash, err := hashPassword(password)
	if err != nil {
		return nil, fmt.Errorf("Cannot signup user: %w", err)
	}
	setter := models.UserSetter{
		DisplayName:      omit.From(name),
		PasswordHash:     omit.From(passwordHash),
		PasswordHashType: omit.From(enums.HashtypeBcrypt14),
		Username:         omit.From(username),
	}
	u, err := models.Users.Insert(&setter).One(context.TODO(), PGInstance.BobDB)
	if err != nil {
		return nil, fmt.Errorf("Failed to create user: %w", err)
	}
	slog.Info("Created user",
		slog.Int("ID", int(u.ID)),
		slog.String("username", u.Username))

	return u, nil
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
	slog.Info("Validating user",
		slog.String("username", username),
		slog.String("password", password),
		slog.String("hash", passwordHash))
	result, err := sql.UserByUsername(username).All(ctx, PGInstance.BobDB)
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
