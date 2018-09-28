package models

import (
	"log"

	"github.com/jmoiron/sqlx"

	// load postgres driver
	_ "github.com/lib/pq"
)

// Datastore is the collection of model handlers available to the server
// a model handler must be referenced here to be accessible by e.g. `api.db.AllProjects()`
type Datastore interface {
	CreateDescription(Description) (Description, error)
	AllProjects() ([]*Project, error)
	CreateProject(Project) (Project, error)
	RetrieveProject(int) (Project, error)
}

// DB represents a database with an open connection
type DB struct {
	*sqlx.DB
}

// NewDB initializes the database connection
func NewDB(connectionConfig string) (*DB, error) {
	open, err := sqlx.Open("postgres", connectionConfig)
	if err != nil {
		return nil, err
	}

	db := &DB{open}

	if err = db.Ping(); err != nil {
		return nil, err
	}

	db.migrate()

	log.Println("Database connection ready.")
	return db, nil
}

func (db *DB) migrate() (migrated bool, err error) {
	check := `SELECT migrated FROM migration WHERE id=$1`
	row := db.QueryRow(check, 1)
	err = row.Scan(&migrated)

	if err == nil && migrated == true {
		// indicate that the migration does not need to occur
		return migrated, err
	}

	createSoilCodes := `CREATE TYPE soil AS ENUM ('sand', 'gravel', 'silt', 'clay', 'cobbles', '')`
	createConsCodes := `CREATE TYPE consistency AS ENUM ('loose', 'soft', 'firm', 'compact', 'hard', 'dense', '')`
	createMoisCodes := `CREATE TYPE moisture AS ENUM ('very dry', 'dry', 'damp', 'moist', 'wet', 'very wet', '')`

	createUserTable := `CREATE TABLE IF NOT EXISTS users(
		id SERIAL PRIMARY KEY,
		username TEXT NOT NULL CHECK (char_length(username) < 40)	
	)`

	createDescriptionTable := `CREATE TABLE IF NOT EXISTS description(
		id SERIAL PRIMARY KEY,
		original TEXT NOT NULL CHECK (char_length(original) < 255),
		"primary" soil NOT NULL,
		secondary soil,
		consistency consistency,
		moisture moisture
	)`

	createProjectTable := `CREATE TABLE IF NOT EXISTS project(
		id SERIAL PRIMARY KEY,
		name TEXT NOT NULL CHECK (char_length(name) < 255),
		location TEXT NOT NULL CHECK (char_length(location) < 255),
		pm INTEGER REFERENCES users(id) ON DELETE SET NULL
	)`

	createMigrationsTable := `CREATE TABLE IF NOT EXISTS migration(
		id INTEGER PRIMARY KEY,
		migrated BOOLEAN NOT NULL
	)`

	registerMigration := `INSERT INTO migration (id, migrated) VALUES (1, TRUE)`

	tx := db.MustBegin()
	tx.MustExec(createSoilCodes)
	tx.MustExec(createConsCodes)
	tx.MustExec(createMoisCodes)
	tx.MustExec(createUserTable)
	tx.MustExec(createProjectTable)
	tx.MustExec(createDescriptionTable)
	tx.MustExec(createMigrationsTable)
	tx.MustExec(registerMigration)
	err = tx.Commit()
	if err != nil {
		return migrated, err
	}
	log.Println("Database migrated.")

	row = db.QueryRow(check, 1)
	err = row.Scan(&migrated)
	return migrated, err
}
