package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/stephenhillier/geoprojects/backend/models"
)

// Server represents the server environment (db and router)
type Server struct {
	db       models.Datastore
	router   chi.Router
	authCert pemCert // defined in auth.go
}

func main() {

	dbuser := os.Getenv("GEO_DBUSER")
	dbpass := os.Getenv("GEO_DBPASS")
	dbname := os.Getenv("GEO_DBNAME")
	dbhost := os.Getenv("GEO_DBHOST")

	// create db connection and router and use them to create a new "Server" instance
	db, err := models.NewDB(fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=disable", dbuser, dbpass, dbhost, dbname))
	if err != nil {
		log.Panic(err)
	}
	r := chi.NewRouter()

	// get new certificate when server initially starts
	cert, err := getCert(nil)
	if err != nil {
		log.Panic(err)
	}

	api := &Server{db, r, cert}

	// register middleware
	api.router.Use(middleware.Logger)

	// register routes from routes.go
	api.routes()

	log.Printf("Starting HTTP server on port 8000.\n")
	log.Printf("Press CTRL+C to stop.")
	log.Fatal(http.ListenAndServe(":8000", api.router))

}

// Health is a simple health check endpoint that returns HTTP 200 OK.
func (api *Server) health(w http.ResponseWriter, req *http.Request) {
	w.WriteHeader(http.StatusOK)
}
