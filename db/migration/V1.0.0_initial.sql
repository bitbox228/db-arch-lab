DO
$$
    BEGIN
        IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'anime_status') THEN
            CREATE TYPE anime_status AS ENUM (
                'ongoing',
                'planned',
                'released'
                );
        END IF;
        IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'notification_type') THEN
            CREATE TYPE notification_type AS ENUM (
                'friend_request',
                'new_episode',
                'new_message'
                );
        END IF;
        IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'user_anime_list') THEN
            CREATE TYPE user_anime_list AS ENUM (
                'watching',
                'watched',
                'want_to_watch',
                'dropped',
                'deferred',
                'revising'
                );
        END IF;
    END
$$;

CREATE TABLE IF NOT EXISTS users
(
    id            SERIAL PRIMARY KEY,
    email         VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255)        NOT NULL,
    nickname      VARCHAR(255),
    avatar_url    VARCHAR(255),
    is_private    BOOLEAN             NOT NULL DEFAULT false
);

CREATE TABLE IF NOT EXISTS anime
(
    id           SERIAL PRIMARY KEY,
    title        VARCHAR(255) NOT NULL,
    release_date DATE,
    rating       FLOAT,
    genre        VARCHAR(255),
    type         anime_type,
    studio       VARCHAR(255),
    status       anime_status,
    age_rating   INTEGER,
    cover_url    VARCHAR(255)
);

CREATE TABLE IF NOT EXISTS user_anime_status
(
    id            SERIAL PRIMARY KEY,
    anime_id      INTEGER REFERENCES anime (id),
    user_id       INTEGER REFERENCES users (id),
    list          user_anime_list,
    is_subscribed BOOLEAN NOT NULL DEFAULT false
);

CREATE TABLE IF NOT EXISTS anime_series
(
    id            SERIAL PRIMARY KEY,
    anime_id      INTEGER REFERENCES anime (id),
    series_url    VARCHAR(255),
    seconds_count INTEGER
);

CREATE TABLE IF NOT EXISTS reviews
(
    id       SERIAL PRIMARY KEY,
    user_id  INTEGER REFERENCES users (id),
    anime_id INTEGER REFERENCES anime (id),
    rating   INTEGER,
    text     VARCHAR(255)
);

CREATE TABLE IF NOT EXISTS friends
(
    id       SERIAL PRIMARY KEY,
    user_id1 INTEGER REFERENCES users (id),
    user_id2 INTEGER REFERENCES users (id)
);

CREATE TABLE IF NOT EXISTS messages
(
    id          SERIAL PRIMARY KEY,
    sender_id   INTEGER REFERENCES users (id),
    receiver_id INTEGER REFERENCES users (id),
    text        VARCHAR(255),
    file_url    VARCHAR(255),
    time        TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS achievements
(
    id          SERIAL PRIMARY KEY,
    anime_id    INTEGER REFERENCES anime (id),
    name        VARCHAR(255),
    description VARCHAR(255)
);

CREATE TABLE IF NOT EXISTS user_achievements
(
    id             SERIAL PRIMARY KEY,
    achievement_id INTEGER REFERENCES achievements (id),
    user_id        INTEGER REFERENCES users (id),
    time           TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS notifications
(
    id      SERIAL PRIMARY KEY,
    user_id INTEGER REFERENCES users (id),
    type    notification_type,
    body    JSON,
    time    TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS reactions
(
    id        SERIAL PRIMARY KEY,
    review_id INTEGER REFERENCES reviews (id),
    user_id   INTEGER REFERENCES users (id),
    is_like   BOOLEAN NOT NULL DEFAULT false
);
