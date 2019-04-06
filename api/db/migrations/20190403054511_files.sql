-- migrate:up
CREATE TYPE project_file_code AS ENUM (
    'report', 'lab_report', 'calculation', 'proposal', 'budget', 'field_data', 'other'
);

CREATE TABLE IF NOT EXISTS project_file(
    id SERIAL PRIMARY KEY,
    project INTEGER NOT NULL REFERENCES project(id) ON DELETE CASCADE,
    category project_file_code NOT NULL,
    file BYTEA NOT NULL,
    filename TEXT NOT NULL CHECK (char_length(filename) < 250),
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    created_by TEXT NOT NULL CHECK (char_length(created_by) < 250),
    expired_at TIMESTAMP NULL
);

-- migrate:down

