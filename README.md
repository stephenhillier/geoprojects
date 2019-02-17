# Geoprojects
Project management platform for geotechnical engineers and geoscientists

# Folders:

`api/` - REST API for CRUD functions (creating, retrieving, updating and deleting projects, boreholes, soil sample data, etc.)

`web/` - Web frontend (Vue.JS)

`logrend/` - Borehole log PDF renderer (Node.js)

`kubernetes/` - Kubernetes resource definitions (yaml files for creating deployments, services etc.)

## API

The REST API server application is located in the `/api` folder. To build and run, use `go build` and then run the executable with `./api`. It requires access to a PostgreSQL database. The app will migrate the database (if required) and fetch public keys from the auth provider.

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

### Resources

#### Projects:
_/api/v1/projects_

A project represents a contract or job. Within Earthworks, each project works like a folder, containing data such as boreholes and soil samples.

#### Boreholes:
_/api/v1/boreholes_

Boreholes are given a name and a project reference, and there may be several soil strata records associated with a borehole record. Each borehole record references a "datapoint" table record, where location information is stored (allows multiple instruments and a borehole in one single location)

#### Soil strata:
_/api/v1/strata_

Soil strata are records of soil encountered during drilling of a borehole.

## Web

To run the web frontend, use `npm run serve` from the `web/` directory and visit `127.0.0.1:8080`.

To build for production, use `npm run build` and serve the assets from an http server. The included Dockerfile creates an nginx image with the built assets.

## Log renderer

`logrend/` contains a node.js service for generating borehole logs from project data automatically. The log renderer can be tested locally using `npm run logrend`, or can be run as a stateless service using the included Dockerfile and k8s config.
