CREATE TABLE IF NOT EXISTS user_anime_status
(
    id            SERIAL PRIMARY KEY,
    anime_id      INTEGER REFERENCES anime (id) NOT NULL,
    user_id       INTEGER REFERENCES users (id) NOT NULL,
    list          user_anime_list               NOT NULL,
    is_subscribed BOOLEAN                       NOT NULL
);