CREATE TABLE IF NOT EXISTS reviews
(
    id       SERIAL PRIMARY KEY,
    user_id  INTEGER REFERENCES users (id) NOT NULL,
    anime_id INTEGER REFERENCES anime (id) NOT NULL,
    rating   FLOAT                         NOT NULL,
    text     VARCHAR(1000)                 NOT NULL
);