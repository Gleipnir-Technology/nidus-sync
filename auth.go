package main

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/aarondl/opt/omit"
	"github.com/aarondl/opt/omitnull"
	"github.com/Gleipnir-Technology/nidus-sync/enums"
	"github.com/Gleipnir-Technology/nidus-sync/models"
	"github.com/Gleipnir-Technology/nidus-sync/sql"
	"golang.org/x/crypto/bcrypt"
)

type NoCredentialsError struct{}
func (e NoCredentialsError) Error() string { return "No credentials were present in the request" }

type InvalidCredentials struct{}
func (e InvalidCredentials) Error() string { return "No username with that password exists" }

type InvalidUsername struct{}
func (e InvalidUsername) Error() string { return "That username doesn't exist" }

func addUserSession(r *http.Request, user *models.User) {
	id := strconv.Itoa(int(user.ID))
	sessionManager.Put(r.Context(), "user_id", id)
	sessionManager.Put(r.Context(), "username", user.Username)
	slog.Info("Created new user session",
		slog.String("username", user.Username),
		slog.String("user_id", id))
}

func getAuthenticatedUser(r *http.Request) (*models.User, error) {
	//user_id := sessionManager.GetInt(r.Context(), "user_id")
	user_id_str := sessionManager.GetString(r.Context(), "user_id")
	user_id, err := strconv.Atoi(user_id_str)
	if err != nil {
		return nil, fmt.Errorf("Failed to convert user_id to int: %v", err)
	}
	username := sessionManager.GetString(r.Context(), "username")
	slog.Info("Current session info",
		slog.Int("user_id", user_id),
		slog.String("username", username))
	if user_id > 0 && username != "" {
		return models.FindUser(r.Context(), PGInstance.BobDB, int32(user_id))
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
		return nil, fmt.Errorf("Cannot signup user: %v", err)
	}
	setter := models.UserSetter{
		DisplayName: omitnull.From(name),
		PasswordHash: omitnull.From(passwordHash),
		PasswordHashType: omitnull.From(enums.HashtypeBcrypt14),
		Username: omit.From(username),
	}
	u, err := models.Users.Insert(&setter).One(context.TODO(), PGInstance.BobDB)
	if err != nil {
		return nil, fmt.Errorf("Failed to create user: %v", err)
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
		return nil, fmt.Errorf("Failed to hash password: %v", err)
	}
	slog.Info("Validating user",
		slog.String("username", username),
		slog.String("password", password),
		slog.String("hash", passwordHash))
	result, err := sql.UserByUsername(username).All(ctx, PGInstance.BobDB)
	if err != nil {
		return nil, fmt.Errorf("Failed to query for user: %v", err)
	}
	switch len(result) {
		case 0:
		return nil, InvalidUsername{}
	case 1:
		row := result[0]
		hash, err := row.PasswordHash.Value()
		if err != nil {
			return nil, err
		}
		if hash == nil {
			return nil, errors.New("Hash is nil")
		}
		 hashStr, ok := hash.(string);
		if !ok {
			return nil, errors.New("Hash isn't a string")
		}
		if !validatePassword(password, hashStr) {
			return nil, InvalidCredentials{}
		}
		user := models.User{
			ID: row.ID,
			ArcgisAccessToken: row.ArcgisAccessToken,
			ArcgisLicense: row.ArcgisLicense,
			ArcgisRefreshToken: row.ArcgisRefreshToken,
			ArcgisRefreshTokenExpires: row.ArcgisRefreshTokenExpires,
			ArcgisRole: row.ArcgisRole,
			DisplayName: row.DisplayName,
			Email: row.Email,
			OrganizationID: row.OrganizationID,
			Username: row.Username,
		}
		return &user, nil
	default:
		return nil, errors.New("More than one matching row, this should be impossible.")

	}
}
