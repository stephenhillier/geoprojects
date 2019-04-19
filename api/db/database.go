package db

import (
	"log"
	"time"

	"github.com/jmoiron/sqlx"

	// register postgres driver
	_ "github.com/lib/pq"
)

// Datastore represents a database with an open connection
type Datastore struct {
	*sqlx.DB
}

// Config holds database connection info
type Config struct {
	Conn   string // connection string e.g. postgres://user:pass@127.0.0.1:5432/mydb?sslmode=disable
	Driver string // name of the database driver e.g. postgres
}

// NewDB initializes the database connection
func NewDB(config Config) (*Datastore, error) {

	var db *sqlx.DB
	var err error

	// Start attempting to make a database connection.
	// Break the loop after a connection is successfully made.
	for {
		db, err = sqlx.Open(config.Driver, config.Conn)
		if err != nil {
			log.Println(err)
		}

		err = db.Ping()
		if err == nil {
			break
		}
		log.Println(err)
		log.Println("Waiting for database to become available")
		time.Sleep(10 * time.Second)
	}

	log.Println("Database connection ready.")
	return &Datastore{db}, nil
}
