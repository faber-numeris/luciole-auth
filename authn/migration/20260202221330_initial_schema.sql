-- +goose Up
CREATE TABLE users
(
    id            CHAR(26) PRIMARY KEY DEFAULT generate_ulid(),
    email         TEXT UNIQUE NOT NULL,
    password_hash BYTEA       NOT NULL,
    first_name    TEXT        NOT NULL DEFAULT '',
    last_name     TEXT        NOT NULL DEFAULT '',
    locale        TEXT        NOT NULL DEFAULT '',
    timezone      TEXT        NOT NULL DEFAULT '',
    created_at    TIMESTAMP            DEFAULT CURRENT_TIMESTAMP,
    updated_at    TIMESTAMP            DEFAULT CURRENT_TIMESTAMP,
    deleted_at    TIMESTAMP
);


-- +goose Down
SELECT 'down SQL query';
drop table users;

