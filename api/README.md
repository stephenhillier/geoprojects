# Earthworks HTTP API

Backend HTTP API service

Running the server locally: `go run .`

To build/run the docker image: `docker build -t IMAGE_NAME .` then `docker run IMAGE_NAME -p 8000:8000`.

## Database

This service requires a database.  It is currently set up to look for a PostgreSQL server at localhost:5432, but it is also configured to use a Google Cloud SQL postgres instance.  Please see the `/kubernetes` dir from the base directory for a working configuration.
