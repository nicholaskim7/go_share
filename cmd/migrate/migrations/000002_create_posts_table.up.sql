CREATE TABLE IF NOT EXISTS posts (
    id BIGSERIAL PRIMARY KEY,
    title text NOT NULL,
    body text NOT NULL,
    user_id bigint NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    files text[],
    date_created TIMESTAMP(0) WITH TIME ZONE DEFAULT NOW()
);