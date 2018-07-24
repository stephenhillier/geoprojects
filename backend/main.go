package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/stephenhillier/geoprojects/backend/models"
)

// Server represents the server environment (db and router)
type Server struct {
	db models.Datastore
}

func main() {

	dbuser := os.Getenv("DBUSER")
	dbpass := os.Getenv("DBPASS")
	dbname := os.Getenv("DBNAME")
	dbhost := os.Getenv("DBHOST")

	db, err := models.NewDB(fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=disable", dbuser, dbpass, dbhost, dbname))
	if err != nil {
		log.Panic(err)
	}

	s := &Server{db}

	log.Printf("Starting HTTP server on port 8000.\n")
	log.Printf("Press CTRL+C to stop.")
	http.HandleFunc("/projects", s.projectsIndex)
	log.Fatal(http.ListenAndServe(":8000", nil))
}

func (s *Server) projectsIndex(w http.ResponseWriter, req *http.Request) {
	defer func(t time.Time) {
		log.Printf("%s: %s (%v)", req.Method, req.URL.Path, time.Since(t))
	}(time.Now())

	if req.Method != "GET" {
		http.Error(w, http.StatusText(405), 405)
		return
	}
	projects, err := s.db.AllProjects()
	if err != nil {
		http.Error(w, http.StatusText(500), 500)
	}

	response, err := json.Marshal(projects)
	if err != nil {
		http.Error(w, http.StatusText(500), 500)
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(response)
}
