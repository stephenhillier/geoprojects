-- migrate:up
CREATE TABLE IF NOT EXISTS project(
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL CHECK (char_length(name) < 255),
    number TEXT NOT NULL CHECK (char_length(number) < 255),
    client TEXT NOT NULL CHECK (char_length(client) < 255),
    pm TEXT NOT NULL CHECK (char_length(pm) < 255),
    location TEXT NOT NULL CHECK (char_length(location) < 255),
    organization TEXT NULL CHECK (char_length(organization) < 255),
    default_coords GEOGRAPHY(POINT,4326) NOT NULL
);

CREATE TABLE IF NOT EXISTS field_program(
    id SERIAL PRIMARY KEY,
    project INTEGER REFERENCES project(id) NOT NULL,
    start_date DATE NOT NULL,
    end_date DATE
);

-- migrate:down

