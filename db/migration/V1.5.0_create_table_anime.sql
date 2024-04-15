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