package main

import (
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/jmoiron/sqlx"

	// register postgres driver
	_ "github.com/lib/pq"
)

// DB represents a database with an open connection
type DB struct {
	*sqlx.DB
}

// NewDB initializes the database connection
func (s *server) NewDB() (*sqlx.DB, error) {

	var db *sqlx.DB
	var err error
	var connectionConfig string

	//
	if s.config.dbdriver == "postgres" {
		connectionConfig = fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s", s.config.dbuser, s.config.dbpass, s.config.dbhost, s.config.dbport, s.config.dbname, s.config.dbsslmode)
	} else if s.config.dbdriver == "cloudsqlpostgres" {
		connectionConfig = fmt.Sprintf("host=%s user=%s dbname=%s password=%s sslmode=%s", s.config.dbhost, s.config.dbuser, s.config.dbname, s.config.dbpass, s.config.dbsslmode)
	} else {
		return db, errors.New("Invalid database driver specified. Supported drivers: postgres, cloudsqlpostgres")
	}

	for {
		db, err = sqlx.Open(s.config.dbdriver, connectionConfig)
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
	return db, nil
}
