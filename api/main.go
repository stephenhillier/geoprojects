package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/stephenhillier/geoprojects/api/projects"
	"github.com/stephenhillier/geoprojects/backend/models"
)

// Server represents the server environment (db and router)
type Server struct {
	router   chi.Router
	authCert pemCert // defined in auth.go
	apps     apps
}

// apps represents the applications available to the API.
// Applications in this list should have a Routes property with a function
// that registers API routes that the application handles.
type apps struct {
	projects projects.App
}

func main() {

	dbuser := os.Getenv("GEO_DBUSER")
	dbpass := os.Getenv("GEO_DBPASS")
	dbname := os.Getenv("GEO_DBNAME")
	dbhost := os.Getenv("GEO_DBHOST")

	// create db connection and router and use them to create a new "Server" instance
	db, err := NewDB(fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=disable", dbuser, dbpass, dbhost, dbname))
	if err != nil {
		log.Panic(err)
	}
	r := chi.NewRouter()

	apps := apps{
		projects: projects.NewApp(db)
	}

	// get new certificate when server initially starts
	// see auth.go
	cert, err := getCert(nil)
	if err != nil {
		log.Panic(err)
	}

	api := &Server{r, cert, apps}

	// register middleware
	api.router.Use(middleware.Logger)

	// register routes from routes.go
	api.routes()

	log.Printf("Starting HTTP server on port 8000.\n")
	log.Printf("Press CTRL+C to stop.")
	log.Fatal(http.ListenAndServe(":8000", api.router))

}

// health is a simple health check handler that returns HTTP 200 OK.
func (api *Server) health(w http.ResponseWriter, req *http.Request) {
	w.WriteHeader(http.StatusOK)
}
