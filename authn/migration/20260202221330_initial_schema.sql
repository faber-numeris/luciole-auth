-- +goose Up
CREATE TABLE users
(
    id            CHAR(26) PRIMARY KEY DEFAULT generate_ulid(),
    username      TEXT UNIQUE NOT NULL,
    email         TEXT UNIQUE NOT NULL,
    password_hash BYTEA       NOT NULL,
    created_at    TIMESTAMP            DEFAULT CURRENT_TIMESTAMP,
    updated_at    TIMESTAMP            DEFAULT CURRENT_TIMESTAMP,
    deleted_at    TIMESTAMP
);


-- +goose Down
SELECT 'down SQL query';
drop table users;

