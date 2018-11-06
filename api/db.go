package main

import (
	"log"

	"github.com/jmoiron/sqlx"

	// load postgres driver
	_ "github.com/lib/pq"
)

// DB represents a database with an open connection
type DB struct {
	*sqlx.DB
}

// NewDB initializes the database connection
func NewDB(connectionConfig string) (*sqlx.DB, error) {
	db, err := sqlx.Open("postgres", connectionConfig)
	if err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil {
		return nil, err
	}

	migrate(db)

	log.Println("Database connection ready.")
	return db, nil
}

func migrate(db *sqlx.DB) (migrated bool, err error) {
	check := `SELECT migrated FROM migration WHERE id=$1`
	row := db.QueryRow(check, 1)
	err = row.Scan(&migrated)

	if err == nil && migrated == true {
		// indicate that the migration does not need to occur
		return migrated, err
	}

	createUserTable := `CREATE TABLE IF NOT EXISTS users(
		id SERIAL PRIMARY KEY,
		username TEXT NOT NULL CHECK (char_length(username) < 40)	
	)`

	createProjectTable := `CREATE TABLE IF NOT EXISTS project(
		id SERIAL PRIMARY KEY,
		name TEXT NOT NULL CHECK (char_length(name) < 255),
		location TEXT NOT NULL CHECK (char_length(location) < 255),
		expired_at DATE
	)`

	// 2018-10-1

	createFieldProgramTable := `CREATE TABLE IF NOT EXISTS field_program(
		id SERIAL PRIMARY KEY,
		project INTEGER REFERENCES project(id) NOT NULL,
		start_date DATE NOT NULL,
		end_date DATE
	)`

	createDatapointTable := `CREATE TABLE IF NOT EXISTS datapoint(
		id SERIAL PRIMARY KEY,
		location GEOGRAPHY(POINT,4326) NOT NULL
	)`

	createBoreholeTable := `CREATE TABLE IF NOT EXISTS borehole(
		id SERIAL PRIMARY KEY,
		project INTEGER REFERENCES project(id) NOT NULL,
		datapoint INTEGER REFERENCES datapoint(id) NOT NULL,
		program INTEGER REFERENCES field_program(id),
		name TEXT NOT NULL CHECK (char_length(name) < 40),
		start_date DATE NOT NULL,
		end_date DATE,
		field_eng TEXT NOT NULL CHECK (char_length(name) < 80),
		UNIQUE (project, name)
	)`

	createStrataTable := `
		CREATE TABLE IF NOT EXISTS strata(
			id SERIAL PRIMARY KEY,
			borehole INTEGER REFERENCES borehole(id),
			start_depth DOUBLE PRECISION NOT NULL,
			end_depth DOUBLE PRECISION NOT NULL,
			description TEXT NOT NULL CHECK (char_length(description) < 800),
			soils TEXT NOT NULL CHECK (char_length(soils) < 200),
			moisture TEXT CHECK (char_length(moisture) < 50),
			consistency TEXT CHECK (char_length(consistency) < 50)
		)
	`

	// migrations
	createMigrationsTable := `CREATE TABLE IF NOT EXISTS migration(
		id INTEGER PRIMARY KEY,
		migrated BOOLEAN NOT NULL
	)`

	registerMigration := `INSERT INTO migration (id, migrated) VALUES (1, TRUE)`

	tx := db.MustBegin()
	tx.MustExec(createUserTable)
	tx.MustExec(createProjectTable)
	tx.MustExec(createMigrationsTable)

	// 2018-10-01

	tx.MustExec(createFieldProgramTable)
	tx.MustExec(createDatapointTable)
	tx.MustExec(createBoreholeTable)
	tx.MustExec(createStrataTable)

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
