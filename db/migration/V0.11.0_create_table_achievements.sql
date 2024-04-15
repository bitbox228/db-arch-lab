CREATE TABLE IF NOT EXISTS achievements
(
    id          SERIAL PRIMARY KEY,
    anime_id    INTEGER REFERENCES anime (id) NOT NULL,
    name        VARCHAR(255)                  NOT NULL,
    description VARCHAR(255)                  NOT NULL
);