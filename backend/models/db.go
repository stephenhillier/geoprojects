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
	CreateDescription(Description) (Description, error)
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

	createDescriptionTable := `CREATE TABLE IF NOT EXISTS description(
		id serial primary key,
		original text not null check(char_length(original) < 255),
		"primary" soil not null,
		secondary soil,
		consistency consistency,
		moisture moisture
	)`

	createMigrationsTable := `CREATE TABLE IF NOT EXISTS migration(
		id int primary key,
		migrated boolean not null
	)`

	registerMigration := `INSERT INTO migration (id, migrated) VALUES (1, TRUE)`

	tx := db.MustBegin()
	tx.MustExec(createSoilCodes)
	tx.MustExec(createConsCodes)
	tx.MustExec(createMoisCodes)
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
