-- +migrate Up
-- SQL in section 'Up' is executed when this migration is applied

CREATE TABLE IF NOT EXISTS sessions (
    id            VARCHAR(255) NOT NULL
        CONSTRAINT sessions_pk
            PRIMARY KEY,
    user_id       INTEGER,
    ip_address    VARCHAR(45),
    user_agent    VARCHAR(255),
    payload       TEXT         NOT NULL,
    last_activity BIGINT       NOT NULL
);

CREATE INDEX IF NOT EXISTS sessions_last_activity_index
    ON sessions (last_activity);

CREATE INDEX IF NOT EXISTS sessions_user_id_index
    ON sessions (user_id);
