# Geoprojects
Project management platform for geotechnical engineers and geoscientists

## API

### Configuration

Configuration can be supplied by flag or environment variable (flags take precedence):

| Flag         | Env         | Default    | Description |
|--------------|-------------|------------|-------------|
|-dbuser       |DBUSER       | geo |database username |
|-dbpass       |DBPASS       |     |database password |
|-dbname       |DBNAME       | geo |database name |
|-dbhost       |DBHOST       | 127.0.0.1 |database service host |
|-auth_audience|AUTH_AUDIENCE| api.earthworksqc.com |aud claim (name of app in auth0) |
|-auth_issuer  |AUTH_ISSUER  | https://earthworks.auth0.com/ |auth issuer |
|-jwks_endpoint|JWKS_ENDPOINT| https://earthworks.auth0.com/.well-known/jwks.json |JWKS endpoint |
