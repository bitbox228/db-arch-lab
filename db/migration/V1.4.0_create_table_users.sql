CREATE TABLE IF NOT EXISTS users
(
    id            SERIAL PRIMARY KEY,
    email         VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255)        NOT NULL,
    nickname      VARCHAR(255)        NOT NULL,
    avatar_url    VARCHAR(255)        NOT NULL,
    is_private    BOOLEAN             NOT NULL
);