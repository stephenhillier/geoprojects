package config

import "github.com/stephenhillier/geoprojects/api/server/auth"

// Config holds server/database/auth service configuration
type Config struct {
	AuthCert         auth.PEMCert // defined in auth.go
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
