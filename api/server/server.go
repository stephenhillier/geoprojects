package server

import (
	"github.com/go-chi/chi"
	"github.com/stephenhillier/geoprojects/api/db"
)

// Server represents the server environment (db and router)
type Server struct {
	router    chi.Router
	datastore db.Datastore
	config    Config
}

// Config holds server/database/auth service configuration
type Config struct {
	AuthCert         PEMCert // defined in auth.go
	AuthAudience     string
	AuthIssuer       string
	AuthJWKSEndpoint string
	DBConn           string
	DBDriver         string
	DefaultPageLimit int
	MaxPageLimit     int
	AuthGroupClaim   string
	AuthRoleClaim    string
}
