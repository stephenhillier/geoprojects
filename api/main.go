package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/namsral/flag"

	"github.com/stephenhillier/geoprojects/api/projects"
)

// server represents the server environment (db and router)
type server struct {
	router chi.Router
	config config
	apps   apps
}

// apps represents the applications available to the API.
// Applications in this list should have a Routes property with a function
// that registers API routes that the application handles.
type apps struct {
	projects *projects.App
}

// config holds server/database/auth service configuration
type config struct {
	authCert         PEMCert // defined in auth.go
	authAudience     string
	authIssuer       string
	authJWKSEndpoint string
	dbuser           string
	dbpass           string
	dbname           string
	dbhost           string
}

func main() {

	conf := config{}
	flag.StringVar(&conf.dbuser, "dbuser", "geo", "database username")
	flag.StringVar(&conf.dbpass, "dbpass", "", "database password")
	flag.StringVar(&conf.dbname, "dbname", "geo", "database name")
	flag.StringVar(&conf.dbhost, "dbhost", "127.0.0.1", "database service host")
	flag.StringVar(&conf.authAudience, "auth_audience", "api.earthworksqc.com", "authentication service audience claim")
	flag.StringVar(&conf.authIssuer, "auth_issuer", "https://earthworks.auth0.com/", "authentication service issuer claim")
	flag.StringVar(&conf.authJWKSEndpoint, "jwks_endpoint", "https://earthworks.auth0.com/.well-known/jwks.json", "authentication JWKS endpoint")

	api := &server{}
	api.config = conf

	// get new certificate when server initially starts
	// see auth.go
	cert, err := api.getCert(nil)
	if err != nil {
		log.Panic(err)
	}

	api.config.authCert = cert

	// create db connection and router and use them to create a new "Server" instance
	db, err := NewDB(fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=disable", api.config.dbuser, api.config.dbpass, api.config.dbhost, api.config.dbname))
	if err != nil {
		log.Panic(err)
	}

	api.router = chi.NewRouter()
	api.apps = apps{
		projects: projects.NewApp(db),
	}

	// register middleware
	api.router.Use(middleware.Logger)

	// register routes from routes.go
	api.routes()

	h := http.Server{Addr: ":8000", Handler: api.router}

	log.Printf("Starting HTTP server on port 8000.\n")
	log.Printf("Press CTRL+C to stop.")
	go func() {
		if err := h.ListenAndServe(); err != nil {
			log.Fatal(err)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	<-stop
	log.Println("Shutting down...")
	h.Shutdown(context.Background())
	log.Println("Server stopped.")
}

// health is a simple health check handler that returns HTTP 200 OK.
func (api *server) health(w http.ResponseWriter, req *http.Request) {
	w.WriteHeader(http.StatusOK)
}
