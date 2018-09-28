package main

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"os"
	"time"

	jwtmiddleware "github.com/auth0/go-jwt-middleware"
	jwt "github.com/dgrijalva/jwt-go"
)

// jwks stores a jsonWebKeys (JWK) set
type jwks struct {
	Keys []jsonWebKeys `json:"keys"`
}

// jsonWebKeys represents one key from a JWK set
type jsonWebKeys struct {
	Kty string   `json:"kty"`
	Kid string   `json:"kid"`
	Use string   `json:"use"`
	N   string   `json:"n"`
	E   string   `json:"e"`
	X5c []string `json:"x5c"`
}

// cert is a public key (PEM) that has been fetched from an auth server.
// if the key is out of date or the key id doesn't match, a new key will be
// fetched
type pemCert struct {
	Cert   string
	Kid    string
	Expiry time.Time
}

var cert pemCert

// jwtAuthentication returns a new JWTMiddleware from the auth0 go-jwt-middleware package.
// the JWTMiddleware can be used with chi middleware using jwtAuthentication().Handler
func (api *Server) jwtAuthentication() *jwtmiddleware.JWTMiddleware {
	var err error

	// get new certificate when server initially starts
	// create a new middleware
	// see https://auth0.com/docs/
	jwtMiddleware := jwtmiddleware.New(jwtmiddleware.Options{
		ValidationKeyGetter: func(token *jwt.Token) (interface{}, error) {
			// Verify 'aud' claim
			aud := os.Getenv("GEO_AUTH_AUDIENCE")
			checkAud := token.Claims.(jwt.MapClaims).VerifyAudience(aud, false)
			if !checkAud {
				return token, errors.New("invalid audience")
			}
			// Verify 'iss' claim
			iss := os.Getenv("GEO_AUTH_ISS")
			checkIss := token.Claims.(jwt.MapClaims).VerifyIssuer(iss, false)
			if !checkIss {
				return token, errors.New("invalid issuer")
			}

			// check if we need a new certificate
			if api.authCert.Cert == "" || api.authCert.Kid != token.Header["kid"] || api.authCert.Expiry.Before(time.Now()) {
				api.authCert, err = getCert(token)
				if err != nil {
					log.Panic(err)
				}
			}

			result, err := jwt.ParseRSAPublicKeyFromPEM([]byte(api.authCert.Cert))
			if err != nil {
				log.Panic(err)
			}
			return result, nil
		},
		SigningMethod: jwt.SigningMethodRS256,
	})
	return jwtMiddleware
}

// getCert makes a request to the jwks endpoint and returns a public key certificate
// original code from from auth0.com/docs/
func getCert(token *jwt.Token) (pemCert, error) {

	// create a new PEM certificate `newCert`.
	// it will not be returned unless we successfully populate it.
	// if function returns an error, the existing `cert` certificate will be returned.
	newCert := pemCert{}
	host := os.Getenv("GEO_AUTH_HOST_JWKS")

	// make a request to the JWKS endpoint specified in `host` above
	log.Println("Fetching new JWKS from", os.Getenv("GEO_AUTH_ISS"))
	response, err := http.Get(host)
	if err != nil {
		return cert, err
	}

	defer response.Body.Close()

	// decode response as a jwks with a set of keys
	var jwks = jwks{}
	err = json.NewDecoder(response.Body).Decode(&jwks)
	if err != nil {
		return cert, err
	}

	// find a key matching the token.
	// if this function is called with no token, only the first key is cached.
	// if a future token requires a different kid (signing key was rotated), then the
	// cached certificate will not match and the JWKS endpoint will be retried.
keys:
	for k := range jwks.Keys {
		if token == nil || token.Header["kid"] == jwks.Keys[k].Kid {
			// set a new certificate and mark expiry for 24 hours from now
			newCert.Cert = "-----BEGIN CERTIFICATE-----\n" + jwks.Keys[k].X5c[0] + "\n-----END CERTIFICATE-----"
			newCert.Kid = jwks.Keys[k].Kid
			newCert.Expiry = time.Now().Add(24 * time.Hour)
			log.Println("New auth certificate obtained.")
			break keys
		}
	}

	if newCert.Cert == "" {
		err := errors.New("unable to find appropriate key")
		// return previously cached cert.
		// this may happen if user has a token signed by an old key that was rotated out.
		return cert, err
	}

	return newCert, nil
}

// todo: add pemCert type receiver to getCert function (e.g. call by api.authCert.Update())
// - also fix use of cert global variable