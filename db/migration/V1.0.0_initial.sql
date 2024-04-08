DO
$$
    BEGIN
        IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'anime_type') THEN
            CREATE TYPE anime_type AS ENUM (
                'SERIES',
                'MOVIE',
                'OVA'
                );
        END IF;
        IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'anime_status') THEN
            CREATE TYPE anime_status AS ENUM (
                'ONGOING',
                'PLANNED',
                'RELEASED'
                );
        END IF;
        IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'notification_type') THEN
            CREATE TYPE notification_type AS ENUM (
                'FRIEND_REQUEST',
                'NEW_EPISODE',
                'NEW_MESSAGE'
                );
        END IF;
        IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'user_anime_list') THEN
            CREATE TYPE user_anime_list AS ENUM (
                'WATCHING',
                'WATCHED',
                'WANT_TO_WATCH',
                'DROPPED',
                'DEFERRED',
                'REVISING'
                );
        END IF;
    END
$$;

CREATE TABLE IF NOT EXISTS users
(
    id            SERIAL PRIMARY KEY,
    email         VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255)        NOT NULL,
    nickname      VARCHAR(255)        NOT NULL,
    avatar_url    VARCHAR(255)        NOT NULL,
    is_private    BOOLEAN             NOT NULL
);

CREATE TABLE IF NOT EXISTS anime
(
    id           SERIAL PRIMARY KEY,
    title        VARCHAR(255) NOT NULL,
    release_date DATE         NOT NULL,
    rating       FLOAT        NOT NULL,
    genre        VARCHAR(255) NOT NULL,
    type         anime_type   NOT NULL,
    studio       VARCHAR(255) NOT NULL,
    status       anime_status NOT NULL,
    age_rating   VARCHAR(5)   NOT NULL,
    cover_url    VARCHAR(255) NOT NULL
);

CREATE TABLE IF NOT EXISTS user_anime_status
(
    id            SERIAL PRIMARY KEY,
    anime_id      INTEGER REFERENCES anime (id) NOT NULL,
    user_id       INTEGER REFERENCES users (id) NOT NULL,
    list          user_anime_list               NOT NULL,
    is_subscribed BOOLEAN                       NOT NULL
);

CREATE TABLE IF NOT EXISTS anime_series
(
    id            SERIAL PRIMARY KEY,
    anime_id      INTEGER REFERENCES anime (id) NOT NULL,
    series_url    VARCHAR(255)                  NOT NULL,
    seconds_count INTEGER                       NOT NULL
);

CREATE TABLE IF NOT EXISTS reviews
(
    id       SERIAL PRIMARY KEY,
    user_id  INTEGER REFERENCES users (id) NOT NULL,
    anime_id INTEGER REFERENCES anime (id) NOT NULL,
    rating   FLOAT                         NOT NULL,
    text     VARCHAR(1000)                 NOT NULL
);

CREATE TABLE IF NOT EXISTS friends
(
    id       SERIAL PRIMARY KEY,
    user_id1 INTEGER REFERENCES users (id) NOT NULL,
    user_id2 INTEGER REFERENCES users (id) NOT NULL
);

CREATE TABLE IF NOT EXISTS messages
(
    id          SERIAL PRIMARY KEY,
    sender_id   INTEGER REFERENCES users (id) NOT NULL,
    receiver_id INTEGER REFERENCES users (id) NOT NULL,
    text        VARCHAR(1000)                 NOT NULL,
    file_url    VARCHAR(255)                  NOT NULL,
    time        TIMESTAMP                     NOT NULL
);

CREATE TABLE IF NOT EXISTS achievements
(
    id          SERIAL PRIMARY KEY,
    anime_id    INTEGER REFERENCES anime (id) NOT NULL,
    name        VARCHAR(255)                  NOT NULL,
    description VARCHAR(255)                  NOT NULL
);

CREATE TABLE IF NOT EXISTS user_achievements
(
    id             SERIAL PRIMARY KEY,
    achievement_id INTEGER REFERENCES achievements (id) NOT NULL,
    user_id        INTEGER REFERENCES users (id)        NOT NULL,
    time           TIMESTAMP                            NOT NULL
);

CREATE TABLE IF NOT EXISTS notifications
(
    id      SERIAL PRIMARY KEY,
    user_id INTEGER REFERENCES users (id) NOT NULL,
    type    notification_type             NOT NULL,
    body    JSON                          NOT NULL,
    time    TIMESTAMP                     NOT NULL
);

CREATE TABLE IF NOT EXISTS reactions
(
    id        SERIAL PRIMARY KEY,
    review_id INTEGER REFERENCES reviews (id) NOT NULL,
    user_id   INTEGER REFERENCES users (id)   NOT NULL,
    is_like   BOOLEAN                         NOT NULL
);
