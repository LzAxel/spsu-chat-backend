CREATE TABLE chats (
    id BIGSERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    creator_id BIGINT NOT NULL REFERENCES users(id),
    type SMALLINT NOT NULL,
    password_hash BYTEA,
    created_at TIMESTAMPTZ NOT NULL
);