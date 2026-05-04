package auth

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/Gleipnir-Technology/nidus-sync/platform"
	"github.com/rs/zerolog/log"
	"golang.org/x/crypto/bcrypt"
)

type InactiveUser struct{}

func (e InactiveUser) Error() string { return "That user is not active" }

type InvalidCredentials struct{}

func (e InvalidCredentials) Error() string { return "No username with that password exists" }

type InvalidUsername struct{}

func (e InvalidUsername) Error() string { return "That username doesn't exist" }

type NoCredentialsError struct{}

func (e NoCredentialsError) Error() string { return "No credentials were present in the request" }

type AuthenticatedHandler func(http.ResponseWriter, *http.Request, platform.User)
type EnsureAuth struct {
	handler AuthenticatedHandler
}

func AddUserSession(ctx context.Context, user *platform.User) {
	id_str := strconv.Itoa(int(user.ID))
	sessionManager.Put(ctx, "user_id", id_str)
	sessionManager.Put(ctx, "username", user.Username)
	log.Debug().Str("id", id_str).Str("username", user.Username).Msg("added user session")
}
func ImpersonateEnd(ctx context.Context) {
	sessionManager.Put(ctx, "impersonated_user_id", "")
}
func ImpersonateUser(ctx context.Context, target_user_id int) {
	target_user_id_str := strconv.Itoa(int(target_user_id))
	sessionManager.Put(ctx, "impersonated_user_id", target_user_id_str)
}
func ImpersonatedUser(ctx context.Context) *int32 {
	i_str := sessionManager.GetString(ctx, "impersonated_user_id")
	if i_str == "" {
		return nil
	}
	i, err := strconv.Atoi(i_str)
	if err != nil {
		log.Error().Err(err).Str("impersonated_user_id", i_str).Msg("failed to parse impersonated_user_id")
		return nil
	}
	result := int32(i)
	return &result
}
func ImpersonatorID(ctx context.Context) *int32 {
	user_id_str := sessionManager.GetString(ctx, "user_id")
	user_id, err := strconv.Atoi(user_id_str)
	if err != nil {
		log.Error().Err(err).Str("user_id", user_id_str).Msg("failed to parse user_id")
		return nil
	}
	result := int32(user_id)
	return &result

}
func GetAuthenticatedUser(r *http.Request) (*platform.User, error) {
	ctx := r.Context()
	user_id_str := sessionManager.GetString(ctx, "user_id")
	impersonated_user_id_str := sessionManager.GetString(ctx, "impersonated_user_id")
	if impersonated_user_id_str != "" {
		user_id_str = impersonated_user_id_str
	}
	if user_id_str != "" {
		user_id, err := strconv.Atoi(user_id_str)
		if err != nil {
			return nil, fmt.Errorf("Failed to convert user_id to int: %w", err)
		}
		username := sessionManager.GetString(ctx, "username")
		if user_id > 0 && username != "" {
			user, err := platform.UserByID(ctx, int32(user_id))
			if err != nil {
				return nil, fmt.Errorf("user by ID: %w", err)
			}
			if !user.IsActive {
				return nil, fmt.Errorf("user is inactive")
			}
			return user, nil
		}
	}
	// If we can't get the user from the session try to get from auth headers
	username, password, ok := r.BasicAuth()
	if !ok {
		return nil, &NoCredentialsError{}
	}
	user, err := validateUser(ctx, username, password)
	if err != nil {
		return nil, err
	}
	AddUserSession(ctx, user)
	return user, nil
}

func NewEnsureAuth(handlerToWrap AuthenticatedHandler) *EnsureAuth {
	return &EnsureAuth{handlerToWrap}
}

func (ea *EnsureAuth) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// If this is an API request respond with a more machine-readable error state
	accept := r.Header.Get("Accept")
	/*
		offers := []string{"application/json", "text/html"}

		content_type := NegotiateContent(accept, offers)
	*/
	user, err := GetAuthenticatedUser(r)
	if err != nil || user == nil {
		var msg []byte
		// Don't send authentication headers for browsers because it forces the authentication popup
		requested_with := r.Header.Get("X-Requested-With")
		//log.Debug().Str("x-requested-with", requested_with).Send()
		if !strings.HasPrefix(requested_with, "nidus-web") && accept != "text/event-stream" {
			w.Header().Set("WWW-Authenticate", `Basic realm="Nidus Sync"`)
			// Separate return codes for different authentication failures
			if _, ok := err.(*NoCredentialsError); ok {
				log.Info().Msg("No credentials present and no session")
				w.Header().Set("WWW-Authenticate-Error", "no-credentials")
				msg = []byte("Please provide credentials.\n")
			} else if _, ok := err.(*platform.NoUserError); ok {
				w.Header().Set("WWW-Authenticate-Error", "invalid-credentials")
				msg = []byte("Invalid credentials provided.\n")
			} else if _, ok := err.(*InvalidCredentials); ok {
				w.Header().Set("WWW-Authenticate-Error", "invalid-credentials")
				msg = []byte("Invalid credentials provided.\n")
			}
		}

		w.WriteHeader(401)
		_, err = w.Write(msg)
		if err != nil {
			log.Error().Err(err).Msg("failed to write response")
		}
		return
	}
	ea.handler(w, r, *user)
}
func SigninUser(r *http.Request, username string, password string) (*platform.User, error) {
	user, err := validateUser(r.Context(), username, password)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New("No matching user")
	}
	AddUserSession(r.Context(), user)
	return user, nil
}

func SignoutUser(r *http.Request, user platform.User) {
	sessionManager.Put(r.Context(), "user_id", "")
	sessionManager.Put(r.Context(), "username", "")
	err := sessionManager.Destroy(r.Context())
	if err != nil {
		log.Error().Err(err).Msg("failed to destroy session for user on signout")
	}
	log.Info().Str("username", user.Username).Int("user_id", (user.ID)).Msg("Ended user session")
}

func SignupUser(ctx context.Context, username string, name string, password string) (*platform.User, error) {
	password_hash, err := HashPassword(password)
	if err != nil {
		return nil, fmt.Errorf("Cannot signup user, failed to create hashed password: %w", err)
	}
	u, err := platform.CreateUser(ctx, username, name, password_hash)
	if err != nil {
		return nil, fmt.Errorf("create user: %s", err)
	}
	return u, nil
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func redact(s string) string {
	if len(s) <= 4 {
		return s
	}

	first_two := s[:2]
	last_two := s[len(s)-2:]
	middle_length := len(s) - 4

	return first_two + strings.Repeat("*", middle_length) + last_two
}

func validatePassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	if err != nil {
		log.Debug().Err(err).Str("password", password).Str("hash", hash).Msg("!validate password")
	}
	return err == nil
}

func validateUser(ctx context.Context, username string, password string) (*platform.User, error) {
	passwordHash, err := HashPassword(password)
	if err != nil {
		return nil, fmt.Errorf("Failed to hash password: %w", err)
	}
	user, err := platform.UserByUsername(ctx, username)
	if err != nil {
		return nil, fmt.Errorf("Failed to query for user: %w", err)
	}
	if user == nil {
		log.Info().Str("username", username).Str("password", redact(password)).Msg("Invalid username")
		return nil, InvalidUsername{}
	}
	if !user.IsActive {
		return nil, InactiveUser{}
	}
	if !validatePassword(password, user.PasswordHash) {
		log.Info().Str("username", username).Str("password", redact(password)).Str("hash", passwordHash).Msg("Invalid password for user")
		return nil, InvalidCredentials{}
	}
	return user, nil
}
