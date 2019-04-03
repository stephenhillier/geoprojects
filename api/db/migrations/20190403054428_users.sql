-- migrate:up
CREATE TABLE IF NOT EXISTS users(
    id SERIAL PRIMARY KEY,
    username TEXT NOT NULL CHECK (char_length(username) < 40)	
);

-- migrate:down

