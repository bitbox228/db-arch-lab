CREATE TABLE IF NOT EXISTS reactions
(
    id        SERIAL PRIMARY KEY,
    review_id INTEGER REFERENCES reviews (id) NOT NULL,
    user_id   INTEGER REFERENCES users (id)   NOT NULL,
    is_like   BOOLEAN                         NOT NULL
);