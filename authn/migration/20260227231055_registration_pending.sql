-- +goose Up
CREATE TABLE registration_pending
(
    id              CHAR(26) PRIMARY KEY DEFAULT generate_ulid(),
    email           TEXT UNIQUE NOT NULL,
    code            TEXT        NOT NULL,
    code_expires_at TIMESTAMP   NOT NULL,
    created_at      TIMESTAMP            DEFAULT CURRENT_TIMESTAMP
);

-- +goose Down
DROP TABLE registration_pending;
