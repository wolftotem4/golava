-- +migrate Up
-- SQL in section 'Up' is executed when this migration is applied

CREATE TABLE IF NOT EXISTS sessions (
    id text constraint session_pk primary key,
    user_id integer null,
    ip_address text null,
    user_agent text null,
    payload text not null,
    last_activity integer not null
);

CREATE INDEX sessions_last_activity_index ON sessions (last_activity);
