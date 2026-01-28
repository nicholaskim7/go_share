CREATE TABLE IF NOT EXISTS users (
    id BIGSERIAL PRIMARY KEY,
    first_name text NOT NULL,
    last_name text NOT NULL,
    user_name text NOT NULL UNIQUE,
    email text NOT NULL UNIQUE,
    password bytea NOT NULL,
    date_created TIMESTAMP(0) WITH TIME ZONE DEFAULT NOW()
);