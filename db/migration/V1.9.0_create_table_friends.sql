CREATE TABLE IF NOT EXISTS friends
(
    id       SERIAL PRIMARY KEY,
    user_id1 INTEGER REFERENCES users (id) NOT NULL,
    user_id2 INTEGER REFERENCES users (id) NOT NULL
);