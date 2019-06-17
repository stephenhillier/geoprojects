-- migrate:up
CREATE TABLE instrument (
    id SERIAL PRIMARY KEY,
    project INTEGER NOT NULL REFERENCES project(id) ON DELETE CASCADE,
    datapoint INTEGER REFERENCES datapoint(id) NOT NULL,
    program INTEGER REFERENCES field_program(id),
    name TEXT NOT NULL CHECK (char_length(name) < 40),
    device_id UUID NULL,
    install_date DATE NOT NULL,
    field_eng TEXT NOT NULL CHECK (char_length(field_eng) < 80),
    UNIQUE (project, name)
);

CREATE TABLE time_series_data (
    id BIGSERIAL PRIMARY KEY,
    device_id UUID NOT NULL,
    instrument INTEGER NULL REFERENCES instrument(id),
    series INTEGER DEFAULT 0,
    time TIMESTAMPTZ,
    value DOUBLE PRECISION
);

-- migrate:down

