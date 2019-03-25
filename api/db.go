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

	installPostGISExtension := `CREATE EXTENSION IF NOT EXISTS POSTGIS`

	createUserTable := `CREATE TABLE IF NOT EXISTS users(
		id SERIAL PRIMARY KEY,
		username TEXT NOT NULL CHECK (char_length(username) < 40)	
	)`

	createProjectTable := `CREATE TABLE IF NOT EXISTS project(
		id SERIAL PRIMARY KEY,
		name TEXT NOT NULL CHECK (char_length(name) < 255),
		number TEXT NOT NULL CHECK (char_length(number) < 255),
		client TEXT NOT NULL CHECK (char_length(client) < 255),
		pm TEXT NOT NULL CHECK (char_length(pm) < 255),
		location TEXT NOT NULL CHECK (char_length(location) < 255),
		organization TEXT NULL CHECK (char_length(organization) < 255),
		default_coords GEOGRAPHY(POINT,4326) NOT NULL
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
		project INTEGER NOT NULL REFERENCES project(id) ON DELETE CASCADE,
		datapoint INTEGER REFERENCES datapoint(id) NOT NULL,
		program INTEGER REFERENCES field_program(id),
		name TEXT NOT NULL CHECK (char_length(name) < 40),
		start_date DATE NOT NULL,
		end_date DATE,
		field_eng TEXT NOT NULL CHECK (char_length(field_eng) < 80),
		UNIQUE (project, name)
	)`

	createStrataTable := `
		CREATE TABLE IF NOT EXISTS strata(
			id SERIAL PRIMARY KEY,
			borehole INTEGER NOT NULL REFERENCES borehole(id) ON DELETE CASCADE,
			start_depth DOUBLE PRECISION NOT NULL,
			end_depth DOUBLE PRECISION NOT NULL,
			description TEXT NOT NULL CHECK (char_length(description) < 800),
			soils TEXT NOT NULL CHECK (char_length(soils) < 200),
			moisture TEXT CHECK (char_length(moisture) < 50),
			consistency TEXT CHECK (char_length(consistency) < 50)
		)
	`
	// 2018-12-09
	createSampleTable := `
		CREATE TABLE IF NOT EXISTS soil_sample(
			id SERIAL PRIMARY KEY,
			name TEXT NOT NULL CHECK (char_length(name) < 100),
			borehole INTEGER NOT NULL REFERENCES borehole(id) ON DELETE CASCADE,
			start_depth DOUBLE PRECISION NOT NULL,
			end_depth DOUBLE PRECISION NOT NULL,
			description TEXT NOT NULL CHECK (char_length(description) < 800),
			uscs TEXT NOT NULL CHECK (char_length(uscs) < 20),
			UNIQUE (name, borehole)
		)
	`

	createTestType := `
		CREATE TYPE lab_test_code AS ENUM (
			'moisture_content', 'grain_size_analysis', 'hydrometer', 'proctor', 'atterberg' 
		)
	`

	createLabTestTable := `
		CREATE TABLE IF NOT EXISTS lab_test(
			id SERIAL PRIMARY KEY,
			name TEXT NULL CHECK (char_length(name) < 100),
			type lab_test_code NOT NULL,
			start_date DATE NULL,
			end_date DATE NULL,
			performed_by TEXT NULL CHECK (char_length(performed_by) < 200),
			sample INTEGER NOT NULL REFERENCES soil_sample(id),
			checked_date DATE NULL,
			checked_by TEXT NULL CHECK (char_length(checked_by) < 200)
		)
	`

	// 2018-12-19

	createMoistureTestTable := `
		CREATE TABLE IF NOT EXISTS moisture_test(
			id INTEGER REFERENCES lab_test(id) ON DELETE CASCADE PRIMARY KEY,
			tare_mass DOUBLE PRECISION NULL,
			sample_plus_tare DOUBLE PRECISION NULL,
			dry_plus_tare DOUBLE PRECISION NULL
		)
	`

	createGSATestTable := `
		CREATE TABLE IF NOT EXISTS gsa_test(
			id INTEGER REFERENCES lab_test(id) ON DELETE CASCADE PRIMARY KEY,
			tare_mass DOUBLE PRECISION NULL,
			dry_plus_tare DOUBLE PRECISION NULL,
			washed_plus_tare DOUBLE PRECISION NULL
		)
	`

	createGSADataTable := `
		CREATE TABLE IF NOT EXISTS gsa_data(
			id SERIAL PRIMARY KEY,
			test INTEGER REFERENCES gsa_test(id) ON DELETE CASCADE NOT NULL,
			pan BOOLEAN NOT NULL,
			size DOUBLE PRECISION NOT NULL,
			mass_retained DOUBLE PRECISION NULL
		)
	`

	createGSADataUniqueIndex := `
		CREATE UNIQUE INDEX pan_idx ON gsa_data (test) WHERE pan
	`

	createFileType := `
		CREATE TYPE project_file_code AS ENUM (
			'report', 'lab_report', 'calculation', 'proposal', 'budget', 'field_data', 'other'
		)
	`

	createFileTable := `
		CREATE TABLE IF NOT EXISTS project_file(
			id SERIAL PRIMARY KEY,
			project INTEGER NOT NULL REFERENCES project(id) ON DELETE CASCADE,
			category project_file_code NOT NULL,
			file BYTEA NOT NULL,
			filename TEXT NOT NULL CHECK (char_length(filename) < 250),
			created_at TIMESTAMP NOT NULL DEFAULT NOW(),
			created_by TEXT NOT NULL CHECK (char_length(created_by) < 250)
		)
	`

	// migrations
	createMigrationsTable := `CREATE TABLE IF NOT EXISTS migration(
		id INTEGER PRIMARY KEY,
		migrated BOOLEAN NOT NULL
	)`

	registerMigration := `INSERT INTO migration (id, migrated) VALUES (1, TRUE)`

	tx := db.MustBegin()
	tx.MustExec(installPostGISExtension)
	tx.MustExec(createUserTable)
	tx.MustExec(createProjectTable)
	tx.MustExec(createMigrationsTable)

	// 2018-10-01

	tx.MustExec(createFieldProgramTable)
	tx.MustExec(createDatapointTable)
	tx.MustExec(createBoreholeTable)
	tx.MustExec(createStrataTable)

	// 2018-12-09
	tx.MustExec(createSampleTable)
	tx.MustExec(createTestType)
	tx.MustExec(createLabTestTable)

	// 2018-12-19
	tx.MustExec(createMoistureTestTable)
	tx.MustExec(createGSATestTable)
	tx.MustExec(createGSADataTable)
	tx.MustExec(createGSADataUniqueIndex)

	// 2019-3-23
	tx.MustExec(createFileType)
	tx.MustExec(createFileTable)

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
