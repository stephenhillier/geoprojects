package repository

import (
	"github.com/stephenhillier/geoprojects/earthworks/db"
)

// NewInstrumentationRepo returns a PostgresRepo with a database connection
// This method can be called with either a sqlx.DB or a sqlx.Tx (transaction)
func NewInstrumentationRepo(database *db.Datastore) *PostgresRepo {
	return &PostgresRepo{
		conn: database,
	}
}

// PostgresRepo has a database connection and methods to interact with instrumentation in
// the database.
type PostgresRepo struct {
	conn *db.Datastore
}
