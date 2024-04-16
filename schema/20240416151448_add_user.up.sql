CREATE TABLE users (
    id BIGSERIAL PRIMARY KEY,
    username TEXT NOT NULL UNIQUE,
    display_name TEXT NOT NULL,
    password_hash BYTEA NOT NULL,
    type SMALLINT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL
)