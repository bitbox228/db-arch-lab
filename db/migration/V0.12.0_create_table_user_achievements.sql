CREATE TABLE IF NOT EXISTS user_achievements
(
    id             SERIAL PRIMARY KEY,
    achievement_id INTEGER REFERENCES achievements (id) NOT NULL,
    user_id        INTEGER REFERENCES users (id)        NOT NULL,
    time           TIMESTAMP                            NOT NULL
);