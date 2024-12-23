CREATE TABLE IF NOT EXISTS users (
    id integer constraint users_pk primary key autoincrement,
    username text not null,
    password text not null,
    remember_token text,
    created_at datetime not null default current_timestamp,
    updated_at datetime not null default current_timestamp
);

CREATE UNIQUE INDEX users_username_uindex ON users (username);
