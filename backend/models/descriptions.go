package models

// Description is a consistent format for describing
type Description struct {
	ID          uint64 `json:"id" db:"id"`
	Original    string `json:"original" db:"original"`
	Primary     string `json:"primary" db:"primary"`
	Secondary   string `json:"secondary" db:"secondary"`
	Consistency string `json:"consistency" db:"consistency"`
	Moisture    string `json:"moisture" db:"moisture"`
}

// CreateDescription saves a Description (after being parsed) along with the original text string
func (db *DB) CreateDescription(r Description) (d Description, err error) {
	query := `INSERT INTO description (original, "primary", secondary, consistency, moisture) VALUES ($1, $2, $3, $4, $5) RETURNING *`
	err = db.QueryRowx(
		query,
		r.Original,
		r.Primary,
		r.Secondary,
		r.Consistency,
		r.Moisture,
	).StructScan(&d)

	return d, err
}
