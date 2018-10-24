package field

// Program represents one fieldwork job within a project.
// For example, a single trip that an engineer makes to collect field samples
// could be represented as a new Program record. There may be several new
// Datapoints (as well as Boreholes, Instruments) associated with that
// individual Program.
type Program struct {
	ID        int64    `json:"id"`
	Project   int64    `json:"project"`
	StartDate NullDate `json:"start_date" db:"start_date"`
	EndDate   NullDate `json:"end_date" db:"end_date"`
}

// ProgramCreateRequest is the data a user should submit to create a program
type ProgramCreateRequest struct {
	Project   int64    `json:"project"`
	StartDate NullDate `json:"start_date" db:"start_date" schema:"start_date"`
	EndDate   NullDate `json:"end_date" db:"end_date" schema:"end_date"`
}

// ProgramRepository is the set of methods that are available for interacting
// with field program records
type ProgramRepository interface {
	ListPrograms() ([]Program, error)
	CreateProgram(fp ProgramCreateRequest) (Program, error)
	GetProgram(programID int) (Program, error)
}

// Field Program database methods

// ListPrograms returns a list of field program records
func (db *datastore) ListPrograms() ([]Program, error) {
	programs := []Program{}
	query := `SELECT id, project, start_date, end_date FROM field_program`

	err := db.Select(&programs, query)
	if err != nil {
		return []Program{}, err
	}

	return programs, nil
}

// CreateProgram creates a field program record
func (db *datastore) CreateProgram(fp ProgramCreateRequest) (Program, error) {
	query := `INSERT INTO field_program (project, start_date, end_date) VALUES ($1, $2, $3) RETURNING id, project, start_date, end_date`
	created := Program{}
	err := db.Get(&created, query, fp.Project, fp.StartDate, fp.EndDate)
	if err != nil {
		return Program{}, err
	}
	return created, nil
}

// GetProgram retrieves a single field program record
func (db *datastore) GetProgram(programID int) (Program, error) {
	p := Program{}
	query := `SELECT id, project, start_date, end_date FROM field_program WHERE id=$1`
	err := db.Get(&p, query, programID)
	if err != nil {
		return Program{}, err
	}
	return p, nil
}
