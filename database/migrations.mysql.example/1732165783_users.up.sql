CREATE TABLE IF NOT EXISTS users (
    id             INT UNSIGNED auto_increment,
    username       VARCHAR(255)                       NOT NULL,
    password       VARCHAR(255)                       NOT NULL,
    remember_token VARCHAR(100)                       NULL,
    created_at     DATETIME default CURRENT_TIMESTAMP NOT NULL,
    updated_at     DATETIME default CURRENT_TIMESTAMP NOT NULL,
    PRIMARY KEY (id),
    UNIQUE users_username_unique (username)
) COLLATE=utf8mb4_unicode_ci;