# Geoprojects
Project management platform for geotechnical engineers and geoscientists

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

#### Programs:

A program is a drilling or sampling trip undertaken to collect data as part of a project.  A program allows project field work to be divided into manageable parts (e.g. Project 1: 1999 Drilling Work)
Programs are not currently implemented.
Todo: possibly replace with tags

#### Datapoints:

Datapoints represent a spatial location where work has occured or data has been collected (e.g. a sample). Each borehole or sample will have a reference to a datapoint, which allows for several data sources (e.g. a borehole and installed instrumentation) to be associated with one location.  Todo:  possibly rename to Location.

#### Boreholes:
_/api/v1/boreholes_

Boreholes are given a name and a project reference, and there may be several soil strata records associated with a borehole record.

#### Soil strata:

Soil strata are records of soil encountered during drilling of a borehole.
