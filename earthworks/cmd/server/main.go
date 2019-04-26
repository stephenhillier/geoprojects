package main

import "log"

// Config holds server/database/auth service configuration
type Config struct {
	AuthCert         PEMCert // defined in auth.go
	AuthAudience     string
	AuthIssuer       string
	AuthJWKSEndpoint string
	DBConn           string
	DBDriver         string
	AuthGroupClaim   string
	AuthRoleClaim    string
}

func main() {
	var a bool
	if !!a {
		log.Println("asdf!")
	}
}
