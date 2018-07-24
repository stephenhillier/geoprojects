package models

import (
	"log"

	"github.com/jmoiron/sqlx"

	// load postgres driver
	_ "github.com/lib/pq"
)

// Datastore is the collection of model handlers available to the server
type Datastore interface {
	AllProjects() ([]*Project, error)
}

// DB represents a database with an open connection
type DB struct {
	*sqlx.DB
}

// NewDB initializes the database connection
func NewDB(connectionConfig string) (*DB, error) {
	db, err := sqlx.Open("postgres", connectionConfig)
	if err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil {
		return nil, err
	}
	log.Println("Database connection ready.")
	return &DB{db}, nil
}
