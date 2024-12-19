-- +migrate Up
-- SQL in section 'Up' is executed when this migration is applied

CREATE TABLE IF NOT EXISTS sessions (
    id            VARCHAR(255)      NOT NULL,
    user_id       INT UNSIGNED      NULL,
    ip_address    VARCHAR(45)       NULL,
    user_agent    VARCHAR(255)      NULL,
    payload       LONGTEXT          NOT NULL,
    last_activity BIGINT UNSIGNED   NOT NULL,
    PRIMARY KEY (id),
    INDEX sessions_user_id_index (user_id),
    INDEX sessions_last_activity_index (last_activity)
) COLLATE=utf8mb4_unicode_ci;
