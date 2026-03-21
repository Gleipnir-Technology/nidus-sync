package auth

import (
	"time"

	"github.com/Gleipnir-Technology/nidus-sync/db"
	"github.com/alexedwards/scs/pgxstore"
	"github.com/alexedwards/scs/v2"
)

var sessionManager *scs.SessionManager

func NewSessionManager() *scs.SessionManager {
	sessionManager = scs.New()
	sessionManager.Store = pgxstore.New(db.PGInstance.PGXPool)
	sessionManager.Lifetime = 24 * time.Hour
	return sessionManager
}
