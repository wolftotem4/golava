CREATE TABLE IF NOT EXISTS users (
    id             SERIAL
        CONSTRAINT users_pk
            PRIMARY KEY,
    username       VARCHAR(255) NOT NULL
        CONSTRAINT users_username_unique
            UNIQUE,
    password       VARCHAR(255) NOT NULL,
    remember_token VARCHAR(100),
    created_at     TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at     TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);