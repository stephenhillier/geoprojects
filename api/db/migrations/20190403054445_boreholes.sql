-- migrate:up
CREATE TABLE IF NOT EXISTS datapoint(
    id SERIAL PRIMARY KEY,
    location GEOGRAPHY(POINT,4326) NOT NULL
);

CREATE TABLE IF NOT EXISTS borehole(
    id SERIAL PRIMARY KEY,
    project INTEGER NOT NULL REFERENCES project(id) ON DELETE CASCADE,
    datapoint INTEGER REFERENCES datapoint(id) NOT NULL,
    program INTEGER REFERENCES field_program(id),
    name TEXT NOT NULL CHECK (char_length(name) < 40),
    start_date DATE NOT NULL,
    end_date DATE,
    field_eng TEXT NOT NULL CHECK (char_length(field_eng) < 80),
    UNIQUE (project, name)
);

CREATE TABLE IF NOT EXISTS strata(
    id SERIAL PRIMARY KEY,
    borehole INTEGER NOT NULL REFERENCES borehole(id) ON DELETE CASCADE,
    start_depth DOUBLE PRECISION NOT NULL,
    end_depth DOUBLE PRECISION NOT NULL,
    description TEXT NOT NULL CHECK (char_length(description) < 800),
    soils TEXT NOT NULL CHECK (char_length(soils) < 200),
    moisture TEXT CHECK (char_length(moisture) < 50),
    consistency TEXT CHECK (char_length(consistency) < 50)
);

-- migrate:down

