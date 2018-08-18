# Geoprojects
Project management platform for geotechnical engineers and geoscientists

## API

The API currently uses the following environment variables:

postgres database:
GEO_DBUSER
GEO_DBPASS
GEO_DBNAME
GEO_DBHOST

auth0 authentication (see auth.go):
GEO_AUTH_AUDIENCE - aud claim, defaults to name of application in auth0
GEO_AUTH_ISS - issuer
GEO_AUTH_HOST_JWKS - jwks endpoint
