-- migrate:up

CREATE TABLE IF NOT EXISTS soil_sample(
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL CHECK (char_length(name) < 100),
    borehole INTEGER NOT NULL REFERENCES borehole(id) ON DELETE CASCADE,
    start_depth DOUBLE PRECISION NOT NULL,
    end_depth DOUBLE PRECISION NOT NULL,
    description TEXT NOT NULL CHECK (char_length(description) < 800),
    uscs TEXT NOT NULL CHECK (char_length(uscs) < 20),
    UNIQUE (name, borehole)
);

CREATE TYPE lab_test_code AS ENUM (
    'moisture_content', 'grain_size_analysis', 'hydrometer', 'proctor', 'atterberg' 
);

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
);

CREATE TABLE IF NOT EXISTS moisture_test(
    id INTEGER REFERENCES lab_test(id) ON DELETE CASCADE PRIMARY KEY,
    tare_mass DOUBLE PRECISION NULL,
    sample_plus_tare DOUBLE PRECISION NULL,
    dry_plus_tare DOUBLE PRECISION NULL
);

CREATE TABLE IF NOT EXISTS gsa_test(
    id INTEGER REFERENCES lab_test(id) ON DELETE CASCADE PRIMARY KEY,
    tare_mass DOUBLE PRECISION NULL,
    dry_plus_tare DOUBLE PRECISION NULL,
    washed_plus_tare DOUBLE PRECISION NULL
);

CREATE TABLE IF NOT EXISTS gsa_data(
    id SERIAL PRIMARY KEY,
    test INTEGER REFERENCES gsa_test(id) ON DELETE CASCADE NOT NULL,
    pan BOOLEAN NOT NULL,
    size DOUBLE PRECISION NOT NULL,
    mass_retained DOUBLE PRECISION NULL
);


CREATE UNIQUE INDEX pan_idx ON gsa_data (test) WHERE pan;


-- migrate:down

